package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/euskadi31/go-tokenizer"
	"github.com/gocolly/colly"
)

type Mono struct {
	origen          string `json:"origen"`
	conteo_palabras int    `json:"conteo_palabras"`
	conteo_enlaces  int    `json:"conteo_enlaces"`
	sha             string `json:"sha"`
	url             string `json:"url"`
	mono            int    `json:"mono"`
}

func getReader(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)
	t, _ := reader.ReadString('\n')
	t = strings.Replace(t, "\n", "", -1)
	return t
}

func main() {
	colaGeneral := []Mono{}
	monos := getReader("Cantidad de monos buscadores: ")
	intMonos, _ := strconv.Atoi(monos)

	cola := getReader("Tamaño de la cola de espera: ")
	intCola, _ := strconv.Atoi(cola)

	nr := getReader("Nr: ")
	intNr, _ := strconv.Atoi(nr)

	url := getReader("URL: ")

	file := getReader("Nombre de archivo: ")

	//fmt.Println(file)
	f, err := os.Create(file + ".json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Cola")
	colaGeneral = llenarCola(intCola, intNr, url, colaGeneral)
	//fmt.Println(hijos)
	done := make(chan struct{})
	no_mono := 1
	for i := 0; i <= intCola; i = i + (intCola / intMonos) + 1 {
		paso := i + (intCola / intMonos)
		if paso > intCola {
			paso = intCola
		}
		go buscarMono(i, paso, f, done, colaGeneral, no_mono)
		no_mono++
		//scraper("0", intCola, intNr, url, intMonos)
	}

	dones := 0
	for dones < intMonos {
		<-done
		dones++
	}
}

func getSha256(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func buscarMono(inicio int, fin int, f *os.File, done chan struct{}, colaGeneral []Mono, no_mono int) {
	for i := inicio; i < fin; i++ {
		fmt.Print("No. Mono: ")

		monoActual := colaGeneral[i]
		monoActual.mono = no_mono
		fmt.Println(no_mono)
		fmt.Print("URL: ")
		fmt.Println(monoActual.url)
		/*monoSerial := &Mono{
			origen:          monoActual.origen,
			conteo_palabras: monoActual.conteo_palabras,
			conteo_enlaces:  monoActual.conteo_enlaces,
			sha:             monoActual.sha,
			url:             monoActual.url,
			mono:            monoActual.mono,
		}
		b, err := json.Marshal(monoSerial)
		if err != nil {
			fmt.Println(err)
			return
		}*/

		_, error := f.WriteString(fmt.Sprint(monoActual))
		if error != nil {
			panic(error)
		}
		time.Sleep(time.Duration(monoActual.conteo_palabras) * time.Second)
	}

	done <- struct{}{}
}

func llenarCola(intCola, intNr int, url string, colaGeneral []Mono) (cola []Mono) {
	links, cola := scraper("0", intCola, url, 0, colaGeneral)
	for index, link := range links {
		if index >= intNr {
			break
		}
		//fmt.Println("****************otro hijo*****************")
		cola = hijos(intCola, intNr, intNr-1, getSha256(url), link, 0, cola)
		if len(cola) >= intCola {
			break
		}
	}
	return cola
}

func hijos(intCola int, nivel int, nivelActual int, origen string, link string, mono_id int, colaGeneral []Mono) (cola []Mono) {
	//fmt.Println(nivelActual)
	if nivelActual == 0 {
		//fmt.Println("saliendo")
		return colaGeneral
	} else {
		links := []string{}
		links, cola = scraper(origen, intCola, link, mono_id, colaGeneral)
		nivelActual--
		if len(cola) <= intCola {
			for index, element := range links {
				if index >= nivelActual {
					break
				}
				cola = hijos(intCola, nivel, nivelActual, getSha256(link), element, mono_id, cola)
				if len(cola) >= intCola {
					break
				}
			}
			return cola
		} else {
			return cola
		}

	}
	return cola
}

func scraper(origen string, tam_cola int, url string, mono_id int, colaGeneral []Mono) (links []string, cola []Mono) {
	base_url := "https://es.wikipedia.org"

	if !strings.Contains(url, base_url) {
		url = base_url + url
	}

	c := colly.NewCollector(
		colly.AllowedDomains("es.wikipedia.org", "en.wikipedia.org"),
	)
	// CONTAR LA CANTIDAD DE PALABRAS

	tokens := 0
	cantLinks := 0

	// CONTAR LA CANTIDAD DE ENLACES

	c.OnHTML("p", func(e *colly.HTMLElement) {

		// CONTAR LA CANTIDAD DE PALABRAS
		t := tokenizer.New()
		tokens = tokens + len(t.Tokenize(e.Text))

		// CONTAR LA CANTIDAD DE ENLACES

		links = append(links, e.ChildAttrs("a", "href")...)
		//fmt.Println(links)
		cantLinks = cantLinks + len(links)

	})
	mono := Mono{
		origen:          origen,
		conteo_palabras: tokens,
		conteo_enlaces:  cantLinks,
		sha:             getSha256(url),
		url:             url,
		mono:            mono_id,
	}
	if len(colaGeneral) < tam_cola {
		cola = append(colaGeneral, mono)
	}
	fmt.Println(mono)
	c.Visit(url)
	return links, cola
}
