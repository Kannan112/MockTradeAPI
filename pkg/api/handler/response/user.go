package response

type Token struct {
	AccessToken string
}

type OrderResponse struct {
	OrderID   uint    `json:"orderId"`
	OrderUUID string  `json:"orderUUID"`
	Symbol    string  `json:"symbol"`
	Volume    float32 `json:"volume"`
	Price     float64 `json:"price"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
}

type MarketData struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"` // Note the string tag
	AskPrice float64 `json:"askPrice,string"` // Note the string tag
	BidQty   float64 `json:"bidQty,string"`
	AskQty   float64 `json:"askQty,string"`
}
