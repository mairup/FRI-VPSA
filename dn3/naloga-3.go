package main

import (
	"fmt"
	"math/rand"
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

func (izdelek *izdelek) obdelaj() error {
	if len(izdelek.imeIzdelka) == 0 {
		return fmt.Errorf("izdelek nima imena")
	}
	if izdelek.cena <= 0 {
		return fmt.Errorf("izdelek nima veljavne cene")
	}
	if izdelek.teza <= 0 {
		return fmt.Errorf("izdelek nima veljavne teže")
	}

	fmt.Printf("Veljaven izdelek: ")
	izdelek.izpisi()
	return nil
}

func (eknjiga *eknjiga) obdelaj() error {
	if len(eknjiga.naslovKnjige) == 0 {
		return fmt.Errorf("eknjiga nima imena")
	}
	if eknjiga.cena <= 0 {
		return fmt.Errorf("eknjiga nima veljavne cene")
	}
	return nil
}

func (spletniTecaj *spletniTecaj) obdelaj() error {
	if len(spletniTecaj.imeTecaja) == 0 {
		return fmt.Errorf("spletni tečaj nima imena")
	}
	if spletniTecaj.cenaUre <= 0 {
		return fmt.Errorf("spletni tečaj nima veljavne cene")
	}
	if spletniTecaj.trajanjeUre <= 0 {
		return fmt.Errorf("spletni tečaj nima veljavnega trajanja")
	}

	return nil
}

func kupec()  {
	randn := int(rand.Intn(9))
	randn2 := 0

	for randn != 0 {
		randn2 = int(rand.Intn(2))
		if (randn2 == 0) {
			gen_izdelek()
		}
		randn2 = int(rand.Intn(2))
		if (randn2 == 0) {
			gen_eknjigo()
		}
		randn2 = int(rand.Intn(2))
		if (randn2 == 0) {
			gen_tecaj()
		}

		randn := int(rand.Intn(9))
	}
}

func kupi(){

}

func main() {



	// for range 7 {
	// 	izdelek := gen_izdelek()
	// 	izdelek.izpisi()
	// }

	// for range 7 {
	// 	eknjiga := gen_eknjigo()
	// 	eknjiga.izpisi()
	// }

	// for range 7 {
	// 	tecaj := gen_tecaj()
	// 	tecaj.izpisi()
	// }

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
