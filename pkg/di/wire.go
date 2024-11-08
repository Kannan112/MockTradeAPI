//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "github.com/kannan112/mock-trading-platform-api/pkg/api"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/middleware"
	"github.com/kannan112/mock-trading-platform-api/pkg/config"
	"github.com/kannan112/mock-trading-platform-api/pkg/db"
	"github.com/kannan112/mock-trading-platform-api/pkg/repository"
	"github.com/kannan112/mock-trading-platform-api/pkg/service/token"
	"github.com/kannan112/mock-trading-platform-api/pkg/usecase"
)

func InitializeApi(cfg config.Config) (*http.ServerHTTP, error) {

	wire.Build(db.ConnectDatabase,
		//external
		token.NewTokenService,

		// repository
		middleware.NewMiddleware,
		repository.NewUserRepository,

		//usecase
		usecase.NewUserUseCase,

		// handler
		handler.NewUserHandler,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
