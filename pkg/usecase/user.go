package usecase

import (
	"context"
	"errors"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/repository/interfaces"
	service "github.com/kannan112/mock-trading-platform-api/pkg/usecase/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/utils"
)

type userUserCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository) service.UserUseCase {
	return &userUserCase{
		userRepo: userRepo,
	}
}

func (c *userUserCase) CeateNewUser(ctx context.Context, body request.RegisterUserRequest) error {

	exists, err := c.userRepo.FindUserByEmail(ctx, body.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	password, err := utils.GenerateHashFromPassword(body.Password)
	if err != nil {
		return err
	}
	body.Password = password

	_, err = c.userRepo.SaveUser(ctx, body)
	if err != nil {
		return err
	}

	return nil

}
