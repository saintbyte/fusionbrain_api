// Package fusionbrain_api Предоставляет доступ API на Fusionbrain.ai
//
// Этот пакет сделан для того чтоб генерить картинки в основном нейросетью Кадински от сбера.
package fusionbrain_api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Fusionbrain struct {
	ApiHost   string
	ApiKey    string
	SecretKey string
}

func NewFusionbrain() *Fusionbrain {
	return &Fusionbrain{
		ApiHost:   fusionbrainApiHost,
		ApiKey:    "",
		SecretKey: "",
	}
}

func (f *Fusionbrain) getUrl(apiPath string) string {
	return "https://" + fusionbrainApiHost + apiPath
}

func (f *Fusionbrain) getSecretKey() string {
	value, exists := os.LookupEnv(fusionbrainSecretKeyEnv)
	if exists {
		return value
	}
	if f.SecretKey != "" {
		return f.SecretKey
	}
	return ""
}

func (f *Fusionbrain) getApiKey() string {
	value, exists := os.LookupEnv(fusionbrainApiKeyEnv)
	if exists {
		return value
	}
	if f.SecretKey != "" {
		return f.ApiKey
	}
	return ""
}

func (f *Fusionbrain) getRequest(url string, method string, data io.Reader) (*http.Client, *http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Key", "Key "+f.getApiKey())
	request.Header.Set("X-Secret", "Secret "+f.getSecretKey())
	client := &http.Client{}
	return client, request, nil
}

func (f *Fusionbrain) GetModels() (ModelsResponse, error) {
	url := f.getUrl("/key/api/v1/models")
	log.Println(url)
	client, request, _ := f.getRequest(url, "GET", nil)
	response, e := client.Do(request)
	if e != nil {
		log.Fatal(e)
	}
	if response.StatusCode != http.StatusOK {
		return ModelsResponse{}, errors.New("Http error:" + strconv.Itoa(response.StatusCode))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	var result ModelsResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	return result, nil
	//bytes.NewBufferString("scope=GIGACHAT_API_PERS")
}
func (f *Fusionbrain) Generate() {
	url := f.getUrl("/key/api/v1/text2image/run")
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	params := w.CreateFormField("params")
	params.Write()
	model := w.CreateFormField("model_id")

	req, err := http.NewRequest("POST", url, &b)

}
