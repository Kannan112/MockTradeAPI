package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
)

var binanceURL = "wss://stream.binance.com:9443/ws/btcusdt@ticker"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket Market Data Stream godoc
// @Summary Get real-time market data stream
// @Description Opens a WebSocket connection to stream real-time market data
// @Tags market-data
// @Accept  json
// @Produce  json
// @Router /api/market-data [get]
func (h *UserHandler) StreamMarketData(c *gin.Context) {

	if h.clients == nil {
		h.clientsMux.Lock()
		h.clients = make(map[*websocket.Conn]bool)
		h.clientsMux.Unlock()
	}
	clientWS, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade client connection: %v", err)
		return
	}
	defer clientWS.Close()

	// Binance WebSocket connection
	binanceWS, _, err := websocket.DefaultDialer.Dial(binanceURL, nil)
	if err != nil {
		log.Printf("Failed to connect to Binance WebSocket: %v", err)
		return
	}
	defer binanceWS.Close()

	// Listen for messages from Binance WebSocket and forward to client
	for {
		_, message, err := binanceWS.ReadMessage()
		if err != nil {
			log.Printf("Error reading Binance WebSocket message: %v", err)
			return
		}

		// Send the Binance message to the connected client
		err = clientWS.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Failed to send message to client: %v", err)
			return
		}
	}
}

// WebSocket Test Page godoc
// @Summary WebSocket Test Page
// @Description HTML page to test WebSocket connection
// @Tags market-data
// @Accept  html
// @Produce  html
// @Success 200 {string} string "HTML page"
// @Router /api/market-live [get]
func (h *UserHandler) WebSocketTestPage(c *gin.Context) {
	c.HTML(http.StatusOK, "websocket.html", nil)
}

// func (h *UserHandler) handleIncomingMessages(ws *websocket.Conn) {
// 	for {
// 		messageType, message, err := ws.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("WebSocket read error: %v", err)
// 			}
// 			return
// 		}

// 		// Handle incoming messages if needed
// 		switch messageType {
// 		case websocket.TextMessage:
// 			// Handle text message
// 			log.Printf("Received message: %s", message)
// 		case websocket.BinaryMessage:
// 			// Handle binary message
// 			log.Printf("Received binary message")
// 		}
// 	}
// }

// OrderHandler godoc
// @Summary Place an order
// @Description Place a buy/sell order with the given details and fetch market data from Binance API.
// @Tags orders
// @Accept json
// @Produce json
// @Param orderRequest body request.OrderRequest true "Order request details"
// @Success 200 {object} response.Response "Order placed successfully"
// @Failure 400 {object} response.Response "Invalid order type"
// @Failure 500 {object} response.Response "Failed to fetch market data"
// @Router /api/order [post]
func (h *UserHandler) OrderHandler(ctx *gin.Context) {
	var orderRequest request.OrderRequest
	if err := ctx.BindJSON(&orderRequest); err != nil {
		response.ErrorResponse(ctx, "Failed to bind JSON", err, nil)
		return
	}

	// Validate order type early
	if orderRequest.Type != "buy" && orderRequest.Type != "sell" {
		response.ErrorResponse(ctx, "Invalid order type", fmt.Errorf("invalid order type: %s", orderRequest.Type), nil)
		return
	}

	marketData, err := FetchMarketData(orderRequest.Symbol)
	if err != nil {
		response.ErrorResponse(ctx, "Failed to fetch market data", err, nil)
		return
	}

	price, err := GetMarketPrice(marketData, orderRequest.Type)
	if err != nil {
		response.ErrorResponse(ctx, "Failed to get market price", err, nil)
		return
	}

	orderID := uuid.New().String()
	orderResponse := response.OrderResponse{
		OrderID: orderID,
		Symbol:  orderRequest.Symbol,
		Volume:  orderRequest.Volume,
		Price:   price,
		Type:    orderRequest.Type,
		Status:  "accepted",
	}

	response.SuccessResponse(ctx, "Order completed", orderResponse)
}

// Example function to fetch market data
func FetchMarketData(symbol string) (MarketData, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/bookTicker?symbol=%s", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return MarketData{}, fmt.Errorf("failed to fetch market data: %w", err)
	}
	defer resp.Body.Close()

	var data MarketData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return MarketData{}, fmt.Errorf("failed to decode market data: %w", err)
	}

	return data, nil
}

type MarketData struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"` // Use string tag to handle string JSON input
	AskPrice float64 `json:"askPrice,string"` // Use string tag to handle string JSON input
}

func GetMarketPrice(marketData MarketData, orderType string) (float64, error) {
	switch orderType {
	case "buy":
		return marketData.AskPrice, nil
	case "sell":
		return marketData.BidPrice, nil
	default:
		return 0, fmt.Errorf("invalid order type: %s", orderType)
	}
}
