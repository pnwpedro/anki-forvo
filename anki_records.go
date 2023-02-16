// package anki_forvo_plugin
package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

type AnkiRecord struct {
	Word             string `csv:"Word"`
	Meaning          string `csv:"Meaning"`
	Article          string `csv:"Article"`
	AudioLocation    string `csv:"AudioLocation"`
	Pronunciation    string `csv:"Pronunciation"`
	Example          string `csv:"Example"`
	Translations     string `csv:"Translations"`
	Gender           string `csv:"Gender"`
	Declension       string `csv:"Declension"`
	GenitiveSingular string `csv:"GenetiveSingular"`
	Conjugation      string `csv:"Conjugation"`
	ForvoSearchField string `csv:"ForvoSearchField"`
	Headword         string `csv:"Headword"`
	Sentence         string `csv:"Sentence"`
}

func (ar *AnkiRecord) ToSlice() []string {
	return []string{
		ar.Word,
		ar.Meaning,
		ar.Article,
		ar.AudioLocation,
		ar.Pronunciation,
		ar.Example,
		ar.Translations,
		ar.Gender,
		ar.Declension,
		ar.GenitiveSingular,
		ar.Conjugation,
		ar.ForvoSearchField,
		ar.Headword,
		ar.Sentence,
	}
}

func ReadFromCsv(lang string) []*AnkiRecord {
	filename := "fixtures/Anki Import-" + lang + ".csv"
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("Error opening file")
		return nil
	}
	defer f.Close()

	inputWords := []*AnkiRecord{}

	if err := gocsv.UnmarshalFile(f, &inputWords); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling file")
	}

	for _, c := range inputWords {
		fmt.Println("Input record: ", c.ToSlice())
	}

	return inputWords
}

func WriteToCsv(filepath string, records []*AnkiRecord) {

	csvFile, err := os.Create(filepath)
	defer csvFile.Close()

	if err != nil {
		println("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	for _, ar := range records {
		err := csvwriter.Write(ar.ToSlice())
		if err != nil {
			println("error writing record to file", err)
		}
	}

	csvwriter.Flush()
}
