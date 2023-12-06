package models

type Stock struct {
	StockID int64   `json:"stockid"`
	Name    string  `json:"name"`
	Price   float32 `json:"price"`
	Company string  `json:"company"`
}
