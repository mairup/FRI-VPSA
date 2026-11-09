package main

import (
	"fmt"
	rand "math/rand"
)

const TEMP_MIN = -25.0
const TEMP_MAX = 50.0

type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

type temperature struct {
	kraj_slug string
	mean      float64
	sigma     float64
}

// Podatki o temperaturah (srednja vrednost in standardni odklon) so izračunani iz 7-dnevnih urnih meritev
// prek Open-Meteo API (https://open-meteo.com/):
// Primer za Ljubljano: https://api.open-meteo.com/v1/forecast?latitude=46.0569&longitude=14.5058&past_days=7&forecast_days=0&hourly=temperature_2m
// Format poizvedbe: https://api.open-meteo.com/v1/forecast?latitude={lat}&longitude={lon}&past_days=7&forecast_days=0&hourly=temperature_2m
// Datum pridobitve podatkov: 27. 11. 2026

var postaje_code_names = []struct {
	slug  string
	kraj  string
	count int
	temp  temperature
}{
	{"LJ", "Ljubljana", 0, temperature{"LJ", 9.9, 3.6}},
	{"KR", "Kranj", 0, temperature{"KR", 9.0, 3.6}},
	{"NM", "Novo mesto", 0, temperature{"NM", 9.6, 3.8}},
	{"KP", "Koper", 0, temperature{"KP", 14.2, 4.3}},
	{"CE", "Celje", 0, temperature{"CE", 9.5, 4.1}},
	{"GO", "Goriška Brda", 0, temperature{"GO", 12.7, 3.9}},
	{"MB", "Maribor", 0, temperature{"MB", 10.4, 4.0}},
	{"SG", "Slovenj Gradec", 0, temperature{"SG", 8.9, 4.3}},
	{"MS", "Murska Sobota", 0, temperature{"MS", 9.6, 4.3}},
}

func formatted_appendix(number int) string {
	return fmt.Sprintf("%02d", number)
}

func random_normal(mean float64, sigma float64) float64 {
	return mean + sigma*rand.NormFloat64()
}

func dodaj_meritev(postaje map[string]Postaja, id string, meritev float64) {
	_, obstaja := postaje[id]
	if !obstaja {
		fmt.Printf("Postaja ne obstaja: %s \n", id)
		return
	}

	if meritev > TEMP_MAX || meritev < TEMP_MIN {
		fmt.Printf("Neveljavna meritev: %.1f \n", meritev)
		return
	}

	postaja := postaje[id]
	postaja.meritve = append(postaja.meritve, meritev)
	postaje[id] = postaja
}

func izpis_meritev(postaje map[string]Postaja) {
	for id, postaja := range postaje {
		fmt.Printf("[%s] %s [ ", id, postaja.kraj)
		for _, meritev := range postaja.meritve {
			fmt.Printf("%.1f ", meritev)
		}
		fmt.Printf("]\n")
	}
}

func statistika(postaje map[string]Postaja, id string) (min, max, mean float64, ok bool) {
	postaja, obstaja := postaje[id]
	if !obstaja {
		return 0, 0, 0, false
	}

	if len(postaja.meritve) < 3 {
		return 0, 0, 0, false
	}

	min = postaja.meritve[0]
	max = postaja.meritve[0]
	mean = 0

	for _, meritev := range postaja.meritve {
		if meritev < min {
			min = meritev
		}
		if meritev > max {
			max = meritev
		}
		mean += meritev
	}

	mean /= float64(len(postaja.meritve))
	ok = true
	return
}

func ocena_temperature (temp float64) string {

	switch {
	case temp >= 25:
		return "vroče"
	case temp >= 15:
		return "toplo"
	case temp >= 5:
		return "hladno"
	default:
		return "mrzlo"
	}

}

func izpis_porocila(postaje map[string]Postaja) {

	fmt.Printf("POROČILO STATISTIKE:\n\n")

	for id, postaja := range postaje {
		min, max, mean, ok := statistika(postaje, id)
		if !ok {
			fmt.Printf("[%s] %s - premalo meritev\n", id, postaja.kraj)
		} else {
			fmt.Printf("[%s] %s min=%.1f max=%.1f avg=%.1f --> %s\n", id, postaja.kraj, min, max, mean, ocena_temperature(mean))
		}

	}

}

func main() {

	var postaje = make(map[string]Postaja, 9)
	var postaja Postaja
	var id string
	var n_postaj int
	var n_meritev int

	for i := range postaje_code_names {

		n_postaj = rand.Intn(3)

		for j := 1; j <= n_postaj; j++ {

			id = postaje_code_names[i].slug + formatted_appendix(j)

			postaja = Postaja{
				ime:     postaje_code_names[i].kraj + " " + formatted_appendix(j),
				kraj:    postaje_code_names[i].kraj,
				meritve: make([]float64, 0),
			}
			n_meritev = rand.Intn(6) + 1

			for k := 0; k < n_meritev; k++ {
				postaja.meritve = append(
					postaja.meritve,
					random_normal(
						postaje_code_names[i].temp.mean,
						postaje_code_names[i].temp.sigma,
					),
				)
			}

			postaje[id] = postaja
			fmt.Printf("postaja dodana: [%s] ==> { %s } [ ", id, postaje[id].ime)
			for _, meritev := range postaja.meritve {
				fmt.Printf("%.1f ", meritev)
			}
			fmt.Printf("]\n")
		}
	}

	fmt.Printf("\nIZPIS VSEH TRENUTNIH POSTAJ IN MERITEV:\n")
	izpis_meritev(postaje)
	fmt.Printf("\n======================\n\n")

	izpis_porocila(postaje)
}
