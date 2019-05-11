package main

import (
	"bufio"
	"fmt"
	"github.com/antchfx/xmlquery"
	"os"
	"sort"
	"strings"
	"unicode"
)

func main() {
	results := parse("tmp/Streamlingo/Friends/Season 1/en_us/S01E01", 30, 4)

	for _, el := range results {
		fmt.Printf("%d %s %s %s\n", el.Count, el.Word, el.Begin[0], el.End[0])
	}
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

type Cnt struct {
	Count int
	Begin []string
	End   []string
	Word string
}

func parse(filepath string, results int, minLength int) []*Cnt {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	doc, err := xmlquery.Parse(reader)

	list := xmlquery.Find(doc, "//tt/body/div/p")

	counts := make(map[string]*Cnt)

	for _, el := range list {
		var begin string
		var end string
		for _, attr := range el.Attr {
			if attr.Name.Local == "begin" {
				begin = attr.Value
				// 34265481249 ÷10000000 ÷ 60
			} else if attr.Name.Local == "end" {
				end = attr.Value
			}
		}

		text := el.InnerText()

		words := strings.Fields(text)

		for _, word := range words {
			if ! IsLetter(word) {
				continue
			}

			if len(word) < minLength {
				continue
			}

			if _, ok := counts[word]; !ok {
				counts[word] = &Cnt{
					Count: 1,
					Begin: []string{begin},
					End:   []string{end},
					Word: word,
				}
			} else {
				counts[word].Count = counts[word].Count + 1
				counts[word].Begin = append(counts[word].Begin, begin)
				counts[word].End = append(counts[word].End, end)
			}
		}
	}

	ranks := make([]*Cnt, 0, len(counts))

	for key := range counts {
		ranks = append(ranks, counts[key])
	}

	sort.Slice(ranks[:], func(i, j int) bool {
		return ranks[i].Count > ranks[j].Count
	})


	return ranks[:results]
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}