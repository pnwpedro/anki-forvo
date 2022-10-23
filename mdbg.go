package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResultObj struct {
	Simplified  string
	Traditional []string
	Reading     []string
	English     []string
}
type Result struct {
	Result []ResultObj
}

type Mdgb interface {
	Get(word string) ResultObj
}

type MdgbWeb struct{}

func NewMdgbWeb() *MdgbWeb {
	return &MdgbWeb{}
}

func (m *MdgbWeb) Get(word string) ResultObj {
	httpposturl := "https://zhres.herokuapp.com/api/vocab/match"

	var jsonData = []byte(`{
		"entry": "` + word + `"
	}`)

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	var result Result
	json.Unmarshal([]byte(body), &result)
	if len(result.Result) != 0 {
		return result.Result[0]
	} else {
		var r ResultObj
		return r
	}
}
