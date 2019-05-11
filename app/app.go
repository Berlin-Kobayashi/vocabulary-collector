package main

import (
	"cloud.google.com/go/translate"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type Vocabulary struct {
	NativeLanguages []struct {
		NativeLanguage  string `json:"nativeLanguage"`
		TargetLanguages []struct {
			TargetLanguage string `json:"targetLanguage"`
			Series         []struct {
				Name    string `json:"name"`
				Seasons struct {
					Season   int `json:"season"`
					Episodes []struct {
						Episode    int `json:"episode"`
						Vocabulary []struct {
							Original    string `json:"original"`
							Translation string `json:"translation"`
							Seconds     int    `json:"seconds"`
						} `json:"vocabulary"`
					} `json:"episodes"`
				} `json:"seasons"`
			} `json:"series"`
		} `json:"targetLanguages"`
	} `json:"nativeLanguages"`
}

func main() {
	ctx := context.Background()

	_, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	result := Vocabulary{}

	marshaled, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshaled))
}
