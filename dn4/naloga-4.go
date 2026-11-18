package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const n_kupcev = 1000000
const n_obdelovalcev = 2000
const mean_cas_obdelave = 0.0
const nakupovalni_cas = 100.0
const kupec_hp = 70
// const n_kupcev = 5
// const n_obdelovalcev = 2
// const mean_cas_obdelave = 500.0
// const nakupovalni_cas = 250.0
// const kupec_hp = 9


type narocilo interface {
	obdelaj(id_obdelovalca int) error
}

type izdelek struct {
	imeIzdelka string
	cena       float64
	teza       float64
}

func (i izdelek) izpis() string {
	return fmt.Sprintf("IZDELEK: Cena: %.2f €, Teža: %.2f kg | Ime: %s", i.cena, i.teza, i.imeIzdelka)
}

type eknjiga struct {
	naslovKnjige string
	cena         float64
}

func (e eknjiga) izpis() string {
	return fmt.Sprintf("EKNJIGA: Cena: %.2f € | Naslov: %s", e.cena, e.naslovKnjige)
}

type spletniTecaj struct {
	imeTecaja   string
	trajanjeUre int
	cenaUre     float64
}

func (t spletniTecaj) izpis() string {
	return fmt.Sprintf("TECAJ: Trajanje: %3d h, Cena/h: %.2f -> skupaj: %.2f | Ime: %s", t.trajanjeUre, t.cenaUre, t.cenaUre*float64(t.trajanjeUre), t.imeTecaja)
}

func random_normal(mean float64, sigma float64) float64 {
	return mean + sigma*rand.NormFloat64()
}

func gen_izdelek() izdelek {
	izdelek := izdelek{
		imeIzdelka: imena_izdelkov[rand.Intn(len(imena_izdelkov))],
		cena:       random_normal(10.0, 5.0),
		teza:       random_normal(1.0, 0.5),
	}

	return izdelek
}

func gen_eknjigo() eknjiga {
	eknjiga := eknjiga{
		naslovKnjige: imena_knjig[rand.Intn(len(imena_knjig))],
		cena:         random_normal(10.0, 5.0),
	}

	return eknjiga
}

func gen_tecaj() spletniTecaj {
	spletniTecaj := spletniTecaj{
		imeTecaja:   imena_tecajev[rand.Intn(len(imena_tecajev))],
		trajanjeUre: int(random_normal(12, 6)) * 5,
		cenaUre:     random_normal(30, 15),
	}

	return spletniTecaj
}

func (izdelek *izdelek) obdelaj(id_obdelovalca int) error {
	fmt.Printf("[%d] Obdelujem...\n", id_obdelovalca)
	cas_obdelave := time.Duration(random_normal(mean_cas_obdelave, mean_cas_obdelave/2))
	if cas_obdelave < 0 {
		cas_obdelave = 1
	}
	time.Sleep(cas_obdelave)

	if len(izdelek.imeIzdelka) == 0 {
		return fmt.Errorf("izdelek nima imena [%s, %.2f, %.2f]", izdelek.imeIzdelka, izdelek.cena, izdelek.teza)
	}
	if izdelek.cena <= 0 {
		return fmt.Errorf("izdelek nima veljavne cene [%s, %.2f, %.2f]", izdelek.imeIzdelka, izdelek.cena, izdelek.teza)
	}
	if izdelek.teza <= 0 {
		return fmt.Errorf("izdelek nima veljavne teže [%s, %.2f, %.2f]", izdelek.imeIzdelka, izdelek.cena, izdelek.teza)
	}

	rlock.Lock()
	promet += izdelek.cena
	st_narocil++
	fmt.Printf("[%d] Obdelal narocilo [%d] v %d ms: %s\n", id_obdelovalca, st_narocil, cas_obdelave, izdelek.izpis())
	rlock.Unlock()
	return nil
}

