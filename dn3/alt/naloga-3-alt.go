package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type narocilo interface {
	obdelaj() error
}

type izdelek struct {
	imeIzdelka string
	cena       float64
	teza       float64
}

type eknjiga struct {
	naslovKnjige string
	cena         float64
}

type spletniTecaj struct {
	imeTecaja   string
	trajanjeUre int
	cenaUre     float64
}

var rlock sync.RWMutex
var promet float64
var st_narocil int

func posodobi_promet(vrednost float64) {
	rlock.Lock()
	promet += vrednost
	st_narocil++
	rlock.Unlock()
}

func (i *izdelek) obdelaj() error {
	if len(i.imeIzdelka) == 0 {
		return fmt.Errorf("izdelek nima imena [%s, %.2f, %v]", i.imeIzdelka, i.cena, i.teza)
	}
	if i.cena <= 0 {
		return fmt.Errorf("izdelek nima veljavne cene [%s, %.2f, %v]", i.imeIzdelka, i.cena, i.teza)
	}
	if i.teza <= 0 {
		return fmt.Errorf("izdelek nima veljavne teže [%s, %.2f, %v]", i.imeIzdelka, i.cena, i.teza)
	}

	fmt.Printf("Veljaven izdelek: Cena: %6.2f €, Teža: %5.2f kg | Ime: %s\n", i.cena, i.teza, i.imeIzdelka)
	posodobi_promet(i.cena)
	return nil
}

func (e *eknjiga) obdelaj() error {
	if len(e.naslovKnjige) == 0 {
		return fmt.Errorf("eknjiga nima imena [%s, %.2f]", e.naslovKnjige, e.cena)
	}
	if e.cena <= 0 {
		return fmt.Errorf("eknjiga nima veljavne cene [%s, %.2f]", e.naslovKnjige, e.cena)
	}

	fmt.Printf("Veljavna eknjiga: Cena: %6.2f € | Ime: %s\n", e.cena, e.naslovKnjige)
	posodobi_promet(e.cena)
	return nil
}

func (t *spletniTecaj) obdelaj() error {
	if len(t.imeTecaja) == 0 {
		return fmt.Errorf("spletni tečaj nima imena [%s, %.2f, %v]", t.imeTecaja, t.cenaUre, t.trajanjeUre)
	}
	if t.cenaUre <= 0 {
		return fmt.Errorf("spletni tečaj nima veljavne cene [%s, %.2f, %v]", t.imeTecaja, t.cenaUre, t.trajanjeUre)
	}
	if t.trajanjeUre <= 0 {
		return fmt.Errorf("spletni tečaj nima veljavnega trajanja [%s, %.2f, %v]", t.imeTecaja, t.cenaUre, t.trajanjeUre)
	}

	vrednost := t.cenaUre * float64(t.trajanjeUre)
	fmt.Printf("Veljaven spletni tečaj: Cena: %6.2f € | Ime: %s\n", vrednost, t.imeTecaja)
	posodobi_promet(vrednost)
	return nil
}

func random_normal(mean float64, sigma float64) float64 {
	return mean + sigma*rand.NormFloat64()
}

func gen_izdelek() izdelek {
	return izdelek{
		imeIzdelka: imena_izdelkov[rand.Intn(len(imena_izdelkov))],
		cena:       random_normal(10.0, 10.0),
		teza:       random_normal(1.0, 1.0),
	}
}

func gen_eknjigo() eknjiga {
	return eknjiga{
		naslovKnjige: imena_knjig[rand.Intn(len(imena_knjig))],
		cena:         random_normal(10.0, 10.0),
	}
}

func gen_tecaj() spletniTecaj {
	return spletniTecaj{
		imeTecaja:   imena_tecajev[rand.Intn(len(imena_tecajev))],
		trajanjeUre: int(random_normal(12, 7)) * 5,
		cenaUre:     random_normal(30, 30),
	}
}

func gen_narocilo() narocilo {
	switch rand.Intn(3) {
	case 0:
		i := gen_izdelek()
		return &i
	case 1:
		e := gen_eknjigo()
		return &e
	default:
		t := gen_tecaj()
		return &t
	}
}

func main() {
	narocila := make([]narocilo, 0, 12)
	for range 10 {
		narocila = append(narocila, gen_narocilo())
	}
	// vsaj dve naročili sta zagotovo neveljavni (naloga to zahteva)
	narocila = append(narocila, &izdelek{cena: 19.99, teza: 1.2})
	narocila = append(narocila, &eknjiga{naslovKnjige: "Brez naslova", cena: -4.20})

	napake := make(chan error, len(narocila))
	var wg sync.WaitGroup

	for _, n := range narocila {
		wg.Add(1)
		go func(n narocilo) {
			defer wg.Done()
			napake <- n.obdelaj()
		}(n)
	}

	wg.Wait()
	close(napake)

	st_napak := 0
	for err := range napake {
		if err != nil {
			st_napak++
			fmt.Println(err)
		}
	}

	fmt.Printf("\nObdelanih naročil: %d, promet: %.2f €, napak: %d\n", st_narocil, promet, st_napak)
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
	"Zločin in kazan - Fjodor Dostojevski",
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
