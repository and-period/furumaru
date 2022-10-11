package database

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testOrderItem(id, orderID, productID string, now time.Time) *entity.OrderItem {
	return &entity.OrderItem{
		ID:         id,
		OrderID:    orderID,
		ProductID:  productID,
		Price:      100,
		Quantity:   1,
		Weight:     1000,
		WeightUnit: entity.WeightUnitGram,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
