package usecase

import (
	"github.com/kannan112/mock-trading-platform-api/pkg/repository/interfaces"
	service "github.com/kannan112/mock-trading-platform-api/pkg/usecase/interfaces"
)

type orderUseCase struct {
	orderRepo interfaces.OrderRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository) service.OrderUseCase {
	return &orderUseCase{
		orderRepo: orderRepo,
	}
}
