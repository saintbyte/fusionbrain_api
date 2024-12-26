// Package fusionbrain_api Предоставляет доступ API на Fusionbrain.ai
//
// Этот пакет сделан для того чтоб генерить картинки в основном нейросетью Кадински от сбера.
package fusionbrain_api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
)

type Fusionbrain struct {
	ApiHost      string
	ApiKey       string
	SecretKey    string
	Style        string
	CurrentModel string
}

func NewFusionbrain() *Fusionbrain {
	return &Fusionbrain{
		ApiHost:      fusionbrainApiHost,
		ApiKey:       "",
		SecretKey:    "",
		Style:        "",
		CurrentModel: "",
	}
}

func validateSize() {
	//1:1 / 2:3 / 3:2 / 9:16 / 16:9
	//1024
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
}

// (prompt, model string, images, width, height int)
func CreateAudioFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", "audio/wav;rate=8000")
	return w.CreatePart(h)
}
func WriteField(w *multipart.Writer, fieldname, value string) error {
	p, err := w.CreateFormField(fieldname)
	if err != nil {
		return err
	}
	_, err = p.Write([]byte(value))
	return err
}

//func CreateFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {

// }
func (f *Fusionbrain) generateParams(gr GenerateRequest) (bytes.Buffer, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// Добавление model_id
	err := writer.WriteField("model_id", "4")
	if err != nil {
		return buf, "", err
	}
	// ------------------------------------------
	paramsJSON, err := json.Marshal(gr)
	if err != nil {
		return buf, "", err
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "form-data; name=\"params\"")
	h.Set("Content-Type", "application/json")
	p, err := writer.CreatePart(h)
	_, err = p.Write([]byte(paramsJSON))
	if err != nil {
		return buf, "", err
	}
	writer.Close()
	return buf, writer.Boundary(), nil
}
func (f *Fusionbrain) Generate(query string, negativeQuery string, style string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	requestUrl := f.getUrl("/key/api/v1/text2image/run")
	//requestUrl := "http://127.0.0.1:8080/key/api/v1/text2image/run" nc -l 8080
	gr := GenerateRequest{
		Type: "GENERATE",
		GenerateParams: GenerateParams{
			Query: query,
		},
	}
	reqBody, boundary, err := f.generateParams(gr)
	client, request, _ := f.getRequest(requestUrl, "POST", &reqBody)
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	request.Header.Set("Content-Length", strconv.Itoa(reqBody.Len()))
	log.Println(request)
	response, e := client.Do(request)
	if e != nil {
		log.Fatal(e)
	}
	log.Println(reqBody)
	log.Println(response)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("read body error:", err)
	}
	defer response.Body.Close()
	log.Println(string(body))
	if response.StatusCode != http.StatusOK {
		log.Fatal(http.StatusOK)
		return "", errors.New("Http error:" + strconv.Itoa(response.StatusCode))
	}

	return string(body), nil
}
