package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// TODO move API keys
type Payload struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

type DeepLResultObj struct {
	Text                   string `json:"text"`
	DetectedSourceLanguage string `json:"detected_source_language"`
}

type DeepLResult struct {
	Translations []DeepLResultObj `json:"translations"`
}

type DeepL interface {
	Get(word string, lang string) DeepLResultObj
}

type DeepLWeb struct{}

func NewDeepL() *DeepLWeb {
	return &DeepLWeb{}
}

func (m *DeepLWeb) Get(word string, lang string) DeepLResultObj {
	reqURL, _ := url.Parse("https://api-free.deepl.com")
	reqURL.Path = path.Join(reqURL.Path, "v2", "translate")
	apiKey := os.Getenv("DEEPL_API_KEY")

	data := url.Values{}
	data.Set("text", word)
	data.Set("target_lang", "EN")
	data.Set("source_lang", strings.ToUpper(lang))

	req, err := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)

	var result DeepLResult
	json.Unmarshal([]byte(responseBody), &result)
	fmt.Println(result)
	return result.Translations[0]
}
