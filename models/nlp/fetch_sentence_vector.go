package nlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
)

type Response struct {
	EncodedText string `json:"encoded_text"`
}

func FetchSentenceVector(text string, wg *sync.WaitGroup, ch chan []float64) {
	defer wg.Done()
	jsonStr := []byte(fmt.Sprintf(`{"text": "%s"}`, text))
	req, err := http.NewRequest("POST",
		"http://localhost:8000/generate",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Fatal(err)
	}
	input_vec := []float64{}
	json.Unmarshal([]byte(response.EncodedText), &input_vec)
	defer resp.Body.Close()
	ch <- input_vec
}
