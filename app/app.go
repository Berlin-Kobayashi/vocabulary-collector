package main

import (
	"bufio"
	"cloud.google.com/go/translate"
	"encoding/json"
	"github.com/antchfx/xmlquery"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Vocabulary struct {
	NativeLanguages []NativeLanguage `json:"nativeLanguages"`
}

type NativeLanguage struct {
	NativeLanguage  string           `json:"nativeLanguage"`
	TargetLanguages []TargetLanguage `json:"targetLanguages"`
}

type TargetLanguage struct {
	TargetLanguage string   `json:"targetLanguage"`
	Series         []Series `json:"series"`
}

type Series struct {
	Name    string   `json:"name"`
	Seasons []Season `json:"seasons"`
}

type Season struct {
	Season   int       `json:"season"`
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	Episode    int       `json:"episode"`
	Vocabulary []Vocable `json:"vocabulary"`
}

type Vocable struct {
	Original    string `json:"original"`
	Translation string `json:"translation"`
	Seconds     int    `json:"seconds"`
}

type Path struct {
	NativeLanguage, TargetLanguage, Series string
	Season, Episode                        int
}

func main() {
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	result := getVocabulary("/Users/danshu/go/src/vocabulary-collector/netflix/", "de", "en_us", client, ctx)

	marshaled, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	e := ioutil.WriteFile("/Users/danshu/go/src/vocabulary-collector/output/vocabulary.json", marshaled, 0666)
	if e != nil {
		panic(e)
	}
}

func getVocabulary(inputPath, nativeLanguage, targetLanguage string, client *translate.Client, ctx context.Context) Vocabulary {
	series := []Series{}
	seriesNames, err := ioutil.ReadDir(inputPath)
	if err != nil {
		panic(err)
	}

	for _, seriesName := range seriesNames {
		if seriesName.IsDir() {
			seasonNames, err := ioutil.ReadDir(inputPath + seriesName.Name())
			if err != nil {
				panic(err)
			}

			seasons := []Season{}
			for _, seasonName := range seasonNames {
				if seasonName.IsDir() {
					seasonNumber := getSeasonNumber(seasonName.Name())
					episodeNames, err := ioutil.ReadDir(inputPath + seriesName.Name() + "/" + seasonName.Name() + "/" + targetLanguage)
					if err != nil {
						panic(err)
					}

					episodes := []Episode{}
					for _, episodeName := range episodeNames {
						if !episodeName.IsDir() && episodeName.Name() != ".DS_Store" {
							path := inputPath + seriesName.Name() + "/" + seasonName.Name() + "/" + targetLanguage + "/" + episodeName.Name()
							vocabulary := getVocables(path, client, ctx)
							episodes = append(episodes, Episode{Episode: getEpisodeNumber(episodeName.Name()), Vocabulary: vocabulary})
						}
					}
					seasons = append(seasons, Season{Season: seasonNumber, Episodes: episodes})
				}
			}
			newSeries := Series{Name: seriesName.Name(), Seasons: seasons}
			series = append(series, newSeries)
		}
	}

	vocabulary := Vocabulary{NativeLanguages: []NativeLanguage{{NativeLanguage: nativeLanguage, TargetLanguages: []TargetLanguage{{TargetLanguage: targetLanguage, Series: series}}}}}

	return vocabulary
}

func getVocables(path string, client *translate.Client, ctx context.Context) []Vocable {
	vocables := []Vocable{}

	topVocables := parse(path, 50, 6)
	i := 0
	for _, el := range topVocables {
		w := strings.ToLower(el.Word)
		translation := translateText(w, "en", "de", client, ctx)
		if strings.ToLower(translation) != w {
			seconds := parseSeconds(el.Begin[0])
			vocable := Vocable{Original: w, Translation: translation, Seconds: seconds}
			vocables = append(vocables, vocable)
			i++
		}
		if i >= 10 {
			break
		}
	}

	return vocables
}

func parseSeconds(s string) int {
	s = s[:len(s)-5]
	i, _ := strconv.Atoi(s)

	return i / 1000
}

func getSeasonNumber(seasonName string) int {
	seasonName = strings.Replace(seasonName, "Season ", "", -1)

	seasonNumber, _ := strconv.Atoi(seasonName)

	return seasonNumber
}

func getEpisodeNumber(episodeName string) int {
	re := regexp.MustCompile(`.*E(.*)$`)

	match := re.ReplaceAllString(episodeName, "$1")

	episodeNumber, _ := strconv.Atoi(match)

	return episodeNumber
}

func translateText(text, from, to string, client *translate.Client, ctx context.Context) string {
	target, err := language.Parse(to)
	if err != nil {
		log.Fatalf("Failed to parse target language: %v", err)
	}

	source, err := language.Parse(from)
	if err != nil {
		log.Fatalf("Failed to parse source language: %v", err)
	}

	opts := translate.Options{Format: translate.Text, Model: "nmt", Source: source}

	translations, err := client.Translate(ctx, []string{strings.ToLower(text)}, target, &opts)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	return translations[0].Text
}

type Cnt struct {
	Count int
	Begin []string
	End   []string
	Word  string
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
			if !IsLetter(word) {
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
					Word:  word,
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

	if len(ranks) > results {
		ranks = ranks[:results]
	}

	return ranks
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
