package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/euskadi31/go-tokenizer"
	"github.com/gocolly/colly"
)

var Mono struct {
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

	fmt.Println(intMonos)

	cola := getReader("Tamaño de la cola de espera: ")
	intCola, _ := strconv.Atoi(cola)

	fmt.Println(intCola)

	nr := getReader("Nr: ")
	intNr, _ := strconv.Atoi(nr)

	fmt.Println(intNr)

	url := getReader("URL: ")

	fmt.Println(url)

	file := getReader("Nombre de archivo: ")

	fmt.Println(file)

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnHTML("p", func(e *colly.HTMLElement) {

		// CONTAR LA CANTIDAD DE PALABRAS
		t := tokenizer.New()
		tokens := t.Tokenize(e.Text)

		// CONTAR LA CANTIDAD DE ENLACES

		links := e.ChildAttrs("a", "href")
		fmt.Println("{")
		fmt.Println("cantidad de palabras: ", len(tokens))
		fmt.Println("cantidad de enlaces: ", len(links))
		fmt.Println("},")
		fmt.Println(links)

	})

	c.Visit("https://en.wikipedia.org/wiki/Web_scraping")

}
