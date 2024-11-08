package interfaces

import (
	"context"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/domain"
)

//go:generate mockgen -destination=../../mock/mockrepo/user_mock.go -package=mockrepo . UserRepository
type UserRepository interface {
	FindUserByUserID(ctx context.Context, userID uint) (user domain.User, err error)
	FindUserByEmail(ctx context.Context, email string) (bool, error)
	ExtractPassword(ctx context.Context, email string) (string, error)

	GetUserId(ctx context.Context, email string) (int, error)
	SaveUser(ctx context.Context, user request.RegisterUserRequest) (userID uint, err error)
}
