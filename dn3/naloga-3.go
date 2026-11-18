package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type narocilo interface {
	obdelaj() error
}

type izdelek struct {
	imeIzdelka string
	cena       float64
	teza       float64
}

func (i izdelek) izpisi() {
	fmt.Printf("IZDELEK:  Cena: %6.2f €, Teža: %5.2f kg | Ime: %s\n", i.cena, i.teza, i.imeIzdelka)
}

type eknjiga struct {
	naslovKnjige string
	cena         float64
}

func (e eknjiga) izpisi() {
	fmt.Printf("EKNJIGA:  Cena: %3.2f € | Naslov: %s\n", e.cena, e.naslovKnjige)
}

type spletniTecaj struct {
	imeTecaja   string
	trajanjeUre int
	cenaUre     float64
}

func (t spletniTecaj) izpisi() {
	fmt.Printf("TECAJ:    Trajanje: %3d h, Cena/h: %6.2f € | Ime: %s\n", t.trajanjeUre, t.cenaUre, t.imeTecaja)
}



func random_normal(mean float64, sigma float64) float64 {
	return mean + sigma*rand.NormFloat64()
}

func gen_izdelek() izdelek {
	izdelek := izdelek{
		imeIzdelka: imena_izdelkov[rand.Intn(len(imena_izdelkov))],
		cena:       random_normal(10.0, 10.0),
		teza:       random_normal(1.0, 1.0),
	}

	return izdelek
}

func gen_eknjigo() eknjiga {
	eknjiga := eknjiga{
		naslovKnjige: imena_knjig[rand.Intn(len(imena_knjig))],
		cena:         random_normal(10.0, 10.0),
	}

	return eknjiga
}


func gen_tecaj() spletniTecaj {
	spletniTecaj := spletniTecaj{
		imeTecaja:   imena_tecajev[rand.Intn(len(imena_tecajev))],
		trajanjeUre: int(random_normal(12, 7)) * 5,
		cenaUre:     random_normal(30, 30),
	}

	return spletniTecaj
}

func (izdelek *izdelek) obdelaj(id int, napake chan error) error {
	defer wg.Done()
	if len(izdelek.imeIzdelka) == 0 {
		napake <- fmt.Errorf("izdelek nima imena [%s, %.2f, %.2f]", izdelek.imeIzdelka, izdelek.cena, izdelek.teza)
		return nil
	}
	if izdelek.cena <= 0 {
		napake <- fmt.Errorf("izdelek nima veljavne cene [%s, %.2f, %.2f]", izdelek.imeIzdelka, izdelek.cena, izdelek.teza)
		return nil
	}
	if izdelek.teza <= 0 {
		napake <- fmt.Errorf("izdelek nima veljavne teže [%s, %.2f, %.2f]", izdelek.imeIzdelka, izdelek.cena, izdelek.teza)
		return nil
	}

	rlock.Lock()
	promet += izdelek.cena
	st_narocil++
	fmt.Printf("Veljaven izdelek: Cena: %6.2f €, Teža: %5.2f kg | Ime: %s\n", izdelek.cena, izdelek.teza, izdelek.imeIzdelka)
	fmt.Printf("[%d]Narocilo [%d]: Cena: %.2f\n", id, st_narocil, izdelek.cena)
	rlock.Unlock()
	return nil
}

func (eknjiga *eknjiga) obdelaj(id int, napake chan error) error {
	defer wg.Done()
	if len(eknjiga.naslovKnjige) == 0 {
		napake <- fmt.Errorf("eknjiga nima imena [%s, %.2f]", eknjiga.naslovKnjige, eknjiga.cena)
		return nil
	}
	if eknjiga.cena <= 0 {
		napake <- fmt.Errorf("eknjiga nima veljavne cene [%s, %.2f]", eknjiga.naslovKnjige, eknjiga.cena)
		return nil
	}

	rlock.Lock()
	promet += eknjiga.cena
	st_narocil++
	fmt.Printf("Veljavna eknjiga: Cena: %6.2f € | Ime: %s\n", eknjiga.cena, eknjiga.naslovKnjige)
	fmt.Printf("[%d]Narocilo [%d]: Cena: %.2f\n", id, st_narocil, eknjiga.cena)
	rlock.Unlock()
	return nil
}

