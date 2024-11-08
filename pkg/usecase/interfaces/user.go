package interfaces

import (
	"context"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
)

type UserUseCase interface {
	CeateNewUser(ctx context.Context, body request.RegisterUserRequest) error
}
