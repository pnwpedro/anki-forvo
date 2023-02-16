package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"anki-forvo-plugin/forvo"
)

// word,translation,audio,pronunciation,example,declension,conjugation,gender,genetive singular,forvolookup

func main() {

	// lang := "zh"
	// lang := "ja"
	// lang := "es" //"es_es"
	lang := "grc"
	// lang := "la"
	// lang := "nl"
	// lang := "fr"

	secondPos :=
		false
	// true // for example when the front is "accusative masc singular of dives" or similar

	in := ReadFromCsv(lang)
	out := []*AnkiRecord{}

	var mdgb Mdgb
	mdgb = NewMdgbLocal()

	for _, c := range in {
		word := strings.TrimSpace(c.Word)
		example := c.Example
		audioPath := c.AudioLocation
		translation := c.Translation
		pronunciation := c.Pronunciation
		declension := c.Declension
		conjugation := c.Conjugation
		gender := c.Gender
		genitiveSingular := c.GenitiveSingular
		forvoSearchField := strings.TrimSpace(c.ForvoSearchField)

		if audioPath == "" {
			var searchTerm = forvoSearchField
			if searchTerm == "" {
				searchTerm = word
			}
			audioPath = get_audio(searchTerm, lang)
		}
		if audioPath == "" && secondPos {
			audioPath = get_audio(translation, lang)
		}

		if lang == "zh" {
			// skip dictionary for sentences (ie words with translations already)
			if translation == "" {
				dictionaryEntry := mdgb.Get(word)
				translation = strings.Join(dictionaryEntry.English[:], ", ")
				pronunciation = strings.Join(dictionaryEntry.Reading[:], ", ")
			}

			if pronunciation == "" {
				pronunciation = mdgb.GetPinyin(word)
			}
		}

		outRecord := &AnkiRecord{
			Word:             word,
			Translation:      translation,
			AudioLocation:    audioPath,
			Pronunciation:    pronunciation,
			Example:          example,
			Declension:       declension,
			Conjugation:      conjugation,
			Gender:           gender,
			GenitiveSingular: genitiveSingular,
			ForvoSearchField: forvoSearchField,
		}

		out = append(out, outRecord)

	}

	WriteToCsv("fixtures/output_"+lang+".csv", out)
}

//TODO handle moveable ν https://en.wikipedia.org/wiki/Movable_nu

func strip_definite_article(word string, language string) string {
	es_articles := []string{"el", "la", "los", "las"}
	grc_articles := []string{"ὁ", "οἱ", "ἡ", "αἱ", "τό", "τά", "oἱ", "τό"}
	fr_articles := []string{"le", "la", "les", "des", "un", "une"}
	nl_articles := []string{"de", "het"}

	language_def_article_map := map[string][]string{
		"es":  es_articles,
		"grc": grc_articles,
		"fr":  fr_articles,
		"nl":  nl_articles,
	}

	for _, v := range language_def_article_map[language] {
		start_index := len(v) + 1
		if strings.HasPrefix(word, v+" ") {
			return strings.TrimSpace(word[start_index:])
		}

		suffix := ", " + v
		// println("trying " + suffix + " as a suffix")

		if strings.HasSuffix(word, suffix) {
			return strings.TrimSpace(word[:len(word)-len(v)-2])
		}
		// println("didn't find any article to remove")
	}
	return word
}

func get_audio(word string, lang string) string {
	filepath, forvoErr := download(word, lang)
	if forvoErr != nil {
		println("Forvo err", forvoErr.Error())
	}

	// For words that include a gender + def article, try without
	if filepath == "" && (lang == "es" || lang == "grc" || lang == "fr" || lang == "nl") {
		word_without_article := strip_definite_article(word, lang)
		if word != word_without_article {
			filepath, forvoErr = download(word_without_article, lang)
		}
	}

	// for latin verbs try the first part
	if filepath == "" && lang == "la" && strings.Contains(word, ",") {
		firstPart := strings.Split(word, ",")[0]
		filepath, forvoErr = download(firstPart, lang)
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
		fmt.Println("Error after get first pronunciation", word)
		return "", err
	}

	out, err := os.Create(filepathPrefix + filepathSuffix)
	if err != nil {
		fmt.Println("Error after create file", word)
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, bytes.NewReader(mp3))
	if err != nil {
		fmt.Println("Error after writing body to file", word)
		return "", err
	}
	// TODO return forvo link for ease of manual fixes

	return filepathSuffix, nil
}
