package model

type Stock struct {
	ID        string  `json:"id"`
	Fullname  string  `json:"fullname"`
	Shortname string  `json:"shortname"`
	Price     float64 `json:"price"`
}
