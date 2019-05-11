package main

import (
	"cloud.google.com/go/translate"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
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
	translation := translateText("shoot", "en", "de", client, ctx)
	result := getVocabulary("/go/src/github.com/DanShu93/vocabulary-collector/netflix/", "de", "en_us")

	marshaled, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshaled))
}

func getVocabulary(inputPath, nativeLanguage, targetLanguage string) Vocabulary {
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
							vocabulary := []Vocable{}
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
