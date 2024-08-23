package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const (
	experienceTable         = "experiences"
	experienceRevisionTable = "experience_revisions"
)

type experience struct {
	db  *mysql.Client
	now func() time.Time
}

func newExperience(db *mysql.Client) database.Experience {
	return &experience{
		db:  db,
		now: jst.Now,
	}
}

func (e *experience) List(ctx context.Context, params *database.ListExperiencesParams, fields ...string) (entity.Experiences, error) {
	// TODO: 詳細の実装
	return entity.Experiences{}, nil
}

func (e *experience) Count(ctx context.Context, params *database.ListExperiencesParams) (int64, error) {
	// TODO: 詳細の実装
	return 0, nil
}

func (e *experience) MultiGet(ctx context.Context, experienceIDs []string, fields ...string) (entity.Experiences, error) {
	// TODO: 詳細の実装
	return entity.Experiences{}, nil
}

func (e *experience) MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Experiences, error) {
	// TODO: 詳細の実装
	return entity.Experiences{}, nil
}

func (e *experience) Get(ctx context.Context, experienceID string, fields ...string) (*entity.Experience, error) {
	// TODO: 詳細の実装
	return &entity.Experience{}, nil
}

func (e *experience) Create(ctx context.Context, experience *entity.Experience) error {
	// TODO: 詳細の実装
	return nil
}

func (e *experience) Update(ctx context.Context, experienceID string, params *database.UpdateExperienceParams) error {
	// TODO: 詳細の実装
	return nil
}

func (e *experience) Delete(ctx context.Context, experienceID string) error {
	// TODO: 詳細の実装
	return nil
}
