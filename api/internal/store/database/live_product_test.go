package database

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

func fillIgnoreLiveProductField(p *entity.LiveProduct, now time.Time) {
	if p == nil {
		return
	}
	p.CreatedAt = now
	p.UpdatedAt = now
}

func fillIgnoreLiveProductsField(ps entity.LiveProducts, now time.Time) {
	for i := range ps {
		fillIgnoreLiveProductField(ps[i], now)
	}
}
