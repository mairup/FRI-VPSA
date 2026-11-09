package main

import (
	"fmt"
	rand "math/rand"
	"slices"
	"strings"
	"sync"
	"time"
)

const TEMP_MIN = -25.0
const TEMP_MAX = 50.0

const TLAK_MIN = 950.0
const TLAK_MAX = 1050.0

const VLAGA_MIN = 0.0
const VLAGA_MAX = 100.0

const MERITEV_INTERVAL = time.Millisecond * 1000

var VRSTE_MERITEV = []string{"temperatura", "tlak", "vlaga"}

type Postaja struct {
	ime     string
	kraj    string
	health  [3]uint8
	cli     [3]CLI
	meritve []Meritev
}

func (p Postaja) is_alive() bool {
	return p.health[0] > 0 || p.health[1] > 0 || p.health[2] > 0
}

type Meritev struct {
	vrsta    string
	vrednost float32
}

type CLI struct {
	mean  float64
	sigma float64
}

type Sporocilo struct {
	source_id string
	meritve   []Meritev
}

// Podatki o meritvah (srednja vrednost in standardni odklon) so izračunani iz 30-dnevnih urnih meritev
// prek Open-Meteo API (https://open-meteo.com/):
// Primer za Ljubljano: https://api.open-meteo.com/v1/forecast?latitude=46.0569&longitude=14.5058&past_days=30&forecast_days=0&hourly=temperature_2m,pressure_msl,relative_humidity_2m
// Format poizvedbe: https://api.open-meteo.com/v1/forecast?latitude={lat}&longitude={lon}&past_days=30&forecast_days=0&hourly=temperature_2m,pressure_msl,relative_humidity_2m
// Datum pridobitve podatkov: 03. 11. 2026

var postaje_code_names = []struct {
	slug  string
	kraj  string
	temp  CLI
	tlak  CLI
	vlaga CLI
}{
	{"LJ", "Ljubljana", CLI{11.0, 4.0}, CLI{1021.8, 5.1}, CLI{82.9, 15.4}},
	{"KR", "Kranj", CLI{10.2, 4.0}, CLI{1021.8, 5.1}, CLI{84.5, 13.2}},
	{"NM", "Novo mesto", CLI{10.7, 4.9}, CLI{1021.8, 5.1}, CLI{84.3, 17.0}},
	{"KP", "Koper", CLI{15.0, 4.1}, CLI{1019.5, 4.8}, CLI{70.1, 16.2}},
	{"CE", "Celje", CLI{10.4, 4.7}, CLI{1021.8, 5.2}, CLI{83.3, 16.0}},
	{"GO", "Goriška Brda", CLI{14.4, 3.8}, CLI{1019.8, 4.8}, CLI{70.2, 13.6}},
	{"MB", "Maribor", CLI{11.4, 4.5}, CLI{1021.5, 5.4}, CLI{75.9, 16.0}},
	{"SG", "Slovenj Gradec", CLI{10.1, 4.6}, CLI{1021.9, 5.3}, CLI{81.3, 13.6}},
	{"MS", "Murska Sobota", CLI{10.4, 5.1}, CLI{1021.6, 5.3}, CLI{79.3, 17.3}},
}

var stop_channel = make(chan struct{})
var strike_channel = make(chan Strike, 10)
var rw_lock sync.RWMutex

type Strike struct {
	id  string
	dmg uint8
}

func try_for_health_pack() uint8 {
	roll := rand.Intn(20)

	if roll == 0 {
		fmt.Printf("\n🎉🎉🎉 MAXIMUM WIN!!! +150 hp 🎉🎉🎉 --> ")
		return 150
	} else if roll > 17 {
		fmt.Printf("\n🎉 bonus health: +50 hp 🎉 --> ")
		return 50
	}
	return 0
}

func keyboard_listener() {
	fmt.Scanln()
	close(stop_channel)
}

