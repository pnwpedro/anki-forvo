package forvo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
)

const MP3_ENCODING_ERROR = "MP3 not downloaded correctly"

const RETRY_BACKOFF = 1 * time.Second

// expected: "{\"attributes\":{\"total\":1},\"items\":[{\"id\":734419,\"word\":\"\\u72d7\",\"original\":\"\\u72d7\",\"addtime\":\"2010-08-22 11:47:12\",\"hits\":6564,\"username\":\"monimonica\",\"sex\":\"f\",\"country\":\"Australia\",\"code\":\"zh\",\"langname\":\"Mandarin Chinese\",\"pathmp3\":\"https:\\/\\/apifree.forvo.com\\/audio\\/393o313e1h1n2o2q2n3l263f3i383f252o1b2n3k2j3m311k1p322a3l3n2q2m3d2d2l281n1i2d252j3b1j382f2m271m1k1g2k3c24363i1f373m2b263f1p362e3j1b3l2f31272i2h1k1b262m3q1m3m1i323521262c3h211t1t_1o2h1h331k3m2a1p24242a2l2a38322q1j1g33283k211t1t\",\"pathogg\":\"https:\\/\\/apifree.forvo.com\\/audio\\/1i1g3h241h3e3f252e263f3h2n2i3b25293k3q2p32271k3e3q313f3g273q2m362328393o212k2c2q253n1p241b3p3k1p3q2a2k382l3o3k26262i3j2e3m2e322j371j3i3d3n3c242e2e333e3d2b2l1f3q1b29351h2m2h1t1t_1n272k3i2j1n1j2b313c2h2p3l343i3o383d3o2k32371t1t\",\"rate\":2,\"num_votes\":2,\"num_positive_votes\":2}]}"

type ListWordPronunciationItem struct {
	ID      int    `json:"id"`
	Word    string `json:"word"`
	User    string `json:"username"`
	PathMP3 string `json:"pathmp3"`
}

type ListWordPronunciationsResult struct {
	Items []ListWordPronunciationItem
}

type Forvo struct {
	client  *http.Client
	baseUrl string
}

func NewForvo(key string) Forvo {
	client := http.DefaultClient
	forvo := Forvo{
		client:  client,
		baseUrl: fmt.Sprintf("https://apifree.forvo.com/key/%s/format/json/action/", key),
	}

	return forvo
}

func (f *Forvo) GetFirstPronunciation(word string, language string) ([]byte, error) {
	body, err := f.getFirstPronunciation(word, language)

	if err != nil && strings.Contains(err.Error(), MP3_ENCODING_ERROR) {
		println("retrying %s in %s", word, RETRY_BACKOFF)
		time.Sleep(RETRY_BACKOFF)
		body, err = f.getFirstPronunciation(word, language)
	}

	return body, err
}

func (f *Forvo) getFirstPronunciation(word string, language string) ([]byte, error) {
	pronunciations, err := f.ListWordPronunciations(word, language)

	if err != nil {
		return nil, err
	}

	println("Got pronunciations from Forvo: ", len(pronunciations))
	if len(pronunciations) == 0 {
		return nil, errors.New("no pronunciations found")
	}

	first := pronunciations[0]
	resp, err := f.client.Get(first.PathMP3)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	_, mp3err := mp3.NewDecoder(bytes.NewReader(body))
	if mp3err != nil {
		fmt.Println(MP3_ENCODING_ERROR)
		return nil, mp3err
	}

	return body, nil
}

func (f *Forvo) ListWordPronunciations(word string, language string) ([]ListWordPronunciationItem, error) {
	urlSuffix := fmt.Sprintf("word-pronunciations/word/%s/language/%s/order/rate-desc/limit/1", word, language)
	get, err := f.client.Get(f.baseUrl + urlSuffix)
	if err != nil {
		return nil, err
	}
	defer get.Body.Close()

	bodyBytes, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, err
	}

	data := ListWordPronunciationsResult{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}

	return data.Items, nil
}
