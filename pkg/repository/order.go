package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
	"github.com/kannan112/mock-trading-platform-api/pkg/repository/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/utils"
	"gorm.io/gorm"
)

type orderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) interfaces.OrderRepository {
	return &orderDatabase{DB: DB}
}

// need to add gorm model to
func (c *orderDatabase) PlaceOrder(ctx context.Context, uid int, data response.OrderResponse) (int, error) {
	fmt.Println(data, "data")
	var oid int
	query := `INSERT INTO orders (order_uuid, user_id, symbol, volume, type, price, status, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)RETURNING id`

	createdAt := time.Now()
	err := c.DB.Raw(query, data.OrderUUID, uid, data.Symbol, data.Volume, data.Type, data.Price, data.Status, createdAt).Scan(&oid).Error
	return oid, err
}

func (c *orderDatabase) GetAllOrders(uid int) ([]utils.OrderResponse, error) {
	var dbOrders []utils.Order

	query := `
        SELECT id, order_uuid, symbol, volume, price, type, status, created_at 
        FROM orders 
        WHERE user_id = $1
        ORDER BY created_at DESC`

	err := c.DB.Raw(query, uid).Scan(&dbOrders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	orders := make([]utils.OrderResponse, len(dbOrders))
	for i, dbOrder := range dbOrders {
		orders[i] = utils.OrderResponse{
			OrderID:   dbOrder.ID,
			OrderUUID: dbOrder.OrderUUID,
			Symbol:    dbOrder.Symbol,
			Volume:    dbOrder.Volume,
			Price:     dbOrder.Price,
			Type:      dbOrder.Type,
			Status:    dbOrder.Status,
			CreatedAt: dbOrder.CreatedAt,
		}
	}

	return orders, nil
}

func (c *orderDatabase) GetOrderByID(oid, uid uint) (utils.Order, error) {
	var order utils.Order

	query := `
        SELECT id, order_uuid, symbol, volume, price, type, status, created_at 
        FROM orders 
        WHERE user_id = $1 AND id = $2
        LIMIT 1`

	result := c.DB.Raw(query, uid, oid).Scan(&order)
	if result.Error != nil {
		return utils.Order{}, fmt.Errorf("failed to fetch order: %w", result.Error)
	}

	// Check if order was found
	if result.RowsAffected == 0 {
		return utils.Order{}, fmt.Errorf("order not found with ID: %d", oid)
	}

	return order, nil
}

func (c *orderDatabase) DeleteOrderById(oid, uid uint) error {
	query := `DELETE FROM orders WHERE user_id = $1 AND id = $2`

	err := c.DB.Exec(query, uid, oid).Error
	return err
}
