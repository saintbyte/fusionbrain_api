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
	uuid "github.com/nu7hatch/gouuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
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
func (f *Fusionbrain) generateParams(gr GenerateRequest) (bytes.Buffer, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	defer writer.Close()
	// Добавление model_id
	err := writer.WriteField("model_id", f.CurrentModel)
	modelPart, err := writer.CreateFormFile("model_id", "")
	if err != nil {
		return buf, err
	}
	_, err = modelPart.Write([]byte("4"))
	if err != nil {
		return buf, err
	}

	// Добавление params
	paramsJSON, err := json.Marshal(gr)
	if err != nil {
		return buf, err
	}
	err = writer.WriteField("params", string(paramsJSON))
	if err != nil {
		return buf, err
	}

	// Создаем буфер для multipart/form-data

	// Добавляем файл в multipart/form-data
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Println("Ошибка при создании части формы для файла:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Ошибка при копировании файла:", err)
		return
	}

	// Добавляем переменную genera в multipart/form-data

	if err != nil {
		fmt.Println("Ошибка при добавлении переменной genera:", err)
		return
	}

	// Закрываем writer, чтобы завершить multipart/form-data
	err = writer.Close()
	if err != nil {
		fmt.Println("Ошибка при закрытии writer:", err)
		return
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", apiURL, &buf)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Устанавливаем заголовок Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка при загрузке файла:", resp.Status)
		return
	}

	// Создаем файл для записи обработанного содержимого
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	// Копируем содержимое ответа в файл
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
		return
	}

	fmt.Println("Файл успешно загружен и обработан.")

	return buf, nil
}
func (f *Fusionbrain) Generate(query string, negativeQuery string, style string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	u, err := uuid.NewV4()
	boundary := "------------------------" + u.String()
	requestUrl := f.getUrl("/key/api/v1/text2image/run")
	gr := GenerateRequest{
		Type: "GENERATE",
		GenerateParams: GenerateParams{
			Query: query,
		},
	}
	reqBody, err := f.generateParams(gr)
	client, request, _ := f.getRequest(requestUrl, "POST", &reqBody)
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	log.Println(request)
	response, e := client.Do(request)
	if e != nil {
		log.Fatal(e)
	}
	log.Println(reqBody)
	if response.StatusCode != http.StatusOK {
		log.Fatal(http.StatusOK)
		return "", errors.New("Http error:" + strconv.Itoa(response.StatusCode))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	log.Println(body)
	return string(body), nil
}
