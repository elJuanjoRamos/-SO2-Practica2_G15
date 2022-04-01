package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/euskadi31/go-tokenizer"
	"github.com/gocolly/colly"
)

type Mono struct {
	origen          string
	conteo_palabras int
	conteo_enlaces  int
	sha             string
	url             string
	mono            int
}

func getReader(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)
	t, _ := reader.ReadString('\n')
	t = strings.Replace(t, "\n", "", -1)
	return t
}

func main() {

	monos := getReader("Cantidad de monos buscadores: ")
	intMonos, _ := strconv.Atoi(monos)

	cola := getReader("Tamaño de la cola de espera: ")
	intCola, _ := strconv.Atoi(cola)

	nr := getReader("Nr: ")
	intNr, _ := strconv.Atoi(nr)

	url := getReader("URL: ")

	file := getReader("Nombre de archivo: ")

	fmt.Println(file)
	f, err := os.Create(file + ".txt")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})

	for i := 0; i <= intCola; i = i + (intCola / intMonos) + 1 {
		paso := i + (intCola / intMonos)
		if paso > intCola {
			paso = intCola
		}
		go recorrerHijos(url, intMonos, intNr, f, done, paso, i)
		//scraper("0", intCola, intNr, url, intMonos)
	}

	dones := 0
	for dones < intMonos {
		<-done
		dones++
	}
}

func recorrerHijos(url string, mono_id int, nivel int, f *os.File, done chan struct{}, paso int, inicio int) {
	for i := inicio; i <= paso; i++ {
		links := scraper("0", url, mono_id, f)
		for _, link := range links {
			hijos(nivel, 1, url, link, mono_id, f)
		}
	}
	done <- struct{}{}
}

func hijos(nivel int, nivelActual int, origen string, link string, mono_id int, f *os.File) {
	if nivelActual == nivel {
		return
	} else {
		links := scraper(origen, link, mono_id, f)
		for _, link := range links {
			hijos(nivel, nivelActual, origen, link, mono_id, f)
		}
		nivelActual++
	}
}

func getSha256(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func scraper(origen string, url string, mono_id int, f *os.File) (links []string) {
	base_url := "https://es.wikipedia.org"

	if !strings.Contains(url, base_url) {
		url = base_url + url
	}

	c := colly.NewCollector(
		colly.AllowedDomains("es.wikipedia.org", "en.wikipedia.org"),
	)

	c.OnHTML("p", func(e *colly.HTMLElement) {

		// CONTAR LA CANTIDAD DE PALABRAS
		t := tokenizer.New()
		tokens := t.Tokenize(e.Text)

		// CONTAR LA CANTIDAD DE ENLACES

		links = e.ChildAttrs("a", "href")

		mono := Mono{
			origen:          origen,
			conteo_palabras: len(tokens),
			conteo_enlaces:  len(links),
			sha:             getSha256(url),
			url:             url,
			mono:            mono_id,
		}
		_, err := f.WriteString(fmt.Sprint(mono))
		if err != nil {
			panic(err)
		}
		fmt.Println(mono)
	})

	c.Visit(url)
	return links
}
