package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type rehearsal struct {
	db  dynamodb.Client
	now func() time.Time
}

func NewRehearsal(db dynamodb.Client) Rehearsal {
	return &rehearsal{
		db:  db,
		now: jst.Now,
	}
}

func (r *rehearsal) Get(ctx context.Context, liveID string) (*entity.Rehearsal, error) {
	rehearsal := &entity.Rehearsal{LiveID: liveID}
	if err := r.db.Get(ctx, rehearsal.PrimaryKey(), rehearsal); err != nil {
		return nil, exception.InternalError(err)
	}
	return rehearsal, nil
}

func (r *rehearsal) Create(ctx context.Context, rehearsal *entity.Rehearsal) error {
	now := r.now()
	rehearsal.CreatedAt, rehearsal.UpdatedAt = now, now

	err := r.db.Insert(ctx, rehearsal)
	return exception.InternalError(err)
}
