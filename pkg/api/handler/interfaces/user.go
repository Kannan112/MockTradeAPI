package interfaces

import "github.com/gin-gonic/gin"

type UserHandler interface {
	RegisterUser(ctx *gin.Context)
	Login(ctx *gin.Context)

	StreamMarketData(c *gin.Context)
	WebSocketTestPage(c *gin.Context)

	OrderHandler(c *gin.Context)
	AllOrders(c *gin.Context)
	OrderDetails(c *gin.Context)
	DeteleTrade(c *gin.Context)
}