func izpis_stanja_postaj(postaje map[string]Postaja) {
	var ids []string
	for id := range postaje {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	fmt.Println("\nSTANJE POSTAJ:")
	for _, id := range ids {
		postaja := postaje[id]
		if !postaja.is_alive() {
			fmt.Printf("  [%s] %s | -- | -- | -- ( 💀 )\n", id, postaja.kraj)
			continue
		}

		var latest [3]*Meritev
		for i := len(postaja.meritve) - 1; i >= 0; i-- {
			m := &postaja.meritve[i]
			for idx, v := range VRSTE_MERITEV {
				if latest[idx] == nil && m.vrsta == v {
					latest[idx] = m
				}
			}
			if latest[0] != nil && latest[1] != nil && latest[2] != nil {
				break
			}
		}

		var sensorVals [3]string
		for idx := range 3 {
			if postaja.health[idx] == 0 {
				sensorVals[idx] = "--"
			} else if latest[idx] != nil {
				sensorVals[idx] = fmt.Sprintf("%.2f %s", latest[idx].vrednost, enota_meritve(latest[idx].vrsta))
			} else {
				sensorVals[idx] = "..."
			}
		}

		fmt.Printf("  [%s] %s | %s | %s | %s (hp: %v)\n",
			id,
			postaja.kraj,
			sensorVals[0],
			sensorVals[1],
			sensorVals[2],
			postaja.health,
		)
	}
	fmt.Println()
}

func sensor_listener(postaje map[string]Postaja, update_channel chan Sporocilo) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	lastMeasurementTime := time.Now()

loop:
	for {
		select {
		case sporocilo := <-update_channel:
			lastMeasurementTime = time.Now()
			rw_lock.Lock()
			for _, m := range sporocilo.meritve {
				dodaj_meritev(postaje, sporocilo.source_id, m)
			}
			rw_lock.Unlock()

		case strike := <-strike_channel:
			rw_lock.Lock()
			postaja, obstaja := postaje[strike.id]
			if obstaja && postaja.is_alive() {
				for j := range 3 {
					hit_or_miss := rand.Intn(2)
					if hit_or_miss == 0 {
						continue
					}
					if postaja.health[j] <= strike.dmg {
						postaja.health[j] = 0
					} else {
						postaja.health[j] -= strike.dmg
					}
				}
				postaje[strike.id] = postaja
			}
			rw_lock.Unlock()

		case <-ticker.C:
			if time.Since(lastMeasurementTime) >= 5*time.Second {
				fmt.Println("\nSistem je neodziven (več kot 5 s ni bilo meritve).")
				break loop
			}
			rw_lock.RLock()
			izpis_stanja_postaj(postaje)
			rw_lock.RUnlock()

		case <-stop_channel:
			fmt.Println("\nProgram ustavljen s pritiskom na Enter.")
			break loop
		}
	}
}

func evil_striker(postaje map[string]Postaja) {
loop:
	for {
		select {
		case <-time.After(400 * time.Millisecond):
			rw_lock.RLock()
			var alive_ids []string
			for id, postaja := range postaje {
				if postaja.is_alive() {
					alive_ids = append(alive_ids, id)
				}
			}
			rw_lock.RUnlock()

			if len(alive_ids) == 0 {
				fmt.Println("\n⚡ EVIL STRIKER: Vse postaje so uničene! Prekinjam napade...")
				break loop
			}

			n_strikes := rand.Intn(7)
			var strikesSummary []string
			for range n_strikes {
				id := alive_ids[rand.Intn(len(alive_ids))]
				dmg := 20 * (rand.Intn(4) + 1)
				strike_channel <- Strike{id: id, dmg: uint8(dmg)}
				strikesSummary = append(strikesSummary, fmt.Sprintf("%s (-%d hp)", id, dmg))
			}
			if len(strikesSummary) > 0 {
				fmt.Printf("⚡ EVIL STRIKER: %s\n", strings.Join(strikesSummary, ", "))
			}

		case <-stop_channel:
			break loop
		}
	}
}

func meritve_generator(id string, postaje map[string]Postaja, update_channel chan Sporocilo) {
loop:
	for {
		select {
		case <-time.After(MERITEV_INTERVAL):
			rw_lock.RLock()
			postaja, obstaja := postaje[id]
			rw_lock.RUnlock()

			if !obstaja || !postaja.is_alive() {
				continue
			}
			var meritve []Meritev
			for i := range postaja.health {
				if postaja.health[i] > 0 {
					meritve = append(meritve, random_meritev(VRSTE_MERITEV[i], postaja.cli[i]))
				}
			}
			if len(meritve) > 0 {
				update_channel <- Sporocilo{
					source_id: id,
					meritve:   meritve,
				}
			}

		case <-stop_channel:
			break loop
		}
	}
}

// func station_reporter(postaja Postaja) {

// }

