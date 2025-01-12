package fusionbrain_api

// Модель
type ModelItem struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Version float64 `json:"version"`
	Type    string  `json:"type"`
}

// Ответ со списком моделей
type ModelsResponse []ModelItem

// Текстовый промпт для генерации
type GenerateParams struct {
	Query string `json:"query"`
}

// Запрос генерации
type GenerateRequest struct {
	Type                 string         `json:"type"`
	Style                string         `json:"style,omitempty"`
	Width                int            `json:"width,omitempty"`
	Height               int            `json:"height,omitempty"`
	NumImages            int            `json:"num_images,omitempty"`
	NegativePromptUnclip string         `json:"negativePromptUnclip,omitempty"`
	GenerateParams       GenerateParams `json:"generateParams"`
}

// Ответ на генерацию
type GenerateResponse struct {
	Uuid             string   `json:"uuid"`
	Status           string   `json:"status"`
	StatusTime       int      `json:"status_time,omitempty"`
	Images           []string `json:"images,omitempty"`
	ErrorDescription string   `json:"errorDescription,omitempty"`
	Censored         bool     `json:"censored,omitempty"`
}

// Стиль
type StyleItem struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	TitleEn string `json:"titleEn"`
	Image   string `json:"image"`
}

// Список стилей
type StyleResponse []StyleItem

// Доступность
type AvailabilityResponse struct {
	ModelStatus string `json:"model_status"`
}

type StatusResponse struct {
	Images []string `json:"images"`
}
