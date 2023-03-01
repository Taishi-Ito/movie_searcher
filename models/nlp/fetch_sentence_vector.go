package nlp

import(
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Response struct {
	EncodedText string  `json:"encoded_text"`
}

func FetchSentenceVector(text string) []float64 {
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
	return input_vec
}
