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

	scraper("0", intCola, intNr, url, intMonos)

}

func getSha256(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func scraper(origen string, tam_cola int, nivel int, url string, mono_id int) {
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

	links := []string{}

	c.OnHTML("p", func(e *colly.HTMLElement) {

		// CONTAR LA CANTIDAD DE PALABRAS
		t := tokenizer.New()
		tokens = tokens + len(t.Tokenize(e.Text))

		// CONTAR LA CANTIDAD DE ENLACES

		links = e.ChildAttrs("a", "href")
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

	fmt.Println(mono)
	c.Visit(url)

}
