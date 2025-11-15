package order

import (
	"database/sql"
	"strings"

	orderModel "github.com/Christyan39/test-eDot/internal/models/order"
)

func (r *orderRepository) CreateOrderItem(tx *sql.Tx, req []orderModel.OrderItem) error {
	placeholders := make([]string, 0, len(req))
	args := make([]interface{}, 0, len(req)*5)
	for _, item := range req {
		placeholders = append(placeholders, "(?, ?, ?, ?, NOW(), NOW())")
		args = append(args, item.OrderID, item.ProductID, item.Quantity, item.Price)
	}

	query := `
		INSERT INTO order_items (
		order_id,
		product_id,
		quantity,
		item_price,
		created_at,
		updated_at)
		VALUES ` + strings.Join(placeholders, ",") + `	
	`

	_, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
