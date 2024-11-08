package handler

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
	"github.com/kannan112/mock-trading-platform-api/pkg/service/token"
	usecaseInterface "github.com/kannan112/mock-trading-platform-api/pkg/usecase/interfaces"
)

const (
	BindJsonFailMessage string = "Failed to bind JSON request"
)

type UserHandler struct {
	userUseCase usecaseInterface.UserUseCase
	clients     map[*websocket.Conn]bool
	clientsMux  sync.RWMutex
}

func NewUserHandler(userUsecase usecaseInterface.UserUseCase, tokenService token.TokenService) interfaces.UserHandler {
	return &UserHandler{
		userUseCase: userUsecase,
		clients:     make(map[*websocket.Conn]bool),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with a username, email, and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body request.RegisterUserRequest true "User registration details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/auth/register [post]
func (c *UserHandler) RegisterUser(ctx *gin.Context) {
	var body request.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, BindJsonFailMessage, err, nil)
		return
	}
	err := c.userUseCase.CeateNewUser(ctx, body)
	if err != nil {
		response.ErrorResponse(ctx, "Failed to register user", err, nil)
		return
	}
	response.SuccessResponse(ctx, "User registered successfully", nil)
}

// Login godoc
// @Summary Login User
// @Description user login  email, and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body request.LoginRequest true "User login details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/auth/login [post]
func (c *UserHandler) Login(ctx *gin.Context) {
	var body request.LoginRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, BindJsonFailMessage, err, nil)
		return
	}
	tokenRes, err := c.userUseCase.UserLogin(ctx, body)
	if err != nil {
		response.ErrorResponse(ctx, "Login failed: invalid credentials or server error", err, nil)
		return
	}

	response.SuccessResponse(ctx, "Login successful", tokenRes)
}
