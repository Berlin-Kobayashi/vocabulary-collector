package main

import (
	"bufio"
	"fmt"
	"github.com/antchfx/xmlquery"
	"os"
	"regexp"
	"strings"
)

func main() {
	parse("/home/ska/lab/go/src/github.com/DanShu93/vocabulary-collector/tmp/Streamlingo/Friends/Season 1/en_us/S01E01")
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func parse(filepath string) {
	file, err := os.Open(filepath)
	if err!= nil {
		panic(err)
	}

	reg := regexp.MustCompile("[^a-zA-Z üöä]")

	reader := bufio.NewReader(file)
	doc, err := xmlquery.Parse(reader)

	list := xmlquery.Find(doc, "//tt/body/div/p")

	counts := make(map[string]int)

	for _, el := range list {
		for _, attr := range el.Attr {
			if attr.Name.Local == "begin" {
				// 34265481249 ÷10000000 ÷ 60
				//fmt.Println(attr.Value)
			} else if attr.Name.Local == "end" {
				//fmt.Println(attr.Value)
			}
		}

		text := el.InnerText()
		res := reg.ReplaceAllString(text, " ")

		cleaned := standardizeSpaces(res)

		words := strings.Fields(cleaned)

		for _, word := range words {

			if _, ok := counts[word] ; !ok {
				counts[word] = 1
			} else {
				counts[word] = counts[word] + 1
			}


		}
		//fmt.Println(cleaned)
	}

	fmt.Println(counts)
}