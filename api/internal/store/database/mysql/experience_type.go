package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const experienceTypeTable = "experience_types"

type experienceType struct {
	db  *mysql.Client
	now func() time.Time
}

func newExperienceType(db *mysql.Client) database.ExperienceType {
	return &experienceType{
		db:  db,
		now: jst.Now,
	}
}

func (t *experienceType) List(
	ctx context.Context, params *database.ListExperienceTypesParams, fields ...string,
) (entity.ExperienceTypes, error) {
	// TODO: 詳細の実装
	return entity.ExperienceTypes{}, nil
}

func (t *experienceType) Count(ctx context.Context, params *database.ListExperienceTypesParams) (int64, error) {
	// TODO: 詳細の実装
	return 0, nil
}

func (t *experienceType) MultiGet(ctx context.Context, experienceIDs []string, fields ...string) (entity.ExperienceTypes, error) {
	// TODO: 詳細の実装
	return entity.ExperienceTypes{}, nil
}

func (t *experienceType) Get(ctx context.Context, experienceID string, fields ...string) (*entity.ExperienceType, error) {
	// TODO: 詳細の実装
	return &entity.ExperienceType{}, nil
}

func (t *experienceType) Create(ctx context.Context, experience *entity.ExperienceType) error {
	// TODO: 詳細の実装
	return nil
}

func (t *experienceType) Update(ctx context.Context, experienceID string, params *database.UpdateExperienceTypeParams) error {
	// TODO: 詳細の実装
	return nil
}

func (t *experienceType) Delete(ctx context.Context, experienceID string) error {
	// TODO: 詳細の実装
	return nil
}
