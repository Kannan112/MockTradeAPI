package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var binanceURL = "wss://stream.binance.com:9443/ws/btcusdt@ticker"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MarketData struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Time     int64   `json:"time"`
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
