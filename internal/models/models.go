package models

type URLPayload struct {
	URL string `json:"url"`
}

type ShortURLResponse struct {
	Result string `json:"result"`
}
