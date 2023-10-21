package mysql

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testOrderItem(orderID, productID string, now time.Time) *entity.OrderItem {
	return &entity.OrderItem{
		OrderID:   orderID,
		ProductID: productID,
		Price:     100,
		Quantity:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