func main() {

	var postaje = make(map[string]Postaja, 9)
	var postaja Postaja
	var id string
	var n_postaj int

	for i := range postaje_code_names {

		n_postaj = rand.Intn(5) - 1

		for j := 1; j <= n_postaj; j++ {

			id = postaje_code_names[i].slug + formatted_appendix(j)

			postaja = Postaja{
				ime:     postaje_code_names[i].kraj + " " + formatted_appendix(j),
				kraj:    postaje_code_names[i].kraj,
				health:  [3]uint8{uint8((rand.Intn(4) + 1) * 25), uint8((rand.Intn(4) + 1) * 25), uint8((rand.Intn(4) + 1) * 25)},
				cli:     [3]CLI{postaje_code_names[i].temp, postaje_code_names[i].tlak, postaje_code_names[i].vlaga},
				meritve: make([]Meritev, 0),
			}

			pack := try_for_health_pack()
			for k := range postaja.health {
				postaja.health[k] += pack
			}

			postaje[id] = postaja
			fmt.Printf("postaja dodana: [%s] ==> { %s | %v hp } \n", id, postaje[id].ime, postaja.health)

			// 	for k := 0; k < n_meritev; k++ {
			// 		dodaj_meritev(postaje, id, random_meritev("temperatura", postaje_code_names[i].temp))
			// 		dodaj_meritev(postaje, id, random_meritev("tlak", postaje_code_names[i].tlak))
			// 		dodaj_meritev(postaje, id, random_meritev("vlaga", postaje_code_names[i].vlaga))
			// 	}
		}
	}

	update_channel := make(chan Sporocilo, 3)

	go keyboard_listener()
	go evil_striker(postaje)
	for id := range postaje {
		go meritve_generator(id, postaje, update_channel)
	}
	sensor_listener(postaje, update_channel)

	izpis_porocila(postaje)
}

//////////////////////////////////// DN1 stuff:

func formatted_appendix(number int) string {
	return fmt.Sprintf("%02d", number)
}

func random_normal(mean float64, sigma float64) float64 {
	return mean + sigma*rand.NormFloat64()
}

func random_meritev(vrsta string, stat CLI) Meritev {
	return Meritev{
		vrsta:    vrsta,
		vrednost: float32(random_normal(stat.mean, stat.sigma)),
	}
}

func enota_meritve(vrsta string) string {
	switch vrsta {
	case "temperatura":
		return "°C"
	case "tlak":
		return "hPa"
	case "vlaga":
		return "%"
	default:
		return ""
	}
}

func veljavna_meritev(m Meritev) bool {
	val := float64(m.vrednost)
	switch m.vrsta {
	case "temperatura":
		return val >= TEMP_MIN && val <= TEMP_MAX
	case "tlak":
		return val >= TLAK_MIN && val <= TLAK_MAX
	case "vlaga":
		return val >= VLAGA_MIN && val <= VLAGA_MAX
	default:
		return false
	}
}

func dodaj_meritev(postaje map[string]Postaja, id string, meritev Meritev) {
	postaja, obstaja := postaje[id]
	if !obstaja {
		fmt.Printf("Postaja ne obstaja: %s \n", id)
		return
	}

	if !veljavna_meritev(meritev) {
		return
	}

	postaja.meritve = append(postaja.meritve, meritev)
	postaje[id] = postaja
}

func izpis_porocila(postaje map[string]Postaja) {
	rw_lock.RLock()
	defer rw_lock.RUnlock()

	var ids []string
	for id := range postaje {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	fmt.Println("\n=================================== KONČNO STATISTIČNO POROČILO ===================================")
	for _, id := range ids {
		postaja := postaje[id]

		var tempSum, tlakSum, vlagaSum float64
		var tempCnt, tlakCnt, vlagaCnt int

		for _, m := range postaja.meritve {
			switch m.vrsta {
			case "temperatura":
				tempSum += float64(m.vrednost)
				tempCnt++
			case "tlak":
				tlakSum += float64(m.vrednost)
				tlakCnt++
			case "vlaga":
				vlagaSum += float64(m.vrednost)
				vlagaCnt++
			}
		}

		var tempAvg, tlakAvg, vlagaAvg string
		if tempCnt > 0 {
			tempAvg = fmt.Sprintf("%.2f °C", tempSum/float64(tempCnt))
		} else {
			tempAvg = "--"
		}
		if tlakCnt > 0 {
			tlakAvg = fmt.Sprintf("%.2f hPa", tlakSum/float64(tlakCnt))
		} else {
			tlakAvg = "--"
		}
		if vlagaCnt > 0 {
			vlagaAvg = fmt.Sprintf("%.2f %%", vlagaSum/float64(vlagaCnt))
		} else {
			vlagaAvg = "--"
		}

		status := fmt.Sprintf("❤️ HP: %v", postaja.health)
		if !postaja.is_alive() {
			status = "💀 UNIČENA       "
		}

		fmt.Printf("[%s] %-14s | %-18s | N=%-3d | povprečja: %10s | %11s | %8s\n",
			id,
			postaja.kraj,
			status,
			len(postaja.meritve),
			tempAvg,
			tlakAvg,
			vlagaAvg,
		)
	}
	fmt.Println("===================================================================================================")
}
