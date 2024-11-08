package interfaces

import (
	"context"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
)

type UserUseCase interface {
	CeateNewUser(ctx context.Context, body request.RegisterUserRequest) error
	UserLogin(ctx context.Context, body request.LoginRequest) (response.Token, error)
}
