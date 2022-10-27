package database

import (
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testLive(id, scheduleID, producerID string, productIDs []string, now time.Time) *entity.Live {
	l := &entity.Live{
		ID:          id,
		ScheduleID:  scheduleID,
		ProducerID:  producerID,
		Title:       "配信のタイトル",
		Description: "配信の説明",
		Published:   false,
		Canceled:    false,
		StartAt:     now,
		EndAt:       now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	ps := make(entity.LiveProducts, len(productIDs))
	for i := range productIDs {
		ps[i] = testLiveProduct(id, productIDs[i], now)
	}
	l.Fill(ps, now)
	return l
}

func testLives(id, scheduleID, producerID string, productIDs []string, now time.Time, length int) entity.Lives {
	lives := make(entity.Lives, length)
	for i := 0; i < length; i++ {
		liveID := fmt.Sprintf("%s-%2d", id, i)
		lives[i] = testLive(liveID, scheduleID, producerID, productIDs, now)
	}
	return lives
}

func fillIgnoreLiveField(l *entity.Live, now time.Time) {
	if l == nil {
		return
	}
	l.StartAt = now
	l.EndAt = now
	l.CreatedAt = now
	l.UpdatedAt = now
}

func fillIgnoreLivesField(ls entity.Lives, now time.Time) {
	for i := range ls {
		fillIgnoreLiveField(ls[i], now)
	}
}
