package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
	"github.com/kannan112/mock-trading-platform-api/pkg/repository/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/service/token"
	service "github.com/kannan112/mock-trading-platform-api/pkg/usecase/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/utils"
)

type userUserCase struct {
	userRepo     interfaces.UserRepository
	tokenService token.TokenService
}

func NewUserUseCase(userRepo interfaces.UserRepository, tokenService token.TokenService) service.UserUseCase {
	return &userUserCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

const (
	AccessTokenDuration = time.Minute * 20
)

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

func (c *userUserCase) UserLogin(ctx context.Context, body request.LoginRequest) (response.Token, error) {

	exists, err := c.userRepo.FindUserByEmail(ctx, body.Email)
	if err != nil {
		return response.Token{}, err
	}
	fmt.Println(exists)
	if !exists {
		return response.Token{}, errors.New("user not exists")
	}

	hashPassword, err := c.userRepo.ExtractPassword(ctx, body.Email)
	if err != nil {
		return response.Token{}, err
	}

	verify := utils.VerifyHashAndPassword(hashPassword, body.Password)
	if !verify {
		return response.Token{}, errors.New("worng password")
	}

	uid, err := c.userRepo.GetUserId(ctx, body.Email)
	if err != nil {
		return response.Token{}, err
	}

	tokenDetails := token.GenerateTokenRequest{
		UserID:   uint(uid),
		UsedFor:  "user",
		ExpireAt: time.Now().Add(24 * time.Hour),
	}

	token, err := c.tokenService.GenerateToken(tokenDetails)
	if err != nil {
		return response.Token{}, err
	}
	return response.Token{
		AccessToken: token.TokenString,
	}, nil

}
