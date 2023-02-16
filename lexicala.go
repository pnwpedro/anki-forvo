package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Headword struct {
	Text           string            `json:"text"`
	Pronunciations map[string]string `json:"pronunciation"`
}
type Translation struct {
	Text string `json:"text"`
}

type Sense struct {
	Definition   string                 `json:"definition"`
	Translations map[string]Translation `json:"translations"`
}

// TODO Headword can be an array or a single value??? (fr vs nl)
type LexicalaResultObj struct {
	Headwords []Headword `json:"headword"`
	Senses    []Sense    `json:"senses"`
}
type LexicalaResult struct {
	Result []LexicalaResultObj `json:"results"`
}

type Lexicala interface {
	Get(word string, lang string) LexicalaResultObj
}

type LexicalaWeb struct{}

func NewLexicala() *LexicalaWeb {
	return &LexicalaWeb{}
}

func (m *LexicalaWeb) Get(word string, lang string) LexicalaResultObj {
	//who knows if morph should be true
	url := "https://lexicala1.p.rapidapi.com/search-entries?text=" + word + "&language=" + lang + "&morph=true"

	req, _ := http.NewRequest("GET", url, nil)

	apiKey := os.Getenv("LEXICALA_API_KEY")
	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", "lexicala1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)

	var result LexicalaResult
	json.Unmarshal([]byte(body), &result)

	if len(result.Result) > 0 {
		//TODO if we start using headwords check the length //&& len(result.Result[0].Headwords) > 0
		return result.Result[0]
	} else {
		var r LexicalaResultObj
		return r
	}
}

// {
// 	"n_results": 1,
// 	"page_number": 1,
// 	"results_per_page": 10,
// 	"n_pages": 1,
// 	"available_n_pages": 1,
// 	"results": [
// 	  {
// 		"id": "FR_DE00110847",
// 		"source": "global",
// 		"language": "fr",
// 		"version": 2,
// 		"frequency": "960",
// 		"headword": [
// 		  {
// 			"text": "rubicond",
// 			"pronunciation": {
// 			  "value": "ʀybikɔ̃"
// 			},
// 			"pos": "adjective",
// 			"additional_inflections": [
// 			  "rubiconde",
// 			  "rubicondes",
// 			  "rubiconds"
// 			]
// 		  },
// 		  {
// 			"text": "rubiconde",
// 			"pronunciation": {
// 			  "value": "ʀybikɔ̃d"
// 			},
// 			"pos": "adjective"
// 		  }
// 		],
// 		"senses": [
// 		  {
// 			"id": "FR_SE00120700",
// 			"definition": "qui est très rouge",
// 			"semantic_subcategory": "rouge",
// 			"translations": {
// 			  "br": {
// 				"text": "rubicundo",
// 				"gender": "masculine",
// 				"inflections": [
// 				  {
// 					"text": "rubicunda",
// 					"gender": "feminine"
// 				  }
// 				]
// 			  },
// 			  "en": {
// 				"text": "ruddy"
// 			  },
// 			  "nl": {
// 				"text": "hoogrood"
// 			  }
// 			},
// 			"examples": [
// 			  {
// 				"text": "un visage rubicond",
// 				"translations": {
// 				  "br": {
// 					"text": "um rosto rubicundo"
// 				  },
// 				  "en": {
// 					"text": "a ruddy face"
// 				  },
// 				  "nl": {
// 					"text": "een hoogrood gezicht"
// 				  }
// 				}
// 			  }
// 			]
// 		  }
// 		]
// 	  }
// 	]
//   }
