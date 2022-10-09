package database

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testOrderActivity(id, orderID, userID string, now time.Time) *entity.OrderActivity {
	return &entity.OrderActivity{
		ID:        id,
		OrderID:   orderID,
		UserID:    userID,
		EventType: entity.OrderEventTypeUnknown,
		Detail:    "支払いが完了しました。",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
