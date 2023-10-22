package order

import (
	"database/sql"
	"fmt"

	"github.com/AdamWalker95/orders-api/model"
)

type SqlRepo struct {
	Client *sql.DB
}

func (s *SqlRepo) FindByID(user int) (map[int]model.Order, error) {

	foundOrders := make(map[int]model.Order)

	query := "SELECT * FROM ORDERS WHERE customer_id = ?;"
	rows, err := s.Client.Query(query, user)
	if err != nil {
		return map[int]model.Order{}, fmt.Errorf("Failed to find any orders on system")
	}
	defer rows.Close()
	for i := 0; rows.Next(); i++ {
		var nextOrder model.Order
		err = rows.Scan(&nextOrder.OrderID,
			&nextOrder.CustomerID,
			&nextOrder.CreatedAt,
			&nextOrder.ShippedAt,
			&nextOrder.CompletedAt)

		if err != nil {
			return map[int]model.Order{}, fmt.Errorf("Error occurred when retrieving orders: %w", err)
		}
		foundOrders[i] = nextOrder
	}

	return foundOrders, nil
}
