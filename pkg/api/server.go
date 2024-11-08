package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/kannan112/mock-trading-platform-api/cmd/api/docs"
	handlerInterface "github.com/kannan112/mock-trading-platform-api/pkg/api/handler/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/middleware"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/routes"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	Engine *gin.Engine
}

// @title					Trading Platform Backend API
// @description				Backend API built with Golang using Clean Code architecture. \nGithub: [https://github.com/kannan112/mock-trading-platform-api].
//
// @contact.name				For API Support
// @contact.email				abhinandarun11@gmail.com
//
// @SecurityDefinitions.apikey	BearerAuth
// @Name						Authorization
// @In							header
// @Description				Add prefix of Bearer before  token Ex: "Bearer token"
// @Query.collection.format	multi
func NewServerHTTP(middleware middleware.Middleware, userHandler handlerInterface.UserHandler) *ServerHTTP {

	engine := gin.New()

	engine.LoadHTMLGlob("views/*.html")

	engine.Use(gin.Logger())

	// swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// set up routes
	routes.UserRoutes(engine.Group("/api"), middleware, userHandler)

	// no handler
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "invalid url go to /swagger/index.html for api documentation",
		})
	})

	return &ServerHTTP{Engine: engine}
}

func (s *ServerHTTP) Start() error {

	return s.Engine.Run(":8080")
}
