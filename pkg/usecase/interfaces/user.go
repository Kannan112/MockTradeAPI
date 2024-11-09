package interfaces

import (
	"context"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
	"github.com/kannan112/mock-trading-platform-api/pkg/utils"
)

type UserUseCase interface {
	CeateNewUser(ctx context.Context, body request.RegisterUserRequest) error
	UserLogin(ctx context.Context, body request.LoginRequest) (response.Token, error)

	FetchMarketData(symbol string) (response.MarketData, error)
	GetMarketPrice(marketData response.MarketData, orderType string) (float64, error)

	CreateOrder(ctx context.Context, uid int, orderData response.OrderResponse) (oid int, err error)
	ListOrders(uid int) ([]utils.OrderResponse, error)
	GetOrderByID(ctx context.Context, uid, oid uint) (utils.Order, error)
	DeleteOrderById(ctx context.Context, uid, oid uint) error
}
