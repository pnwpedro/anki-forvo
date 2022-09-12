// package anki_forvo_plugin
package forvo_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"anki-forvo-plugin/forvo"

	"github.com/stretchr/testify/assert"
)

func init() {
	file, err := os.Open("../.env")
	if err != nil {
		println("error setting env")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		vals := strings.Split(scanner.Text(), "=")
		os.Setenv(vals[0], vals[1])
	}

}

func TestForvo_ListWordPronunciations(t *testing.T) {
	lang := "zh"
	word := "狗"
	key := os.Getenv("FORVO_API_KEY")
	forvo := forvo.NewForvo(key)
	pronunciations, err := forvo.ListWordPronunciations(word, lang)

	assert.Nil(t, err)

	res := pronunciations[0]

	assert.Equal(t, 734419, res.ID)
	assert.Equal(t, "狗", res.Word)
	assert.Equal(t, "monimonica", res.User)
	println(res.PathMP3)
	assert.NotEmpty(t, res.PathMP3)
}

func TestForvo_GetFirstPronunciation_Many(t *testing.T) {
	language := "zh"
	words := []string{
		"之字旁",
		"单立刀",
		"硬耳朵",
		"软耳朵",
		"言字旁",
		"单立人",
		"秃宝盖",
		"女字旁",
		"走之旁",
		"反犬旁",
		"绞丝旁",
		"双立人",
		"饮食旁",
		"宝盖头",
		"竖心旁",
		"广子头",
		"大口",
	}

	key := os.Getenv("FORVO_API_KEY")
	forvo := forvo.NewForvo(key)

	for _, w := range words {
		_, err := forvo.GetFirstPronunciation(w, language)

		assert.NoError(t, err)
	}
}
