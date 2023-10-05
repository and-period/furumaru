package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListByScheduleID(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	lives := entity.Lives{
		{
			ID:         "live-id01",
			ScheduleID: "schedule-id",
			StartAt:    now.AddDate(0, -1, 0),
			EndAt:      now.AddDate(0, 1, 0),
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "live-id02",
			ScheduleID: "schedule-id",
			StartAt:    now.AddDate(0, -1, 0),
			EndAt:      now.AddDate(0, 1, 0),
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListLivesByScheduleIDInput
		expect    entity.Lives
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().ListByScheduleID(ctx, "schedule-id").Return(lives, nil)
			},
			input: &store.ListLivesByScheduleIDInput{
				ScheduleID: "schedule-id",
			},
			expect:    lives,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ListLivesByScheduleIDInput{},
			expect:    nil,
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to list lives",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().ListByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &store.ListLivesByScheduleIDInput{
				ScheduleID: "schedule-id",
			},
			expect:    nil,
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListLivesByScheduleID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetLive(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	live := &entity.Live{
		ID:         "live-id",
		ScheduleID: "schedule-id",
		StartAt:    now.AddDate(0, -1, 0),
		EndAt:      now.AddDate(0, 1, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetLiveInput
		expect    *entity.Live
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
			},
			input: &store.GetLiveInput{
				LiveID: "live-id",
			},
			expect:    live,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetLiveInput{},
			expect:    nil,
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to get live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(nil, assert.AnError)
			},
			input: &store.GetLiveInput{
				LiveID: "live-id",
			},
			expect:    nil,
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetLive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateLive(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	products := entity.Products{
		{
			ID:   "product-id",
			Name: "芽が出たじゃがいも",
		},
	}
	producerIn := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	producer := &uentity.Producer{
		AdminID:  "producer-id",
		Username: "&.農園",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateLiveInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(&entity.Schedule{}, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
				mocks.db.Live.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, live *entity.Live) error {
						expect := &entity.Live{
							ID:           live.ID,
							ScheduleID:   "schedule-id",
							ProducerID:   "producer-id",
							ProductIDs:   []string{"product-id"},
							Comment:      "よろしくお願いします。",
							StartAt:      now.AddDate(0, -1, 0),
							EndAt:        now.AddDate(0, 1, 0),
							LiveProducts: live.LiveProducts,
						}
						assert.Equal(t, expect, live)
						return nil
					})
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateLiveInput{},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(nil, assert.AnError)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInternal,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(&entity.Schedule{}, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, assert.AnError)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInternal,
		},
		{
			name: "failed to unmatch products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(&entity.Schedule{}, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(entity.Products{}, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(&entity.Schedule{}, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(nil, assert.AnError)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInternal,
		},
		{
			name: "failed to create live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(&entity.Schedule{}, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
				mocks.db.Live.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateLive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateLive(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	live := &entity.Live{
		ID:         "live-id",
		ScheduleID: "schedule-id",
		ProducerID: "producer-id",
		ProductIDs: []string{"product-id"},
		Comment:    "よろしくお願いします",
		StartAt:    now.AddDate(0, -1, 0),
		EndAt:      now.AddDate(0, 1, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	products := entity.Products{
		{
			ID:   "product-id",
			Name: "芽が出たじゃがいも",
		},
	}
	params := &database.UpdateLiveParams{
		ProductIDs: []string{"product-id"},
		Comment:    "よろしくお願いします。",
		StartAt:    now.AddDate(0, -1, 0),
		EndAt:      now.AddDate(0, 1, 0),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateLiveInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products, nil)
				mocks.db.Live.EXPECT().Update(ctx, "live-id", params).Return(nil)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateLiveInput{},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to get live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(nil, assert.AnError)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInternal,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(nil, assert.AnError)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInternal,
		},
		{
			name: "failed to unmatch products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(entity.Products{}, nil)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to update live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products, nil)
				mocks.db.Live.EXPECT().Update(ctx, "live-id", params).Return(assert.AnError)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
			},
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateLive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteLive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteLiveInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Delete(ctx, "live-id").Return(nil)
			},
			input: &store.DeleteLiveInput{
				LiveID: "live-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteLiveInput{},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to delete live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Delete(ctx, "live-id").Return(assert.AnError)
			},
			input: &store.DeleteLiveInput{
				LiveID: "live-id",
			},
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteLive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