func (spletniTecaj *spletniTecaj) obdelaj(id int, napake chan error) error {
	defer wg.Done()

	if len(spletniTecaj.imeTecaja) == 0 {
		napake <- fmt.Errorf("spletni tecaj nima imena [%s, %.2f, %v]", spletniTecaj.imeTecaja, spletniTecaj.cenaUre, spletniTecaj.trajanjeUre)
		return nil
	}
	if spletniTecaj.cenaUre <= 0 {
		napake <- fmt.Errorf("spletni tecaj nima veljavne cene [%s, %.2f, %v]", spletniTecaj.imeTecaja, spletniTecaj.cenaUre, spletniTecaj.trajanjeUre)
		return nil
	}
	if spletniTecaj.trajanjeUre <= 0 {
		napake <- fmt.Errorf("spletni tecaj nima veljavnega trajanja [%s, %.2f, %v]", spletniTecaj.imeTecaja, spletniTecaj.cenaUre, spletniTecaj.trajanjeUre)
		return nil
	}

	rlock.Lock()
	promet += spletniTecaj.cenaUre * float64(spletniTecaj.trajanjeUre)
	st_narocil++
	fmt.Printf("Veljaven spletni tečaj: Cena: %6.2f €, Trajanje: %v ur | Ime: %s\n", spletniTecaj.cenaUre * float64(spletniTecaj.trajanjeUre), spletniTecaj.trajanjeUre, spletniTecaj.imeTecaja)
	fmt.Printf("[%d]Narocilo [%d]: Cena: %.2f\n", id, st_narocil, spletniTecaj.cenaUre * float64(spletniTecaj.trajanjeUre))
	rlock.Unlock()
	return nil
}


func kupec(id int, napake chan error)  {
	defer wg.Done()
	randn := int(rand.Intn(9))
	randn2 := 0

	for randn != 0 {
		fmt.Printf("[%d]Brskam...\n", id)
		time.Sleep(time.Millisecond * time.Duration(100 * rand.Intn(5)))

		randn2 = int(rand.Intn(2))
		if (randn2 == 0) {
			time.Sleep(time.Millisecond * time.Duration(100 * rand.Intn(3)))
			izdelek := gen_izdelek()
			wg.Add(1)
			go izdelek.obdelaj(id, napake)
		}
		randn2 = int(rand.Intn(2))
		if (randn2 == 0) {
			time.Sleep(time.Millisecond * time.Duration(100 * rand.Intn(5)))
			eknjiga := gen_eknjigo()
			wg.Add(1)
			go eknjiga.obdelaj(id, napake)
		}
		randn2 = int(rand.Intn(2))
		if (randn2 == 0) {
			time.Sleep(time.Millisecond * time.Duration(100 * rand.Intn(5)))
			tecaj := gen_tecaj()
			wg.Add(1)
			go tecaj.obdelaj(id, napake)
		}


		randn = int(rand.Intn(9))
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
	for _, err := range error_messages {
		fmt.Println(err)
	}
}

var rlock sync.RWMutex
var promet float64
var st_narocil int
var wg sync.WaitGroup

func main() {
	napake := make(chan error, 10)
	promet = 0
	st_narocil = 0
	var error_messages []error

	for i := range 10 {
		wg.Add(1)
		go kupec(i, napake)
	}

	go catch_errors(napake, &error_messages)
	wg.Wait()
	wg.Add(1)
	close(napake)
	fmt.Printf("\nKupci so odšli. Statistika: stevilo narocil: %d, promet: %.2f\n\n", st_narocil, promet)
	fmt.Println("Izpisujem napake:")
	wg.Wait()
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
	"rokomet",
	"košarica",
	"",
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
	"",
}
