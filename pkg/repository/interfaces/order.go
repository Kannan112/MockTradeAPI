package interfaces

import (
	"context"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
	"github.com/kannan112/mock-trading-platform-api/pkg/utils"
)

type OrderRepository interface {
	PlaceOrder(ctx context.Context, uid int, data response.OrderResponse) (int, error)
	GetAllOrders(uid int) ([]utils.OrderResponse, error)
	GetOrderByID(oid, uid uint) (utils.Order, error)
	DeleteOrderById(oid, uid uint) error
}
