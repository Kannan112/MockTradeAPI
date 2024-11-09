package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/middleware"
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
// @Security BearerTokenAuth
// @Produce json
// @Param orderRequest body request.OrderRequest true "Order request details"
// @Success 200 {object} response.Response "Order placed successfully"
// @Failure 400 {object} response.Response "Invalid order type"
// @Failure 500 {object} response.Response "Failed to fetch market data"
// @Router /api/order [post]
func (h *UserHandler) OrderHandler(ctx *gin.Context) {

	uid, err := middleware.GetUserIdFromContext(ctx)
	if err != nil {
		response.ErrorResponse(ctx, "Failed to get userid from context", err, nil)
		return
	}

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

	marketData, err := h.userUseCase.FetchMarketData(orderRequest.Symbol)
	if err != nil {
		response.ErrorResponse(ctx, "Failed to fetch market data", err, nil)
		return
	}

	price, err := h.userUseCase.GetMarketPrice(marketData, orderRequest.Type)
	if err != nil {
		response.ErrorResponse(ctx, "Failed to get market price", err, nil)
		return
	}
	fmt.Println(price, "price")

	orderUUID := uuid.New().String()
	orderResponse := response.OrderResponse{
		OrderUUID: orderUUID,
		Symbol:    orderRequest.Symbol,
		Volume:    orderRequest.Volume,
		Price:     price,
		Type:      orderRequest.Type,
		Status:    "accepted",
	}
	oid, err := h.userUseCase.CreateOrder(ctx, uid, orderResponse)
	if err != nil {
		response.ErrorResponse(ctx, "failed to create order", err, nil)
		return
	}
	orderResponse.OrderID = uint(oid)

	response.SuccessResponse(ctx, "order completed", orderResponse)
}

// AllOrders godoc
// @Summary List all orders
//
// @Security BearerTokenAuth
//
// @Description Retrieve all buy/sell orders for the authenticated user.
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Order list retrieved successfully"
// @Failure 400 {object} response.Response "User ID not found in context"
// @Failure 500 {object} response.Response "Failed to retrieve order list"
// @Router /api/order/trade-history [get]
func (h *UserHandler) AllOrders(c *gin.Context) {

	uid, err := middleware.GetUserIdFromContext(c)
	if err != nil {
		response.ErrorResponse(c, "Faild to get user id from context", err, nil)
		return
	}
	data, err := h.userUseCase.ListOrders(uid)
	if err != nil {
		response.ErrorResponse(c, "Failed to retrieve order list", err, nil)
		return
	}

	response.SuccessResponse(c, "Order list retrieved successfully", data)
}

// OrderDetails godoc
// @Summary Get order details
// @Description Retrieve the details of a specific order by order ID for the authenticated user.
// @Tags orders
// @Accept json
// @Security BearerTokenAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response "Order details retrieved successfully"
// @Failure 400 {object} response.Response "Invalid order ID"
// @Failure 500 {object} response.Response "Failed to retrieve order details"
// @Router /api/order/{id} [get]
func (h *UserHandler) OrderDetails(c *gin.Context) {
	idStr := c.Param("id")
	orderid, err := strconv.Atoi(idStr)
	if err != nil {
		response.ErrorResponse(c, "Faild to get user id from context", err, nil)
		return
	}
	uid, err := middleware.GetUserIdFromContext(c)
	if err != nil {
		response.ErrorResponse(c, "Faild to get user id from context", err, nil)
	}

	data, err := h.userUseCase.GetOrderByID(c, uint(uid), uint(orderid))
	if err != nil {
		response.ErrorResponse(c, "failed to get order details", err, nil)
		return
	}
	response.SuccessResponse(c, "order details", data)
}

// DeteleTrade godoc
// @Summary Delete an order
// @Description Delete a specific order by order ID for the authenticated user.
// @Tags orders
// @Accept json
// @Security BearerTokenAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response "Order deleted successfully"
// @Failure 400 {object} response.Response "Invalid order ID"
// @Failure 500 {object} response.Response "Failed to delete order"
// @Router /api/order/{id} [delete]
func (h *UserHandler) DeteleTrade(c *gin.Context) {
	idStr := c.Param("id")
	orderid, err := strconv.Atoi(idStr)
	if err != nil {
		response.ErrorResponse(c, "Faild to get user id from context", err, nil)
		return
	}
	uid, err := middleware.GetUserIdFromContext(c)
	if err != nil {
		response.ErrorResponse(c, "Faild to get user id from context", err, nil)
	}

	err = h.userUseCase.DeleteOrderById(c, uint(uid), uint(orderid))
	if err != nil {
		response.ErrorResponse(c, "failed to get order details", err, nil)
		return
	}
	response.SuccessResponse(c, "order deleted")
}
