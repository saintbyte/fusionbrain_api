package fusionbrain_api

type GenerateRequest struct {
	Type                 string `json:"type"`
	Style                string `json:"style"`
	Width                int    `json:"width"`
	Height               int    `json:"height"`
	NumImages            int    `json:"num_images"`
	NegativePromptUnclip string `json:"negativePromptUnclip,omitempty"`
	GenerateParams       struct {
		Query string `json:"query"`
	} `json:"generateParams"`
}
