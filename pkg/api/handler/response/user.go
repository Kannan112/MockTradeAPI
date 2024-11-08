package response

type Token struct {
	AccessToken string
}

type OrderResponse struct {
	OrderID string  `json:"orderId"`
	Symbol  string  `json:"symbol"`
	Volume  float32 `json:"volume"`
	Price   float64 `json:"price"`
	Type    string  `json:"type"`
	Status  string  `json:"status"`
}
