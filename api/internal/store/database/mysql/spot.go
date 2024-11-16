package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const spotTable = "spots"

type spot struct {
	db  *mysql.Client
	now func() time.Time
}

func newSpot(db *mysql.Client) database.Spot {
	return &spot{
		db:  db,
		now: jst.Now,
	}
}

func (s *spot) List(ctx context.Context, params *database.ListSpotsParams, fields ...string) (entity.Spots, error) {
	// TODO: 詳細の実装
	return entity.Spots{}, nil
}

func (s *spot) ListByGeolocation(ctx context.Context, params *database.ListSpotsByGeolocationParams, fields ...string) (entity.Spots, error) {
	// TODO: 詳細の実装
	return entity.Spots{}, nil
}

func (s *spot) Count(ctx context.Context, params *database.ListSpotsParams) (int64, error) {
	// TODO: 詳細の実装
	return 0, nil
}

func (s *spot) Get(ctx context.Context, spotID string, fields ...string) (*entity.Spot, error) {
	// TODO: 詳細の実装
	return &entity.Spot{}, nil
}

func (s *spot) Create(ctx context.Context, spot *entity.Spot) error {
	// TODO: 詳細の実装
	return nil
}

func (s *spot) Update(ctx context.Context, spotID string, params *database.UpdateSpotParams) error {
	// TODO: 詳細の実装
	return nil
}

func (s *spot) Delete(ctx context.Context, spotID string) error {
	// TODO: 詳細の実装
	return nil
}

func (s *spot) Approve(ctx context.Context, spotID string, params *database.ApproveSpotParams) error {
	// TODO: 詳細の実装
	return nil
}
