package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromotion(t *testing.T) {
	assert.NotNil(t, NewPromotion(nil))
}

func TestPromotion_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	promotions := make(entity.Promotions, 3)
	promotions[0] = testPromotion("promotion-id01", "code0001", "", now().Add(-time.Hour))
	promotions[0].Status = entity.PromotionStatusFinished
	promotions[1] = testPromotion("promotion-id02", "code0002", "", now())
	promotions[2] = testPromotion("promotion-id03", "code0003", "", now().Add(time.Hour))
	err = db.DB.Create(&promotions).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListPromotionsParams
	}
	type want struct {
		promotions entity.Promotions
		hasErr     bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListPromotionsParams{
					Title:  "夏の採れたて野菜",
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				promotions: promotions[1:2],
				hasErr:     false,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListPromotionsParams{
					Orders: []*database.ListPromotionsOrder{
						{Key: database.ListPromotionsOrderByCreatedAt, OrderByASC: true},
						{Key: database.ListPromotionsOrderByUpdatedAt, OrderByASC: false},
					},
				},
			},
			want: want{
				promotions: promotions,
				hasErr:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.promotions, actual)
		})
	}
}

func TestPromotion_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	promotions := make(entity.Promotions, 3)
	promotions[0] = testPromotion("promotion-id01", "code0001", "", now())
	promotions[1] = testPromotion("promotion-id02", "code0002", "", now())
	promotions[2] = testPromotion("promotion-id03", "code0003", "", now())
	err = db.DB.Create(&promotions).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListPromotionsParams
	}
	type want struct {
		total  int64
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListPromotionsParams{
					Title:  "夏の採れたて野菜",
					Limit:  2,
					Offset: 1,
				},
			},
			want: want{
				total:  3,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestPromotion_MultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	promotions := make(entity.Promotions, 3)
	promotions[0] = testPromotion("promotion-id01", "code0001", "", now())
	promotions[1] = testPromotion("promotion-id02", "code0002", "", now())
	promotions[2] = testPromotion("promotion-id03", "code0003", "", now())
	err = db.DB.Create(&promotions).Error
	require.NoError(t, err)

	type args struct {
		promotionIDs []string
	}
	type want struct {
		promotions entity.Promotions
		hasErr     bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				promotionIDs: []string{
					"promotion-id01",
					"promotion-id02",
					"promotion-id03",
				},
			},
			want: want{
				promotions: promotions,
				hasErr:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.promotionIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.promotions, actual)
		})
	}
}

func TestPromotion_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	p := testPromotion("promotion-id", "code0001", "", now())
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	type args struct {
		promotionID string
	}
	type want struct {
		promotion *entity.Promotion
		hasErr    bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				promotionID: "promotion-id",
			},
			want: want{
				promotion: p,
				hasErr:    false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				promotionID: "other-id",
			},
			want: want{
				promotion: nil,
				hasErr:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.promotionID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.promotion, actual)
		})
	}
}

func TestPromotion_GetByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	p := testPromotion("promotion-id", "code0001", "", now())
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	type args struct {
		code string
	}
	type want struct {
		promotion *entity.Promotion
		err       error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				code: "code0001",
			},
			want: want{
				promotion: p,
				err:       nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				code: "",
			},
			want: want{
				promotion: nil,
				err:       database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			actual, err := db.GetByCode(ctx, tt.args.code)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.promotion, actual)
		})
	}
}

func TestPromotion_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		promotion *entity.Promotion
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				promotion: testPromotion("promotion-id", "code0001", "", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				promotion := testPromotion("promotion-id", "code0001", "", now())
				err = db.DB.Create(&promotion).Error
				require.NoError(t, err)
			},
			args: args{
				promotion: testPromotion("promotion-id", "code0001", "", now()),
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, promotionTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			err = db.Create(ctx, tt.args.promotion)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestPromotion_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		promotionID string
		params      *database.UpdatePromotionParams
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				promotion := testPromotion("promotion-id", "code0001", "", now())
				err := db.DB.Create(&promotion).Error
				require.NoError(t, err)
			},
			args: args{
				promotionID: "promotion-id",
				params: &database.UpdatePromotionParams{
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					DiscountType: entity.DiscountTypeFreeShipping,
					DiscountRate: 0,
					Code:         "code0001",
					CodeType:     entity.PromotionCodeTypeOnce,
					StartAt:      now(),
					EndAt:        now().AddDate(0, 1, 0),
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "code is unique",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				promotion := testPromotion("promotion-id", "code0001", "", now())
				err := db.DB.Create(&promotion).Error
				require.NoError(t, err)
				promotion = testPromotion("other-id", "code0002", "", now())
				err = db.DB.Create(&promotion).Error
				require.NoError(t, err)
			},
			args: args{
				promotionID: "promotion-id",
				params:      &database.UpdatePromotionParams{Code: "code0002"},
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, promotionTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			err = db.Update(ctx, tt.args.promotionID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestPromotion_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		promotionID string
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				promotion := testPromotion("promotion-id", "code0001", "", now())
				err := db.DB.Create(&promotion).Error
				require.NoError(t, err)
			},
			args: args{
				promotionID: "promotion-id",
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, promotionTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &promotion{db: db, now: now}
			err = db.Delete(ctx, tt.args.promotionID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testPromotion(id, code, shopID string, now time.Time) *entity.Promotion {
	targetType := entity.PromotionTargetTypeAllShop
	if shopID != "" {
		targetType = entity.PromotionTargetTypeSpecificShop
	}
	return &entity.Promotion{
		ID:           id,
		ShopID:       shopID,
		Status:       entity.PromotionStatusEnabled,
		Title:        "夏の採れたて野菜マルシェを開催!!",
		Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
		Public:       true,
		TargetType:   targetType,
		DiscountType: entity.DiscountTypeFreeShipping,
		DiscountRate: 0,
		Code:         code,
		CodeType:     entity.PromotionCodeTypeOnce,
		StartAt:      now.Add(-time.Hour),
		EndAt:        now.Add(time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
