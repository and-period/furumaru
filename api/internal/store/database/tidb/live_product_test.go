package tidb

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testLiveProduct(liveID, productID string, now time.Time) *entity.LiveProduct {
	return &entity.LiveProduct{
		LiveID:    liveID,
		ProductID: productID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