func (eknjiga *eknjiga) obdelaj(id_obdelovalca int) error {
	fmt.Printf("[%d] Obdelujem...\n", id_obdelovalca)
	cas_obdelave := time.Duration(random_normal(mean_cas_obdelave, mean_cas_obdelave/2))
	if cas_obdelave < 0 {
		cas_obdelave = 1
	}
	time.Sleep(cas_obdelave)

	if len(eknjiga.naslovKnjige) == 0 {
		return fmt.Errorf("eknjiga nima imena [%s, %.2f]", eknjiga.naslovKnjige, eknjiga.cena)
	}
	if eknjiga.cena <= 0 {
		return fmt.Errorf("eknjiga nima veljavne cene [%s, %.2f]", eknjiga.naslovKnjige, eknjiga.cena)
	}

	rlock.Lock()
	promet += eknjiga.cena
	st_narocil++
	fmt.Printf("[%d] Obdelal narocilo [%d] v %d ms: %s\n", id_obdelovalca, st_narocil, cas_obdelave, eknjiga.izpis())
	rlock.Unlock()
	return nil
}

func (spletniTecaj *spletniTecaj) obdelaj(id_obdelovalca int) error {
	fmt.Printf("[%d] Obdelujem...\n", id_obdelovalca)
	cas_obdelave := time.Duration(random_normal(mean_cas_obdelave, mean_cas_obdelave/2))
	if cas_obdelave < 0 {
		cas_obdelave = 1
	}
	time.Sleep(cas_obdelave)

	if len(spletniTecaj.imeTecaja) == 0 {
		return fmt.Errorf("spletni tecaj nima imena [%s, %.2f, %v]", spletniTecaj.imeTecaja, spletniTecaj.cenaUre, spletniTecaj.trajanjeUre)
	}
	if spletniTecaj.cenaUre <= 0 {
		return fmt.Errorf("spletni tecaj nima veljavne cene [%s, %.2f, %v]", spletniTecaj.imeTecaja, spletniTecaj.cenaUre, spletniTecaj.trajanjeUre)
	}
	if spletniTecaj.trajanjeUre <= 0 {
		return fmt.Errorf("spletni tecaj nima veljavnega trajanja [%s, %.2f, %v]", spletniTecaj.imeTecaja, spletniTecaj.cenaUre, spletniTecaj.trajanjeUre)
	}

	rlock.Lock()
	promet += spletniTecaj.cenaUre * float64(spletniTecaj.trajanjeUre)
	st_narocil++
	fmt.Printf("[%d] Obdelal narocilo [%d] v %d ms: %s\n", id_obdelovalca, st_narocil, cas_obdelave, spletniTecaj.izpis())
	rlock.Unlock()
	return nil
}

func kupec(id int, narocila chan narocilo, wg *sync.WaitGroup, stop_channel chan any) {
	defer wg.Done()
	randn := int(rand.Intn(kupec_hp))
	randn2 := 0

	for randn != 0 {
		fmt.Printf("[%d]Brskam...\n", id)
		select {
		case <-stop_channel:
			return
		case <-time.After(time.Millisecond * time.Duration(nakupovalni_cas*rand.Intn(5))):
		}

		randn2 = int(rand.Intn(2))
		if randn2 == 0 {
			select {
			case <-stop_channel:
				return
			case <-time.After(time.Millisecond * time.Duration(nakupovalni_cas*rand.Intn(3))):
				izdelek := gen_izdelek()
				fmt.Printf("[%d]👀 %s\n", id, izdelek.izpis())
				narocila <- &izdelek
			}
		}
		randn2 = int(rand.Intn(2))
		if randn2 == 0 {
			select {
			case <-stop_channel:
				return
			case <-time.After(time.Millisecond * time.Duration(nakupovalni_cas*rand.Intn(5))):
				eknjiga := gen_eknjigo()
				fmt.Printf("[%d]👀 %s\n", id, eknjiga.izpis())
				narocila <- &eknjiga
			}
		}
		randn2 = int(rand.Intn(2))
		if randn2 == 0 {
			select {
			case <-stop_channel:
				return
			case <-time.After(time.Millisecond * time.Duration(nakupovalni_cas*rand.Intn(5))):
				tecaj := gen_tecaj()
				fmt.Printf("[%d]👀 %s\n", id, tecaj.izpis())
				narocila <- &tecaj
			}
		}

		randn = int(rand.Intn(kupec_hp))
	}
}

