package routes

import (
	"github.com/gin-gonic/gin"
	handlerInterface "github.com/kannan112/mock-trading-platform-api/pkg/api/handler/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup,
	userHandler handlerInterface.UserHandler,

) {

	auth := api.Group("/auth")

	{
		auth.POST("/register", userHandler.RegisterUser)
		auth.POST("/login", userHandler.Login)
	}
	{
		api.GET("/market-data", userHandler.StreamMarketData)
		api.GET("/market-live", userHandler.WebSocketTestPage)
	}

	{
		order := api.Group("/order")
		order.Use(middleware.UserAuth)
		{
			order.POST("", userHandler.OrderHandler)
			order.DELETE(":id", userHandler.DeteleTrade)
			order.GET(":id", userHandler.OrderDetails)
			order.GET("/trade-history", userHandler.AllOrders)
		}

	}

}
