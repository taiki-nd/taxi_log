package model

type Product struct {
	Id      string `json:"id"`
	PriceId string `json:"price_id"`
	Price   int64  `json:"price"`
}

type Customer struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description"`
}
