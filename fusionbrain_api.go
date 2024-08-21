// Package fusionbrain_api Предоставляет доступ API на Fusionbrain.ai
//
// Этот пакет сделан для того чтоб генерить картинки в основном нейросетью Кадински от сбера.
package fusionbrain_api

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"os"
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
	return "https://" + fusionbrainApiHost + "/" + apiPath
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

func (f *Fusionbrain) getRequest(url string, method string, data io.Reader) (*http.Client, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	if method == "GET" {
		request, err := http.NewRequest(method, url, nil)
	} else if method == "POST" {
		request, err := http.NewRequest(method, url, data)
	} else {
		return nil, errors.New("Unknow method")
	}
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Key", "Key "+f.getApiKey())
	request.Header.Set("X-Secret", "Secret "+f.getSecretKey())
	client := &http.Client{}
	return client, nil
}

func (f *Fusionbrain) GetModels() {
	url := f.getUrl("/key/api/v1/models")
	req := f.getRequest(url, "GET", nil)
	//bytes.NewBufferString("scope=GIGACHAT_API_PERS")
}
func (f *Fusionbrain) Generate() {

}
