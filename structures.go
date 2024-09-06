package fusionbrain_api

type ModelItem struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Version float64 `json:"version"`
	Type    string  `json:"type"`
}

type ModelsResponse []ModelItem

type GenerateParams struct {
	Query string `json:"query"`
}

type GenerateRequest struct {
	Type                 string         `json:"type"`
	Style                string         `json:"style,omitempty"`
	Width                int            `json:"width,omitempty"`
	Height               int            `json:"height,omitempty"`
	NumImages            int            `json:"num_images,omitempty"`
	NegativePromptUnclip string         `json:"negativePromptUnclip,omitempty"`
	GenerateParams       GenerateParams `json:"generateParams"`
}

type GenerateResponse struct {
	Uuid             string   `json:"uuid"`
	Status           string   `json:"status"`
	Images           []string `json:"images,omitempty"`
	ErrorDescription string   `json:"errorDescription,omitempty"`
	Censored         string   `json:"censored,omitempty"`
}

type StyleItem struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	TitleEn string `json:"titleEn"`
	Image   string `json:"image"`
}

type StyleResponse []StyleItem
