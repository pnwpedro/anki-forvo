package main

import (
	"anki-forvo-plugin/forvo"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

//TODO instead of "greek and latin" say "declined languages" or similar

// OLD FIELDS word,translation,audio,pronunciation,example,declension,conjugation,gender,genetive singular,forvolookup
// NEW Greek and Latin Fields
//Word,Meaning,PrepositionAudioLocation,Pronunciation,Example,Translations,Gender,Declension,GenetiveSingular,Conjugation,ForvoSearchField,Headword,Sentence
// NEW Other Language fields
//Word,Meaning,Preposition,AudioLocation,Pronunciation,Example,Translations,ForvoSearchField,Headword,Sentence

var supportedLanguages = []string{"zh", "ja", "es", "grc", "la", "nl", "en", "fr"}

func main() {

	// lang := "zh"
	lang := "ja"
	// lang := "es" //"es_es"
	// lang := "grc"
	// lang := "la"
	// lang := "nl"
	// lang := "fr"

	searchMeaningField :=
		false
	// true // for example when the front is "accusative masc singular of dives" or similar
	// TODO rethink this with new template

	in := ReadFromCsv(lang)
	out := []*AnkiRecord{}

	var mdgb = NewMdgbLocal()
	var lexicala Lexicala = NewLexicala()
	var deepL DeepL = NewDeepL()

	for _, c := range in {
		word := strings.TrimSpace(c.Word)
		meaning := strings.TrimSpace(c.Meaning)
		article := strings.TrimSpace(c.Article)
		audioLocation := c.AudioLocation
		pronunciation := c.Pronunciation
		example := strings.TrimSpace(c.Example)
		translations := c.Translations
		gender := c.Gender
		declension := c.Declension
		genitiveSingular := c.GenitiveSingular
		conjugation := c.Conjugation
		forvoSearchField := strings.TrimSpace(c.ForvoSearchField)
		isHeadword, _ := strconv.ParseBool(c.Headword)
		isSentence, _ := strconv.ParseBool(c.Sentence)

		if audioLocation == "" {
			var searchTerm = forvoSearchField
			if searchTerm == "" {
				searchTerm = word
			}
			audioLocation = get_audio(searchTerm, lang)
		}
		if audioLocation == "" && searchMeaningField {
			audioLocation = get_audio(meaning, lang)
		}

		if lang == "zh" {
			// skip dictionary for sentences (ie words with translations already)
			if meaning == "" {
				dictionaryEntry := mdgb.Get(word)
				meaning = strings.Join(dictionaryEntry.English[:], ", ")
				pronunciation = strings.Join(dictionaryEntry.Reading[:], ", ")
			}

			if pronunciation == "" {
				pronunciation = mdgb.GetPinyin(word)
			}
		} else if meaning == "" && isHeadword {
			lexicalaEntry := lexicala.Get(word, lang)
			fmt.Println(lexicalaEntry)
			if len(lexicalaEntry.Senses) != 0 {
				meaning = lexicalaEntry.Senses[0].Definition
				translations = writeTranslationsToField(lexicalaEntry.Senses[0].Translations)
			}
		} else if meaning == "" && isSentence {
			deepLTranslation := deepL.Get(word, lang)
			if deepLTranslation.Text != "" {
				meaning = deepLTranslation.Text
			}
		}

		outRecord := &AnkiRecord{
			Word:             word,
			Meaning:          meaning,
			Article:          article,
			AudioLocation:    audioLocation,
			Pronunciation:    pronunciation,
			Example:          example,
			Translations:     translations,
			Declension:       declension,
			Conjugation:      conjugation,
			Gender:           gender,
			GenitiveSingular: genitiveSingular,
			ForvoSearchField: forvoSearchField,
			Headword:         strconv.FormatBool(isHeadword),
			Sentence:         strconv.FormatBool(isSentence),
		}

		out = append(out, outRecord)

	}

	WriteToCsv("fixtures/output_"+lang+".csv", out)
}

func writeTranslationsToField(translations map[string]Translation) string {
	out := ""
	for k, v := range translations {
		fmt.Printf("key[%s] value[%s]\n", k, v)
		if slices.Contains(supportedLanguages, k) && translations[k].Text != "" {
			out = out + k + ": " + translations[k].Text + "\n"
		}
	}
	return out
}

//TODO (GRC) handle moveable ν https://en.wikipedia.org/wiki/Movable_nu

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
	return filepathSuffix, nil
}