func obdelovalec(id int, napake chan error, narocila chan narocilo, wg *sync.WaitGroup) {
	defer wg.Done()
	for narocilo := range narocila {
		if err := narocilo.obdelaj(id); err != nil {
			napake <- err
		}
	}
}

func catch_errors(napake chan error, error_messages *[]error) {
	defer wg.Done()

	for err := range napake {
		rlock.Lock()
		*error_messages = append(*error_messages, err)
		rlock.Unlock()
	}
}

func print_error_messages(error_messages []error) {
	fmt.Printf("\nStevilo napak: %d, razmerje: %.2f%%\n\n", len(error_messages), float64(len(error_messages))/float64(st_narocil)*100)

	for i, err := range error_messages {
		fmt.Printf("%d. %s\n", i+1, err)
	}
}

func keyboard_listener(stop_channel chan any) {
	fmt.Scanln()
	close(stop_channel)
}

var rlock sync.RWMutex
var promet float64
var st_narocil int
var wg sync.WaitGroup

func main() {
	napake := make(chan error, 10)
	narocila := make(chan narocilo, 10)
	stop_channel := make(chan any)

	promet = 0
	st_narocil = 0
	var error_messages []error

	var wg_kupci sync.WaitGroup
	for i := range n_kupcev {
		wg_kupci.Add(1)
		go kupec(i, narocila, &wg_kupci, stop_channel)
	}

	var wg_obdelovalci sync.WaitGroup
	for i := range n_obdelovalcev {
		wg_obdelovalci.Add(1)
		go obdelovalec(i, napake, narocila, &wg_obdelovalci)
	}

	go keyboard_listener(stop_channel)

	wg.Add(1)
	go catch_errors(napake, &error_messages)

	wg_kupci.Wait()
	close(narocila)
	wg_obdelovalci.Wait()
	close(napake)
	wg.Wait()

	fmt.Printf("\nKupci so odšli. Statistika: stevilo narocil: %d, promet: %.2f\n\n", st_narocil, promet)
	fmt.Println("Izpisujem napake:")
	print_error_messages(error_messages)
}

var imena_izdelkov = []string{
	"banana",
	"jabolko",
	"hruška",
	"bananina",
	"čebula",
	"paradižnik",
	"pes",
	"mačka",
	"avto",
	"kolo",
	"miza",
	"stol",
	"lampa",
	"računalnik",
	"telefon",
	"ura",
	"sveča",
	"zvezek",
	"rokometna žoga",
	"košarica",
	"",
}

var imena_knjig = []string{
	"Kapital - Karl Marx",
	"Dekameron - Giovanni Boccaccio",
	"Cvetje v jeseni - Ivan Tavčar",
	"Visoška kronika - Ivan Tavčar",
	"Na klancu - Ivan Cankar",
	"Kralj na Betajnovi - Ivan Cankar",
	"Hlapci - Ivan Cankar",
	"Pod svobodnim soncem - Fran Saleški Finžgar",
	"Deseti brat - Josip Jurčič",
	"Krst pri Savici - France Prešeren",
	"Martin Krpan z Vrha - Fran Levstik",
	"Butalci - Fran Milčinski",
	"Samorastniki - Prežihov Voranc",
	"Alamut - Vladimir Bartol",
	"Pomladni dan - Ciril Kosmač",
	"Nekropola - Boris Pahor",
	"Zločin in kazen - Fjodor Dostojevski",
	"Vojna in mir - Lev Tolstoj",
	"Preobrazba - Franz Kafka",
	"Tujec - Albert Camus",
	"",
}

var imena_tecajev = []string{
	"Planinarjenje",
	"Kolesarjenje",
	"Plavanje",
	"Ribanje",
	"Plezanje",
	"Tekanje",
	"Nordijska hoja",
	"Jahanje",
	"Smučanje",
	"Deskanje na vodi",
	"Tenis",
	"Košarka",
	"Odbojka",
	"Rokomet",
	"Pikado",
	"Namizni tenis",
	"Boks",
	"Taekwondo",
	"Joga",
	"Pilates",
	"",
}
