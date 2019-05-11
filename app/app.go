package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Name    string `json:"name"`
	Seasons Season `json:"seasons"`
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
	//ctx := context.Background()

	// _, err := translate.NewClient(ctx)
	// if err != nil {
	//	log.Fatalf("Failed to create client: %v", err)
	//}

	result := getVocabulary("/go/src/github.com/DanShu93/vocabulary-collector/netflix/")

	marshaled, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshaled))
}

func getVocabulary(inputPath string) Vocabulary {
	vocabulary := Vocabulary{NativeLanguages: []NativeLanguage{}}
	nativeLanguageNames, err := ioutil.ReadDir(inputPath)
	if err != nil {
		panic(err)
	}

	for _, n := range nativeLanguageNames {
		if n.IsDir() {
			targetLanguages := []TargetLanguage{}
			targetLanguageNames, err := ioutil.ReadDir(inputPath + n.Name())
			if err != nil {
				panic(err)
			}

			for _, t := range targetLanguageNames {
				if t.IsDir() {
					targetLanguages = append(targetLanguages, TargetLanguage{TargetLanguage: t.Name(), Series: []Series{}})

					_, err := ioutil.ReadDir(inputPath + n.Name() + "/" + t.Name())
					if err != nil {
						panic(err)
					}
				}
			}
			vocabulary.NativeLanguages = append(vocabulary.NativeLanguages, NativeLanguage{NativeLanguage: n.Name(), TargetLanguages: targetLanguages})
		}
	}

	return vocabulary
}
