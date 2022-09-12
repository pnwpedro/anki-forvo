package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"anki-forvo-plugin/forvo"
)

// word,translation,audio,pronunciation,example

func main() {

	lang := "zh"
	// lang := "jp"
	// lang := "es"
	// lang := "grc"

	in := ReadFromCsv(lang)
	out := []*AnkiRecord{}

	for _, c := range in {
		word := strings.TrimSpace(c.Word)
		example := c.Example
		audioPath := c.AudioLocation
		translation := c.Translation
		pronunciation := c.Pronunciation

		if audioPath == "" {
			audioPath = get_audio(word, lang)
		}

		if lang == "zh" {
			// skip dictionary for sentences (ie words with translations already)
			if translation == "" {
				dictionaryEntry := CallMdbg(c.Word)
				translation = strings.Join(dictionaryEntry.English[:], ", ")
				pronunciation = strings.Join(dictionaryEntry.Reading[:], ", ")
			}
		}

		outRecord := &AnkiRecord{
			Word:          word,
			Translation:   translation,
			AudioLocation: audioPath,
			Pronunciation: pronunciation,
			Example:       example,
		}

		out = append(out, outRecord)

	}

	WriteToCsv("fixtures/output_"+lang+".csv", out)
}

//TODO handle moveable ν https://en.wikipedia.org/wiki/Movable_nu

func strip_definite_article(word string, language string) string {
	es_articles := []string{"el", "la", "los", "las"}
	grc_articles := []string{"ὁ", "οἱ", "ἡ", "αἱ", "τό", "τά"}

	language_def_article_map := map[string][]string{
		"es":  es_articles,
		"grc": grc_articles,
	}

	for _, v := range language_def_article_map[language] {
		if strings.HasPrefix(word, v+" ") {
			start_index := len(v) + 1
			return word[start_index:]
		}
	}
	return word
}

func get_audio(word string, lang string) string {
	filepath, forvoErr := download(word, lang)
	if forvoErr != nil {
		println("Forvo err", forvoErr.Error())
	}

	if filepath == "" && (lang == "es" || lang == "grc") {
		word_without_article := strip_definite_article(word, lang)
		if word != word_without_article {
			filepath, forvoErr = download(word_without_article, lang)
		}
	}

	var audioPath string
	if filepath != "" {
		audioPath = "[sound:" + filepath + "]"
	} else {
		audioPath = ""
	}
	return audioPath
}

func download(word string, language string) (string, error) {
	fmt.Println("Starting with", word, language)
	key := os.Getenv("FORVO_API_KEY")
	forvo := forvo.NewForvo(key)
	filepathPrefix := "/Users/theapedroza/Library/Application Support/Anki2/User 1/collection.media/"
	filepathSuffix := word + "_" + language + ".mp3"

	mp3, err := forvo.GetFirstPronunciation(word, language)
	if err != nil {
		return "", err
	}

	out, err := os.Create(filepathPrefix + filepathSuffix)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, bytes.NewReader(mp3))
	if err != nil {
		return "", err
	}

	// TODO return forvo link for ease of manual fixes

	return filepathSuffix, nil
}
