package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kannan112/mock-trading-platform-api/pkg/service/token"
)

type Middleware interface {
	AuthenticateUser() gin.HandlerFunc
	TrimSpaces() gin.HandlerFunc
}

type middleware struct {
	tokenService token.TokenService
}

func NewMiddleware(tokenService token.TokenService) Middleware {
	return &middleware{
		tokenService: tokenService,
	}
}
