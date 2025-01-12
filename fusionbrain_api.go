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
	"log/slog"
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
	CurrentModel ModelItem
}

func NewFusionbrain() *Fusionbrain {
	return &Fusionbrain{
		ApiHost:      fusionbrainApiHost,
		ApiKey:       "",
		SecretKey:    "",
		Style:        "",
		CurrentModel: ModelItem{},
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

// Модель по умолчанию , чтоб каждый раз не тянут список и выбирать просто захардкорим это тут
func (f *Fusionbrain) getDefaultModel() ModelItem {
	return ModelItem{
		Id:      4,
		Name:    "Kandinsky",
		Version: 3.1,
		Type:    "TEXT2IMAGE",
	}
}

func (f *Fusionbrain) getModel() ModelItem {
	if f.CurrentModel == (ModelItem{}) {
		return f.getDefaultModel()
	}
	return f.CurrentModel
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
	url := f.getUrl(fusionbrainModelsPath)
	log.Println(url)
	client, request, _ := f.getRequest(url, "GET", nil)
	response, e := client.Do(request)
	if e != nil {
		slog.Error("GetModels client do error:", e)
		return ModelsResponse{}, e
	}
	if response.StatusCode != http.StatusOK {
		return ModelsResponse{}, errors.New("Http error:" + strconv.Itoa(response.StatusCode))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("GetModels read body error:", err)
		return ModelsResponse{}, e
	}
	defer response.Body.Close()

	var result ModelsResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		slog.Error("GetModels json.Unmarshal error:", err2)
		return ModelsResponse{}, e
	}
	return result, nil
}

func (f *Fusionbrain) generateParams(gr GenerateRequest) (bytes.Buffer, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// (prompt, model string, images, width, height int)
	// Добавление model_id
	err := writer.WriteField("model_id", strconv.Itoa(f.getDefaultModel().Id))
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
func (f *Fusionbrain) СheckStatus(uuidOfRequest string) (GenerateResponse, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	requestUrl := f.getUrl(fusionbrainCheckStatusPath + uuidOfRequest)
	client, request, _ := f.getRequest(requestUrl, "GET", nil)
	response, e := client.Do(request)
	if e != nil {
		slog.Error("checkStatus client.Do error:", e)
		return GenerateResponse{}, e
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("checkStatus read body error:", err)
		return GenerateResponse{}, err
	}
	defer response.Body.Close()
	var result GenerateResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		slog.Error("checkStatus json.Unmarshal error:", err2)
		return GenerateResponse{}, err2
	}
	return result, nil
}

// Работает не понятно как , в документации есть а так Unauthorized 401
func (f *Fusionbrain) availability() (AvailabilityResponse, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	requestUrl := f.getUrl(fusionbrainAvailabilityPath)
	client, request, _ := f.getRequest(requestUrl, "GET", nil)
	response, e := client.Do(request)
	if e != nil {
		slog.Error("Availability client do error:", e)
		return AvailabilityResponse{}, e
	}
	log.Println(response)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Availability read body error:", err)
		return AvailabilityResponse{}, e
	}
	slog.Info(string(body))
	defer response.Body.Close()

	var result AvailabilityResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		slog.Error("Availability json.Unmarshal error:", err)
		return AvailabilityResponse{}, e
	}
	return result, nil
}

func (f *Fusionbrain) Generate(query string, negativeQuery string, style string) (GenerateResponse, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	requestUrl := f.getUrl(fusionbrainGeneratePath)
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
	response, e := client.Do(request)
	if e != nil {
		slog.Error("Generate client do error:", e)
		return GenerateResponse{}, e
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Generate  read body error:", err)
		return GenerateResponse{}, e
	}
	defer response.Body.Close()
	log.Println(string(body))
	if response.StatusCode != http.StatusCreated {
		slog.Error("Generate http error:" + strconv.Itoa(response.StatusCode))
		return GenerateResponse{}, errors.New("Http error:" + strconv.Itoa(response.StatusCode))
	}
	var result GenerateResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		slog.Error("Generate json.Unmarshal error:", err)
		return GenerateResponse{}, e
	}
	return result, nil
}
func (f *Fusionbrain) Get(queryId string) (string, error) {
	return "", nil
}
