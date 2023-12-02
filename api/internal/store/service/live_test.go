package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListLives(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	params := &database.ListLivesParams{
		ScheduleIDs: []string{"schedule-id"},
		Limit:       20,
		Offset:      0,
	}
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
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListLivesInput
		expect      entity.Lives
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().List(gomock.Any(), params).Return(lives, nil)
				mocks.db.Live.EXPECT().Count(gomock.Any(), params).Return(int64(2), nil)
			},
			input: &store.ListLivesInput{
				ScheduleIDs: []string{"schedule-id"},
				Limit:       20,
				Offset:      0,
			},
			expect:      lives,
			expectTotal: 2,
			expectErr:   nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListLivesInput{
				Limit:  1000,
				Offset: -1,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list lives",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Live.EXPECT().Count(gomock.Any(), params).Return(int64(2), nil)
			},
			input: &store.ListLivesInput{
				ScheduleIDs: []string{"schedule-id"},
				Limit:       20,
				Offset:      0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count lives",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().List(gomock.Any(), params).Return(lives, nil)
				mocks.db.Live.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListLivesInput{
				ScheduleIDs: []string{"schedule-id"},
				Limit:       20,
				Offset:      0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListLives(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
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
			expectErr: exception.ErrInvalidArgument,
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
			expectErr: exception.ErrInternal,
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

	products := entity.Products{
		{
			ID:   "product-id",
			Name: "芽が出たじゃがいも",
		},
	}
	schedule := &entity.Schedule{
		ID:      "schedule-id",
		StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
		EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
	}
	livesIn := &database.ListLivesParams{
		ScheduleIDs: []string{"schedule-id"},
	}
	lives := entity.Lives{
		{
			ID:         "live-id",
			ScheduleID: "schedule-id",
			StartAt:    jst.Date(2023, 7, 15, 17, 0, 0, 0),
			EndAt:      jst.Date(2023, 7, 15, 18, 30, 0, 0),
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
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
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
							StartAt:      jst.Date(2023, 7, 15, 19, 30, 0, 0),
							EndAt:        jst.Date(2023, 7, 15, 21, 0, 0, 0),
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
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateLiveInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(nil, assert.AnError)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to list lives",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, assert.AnError)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to unmatch products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(entity.Products{}, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(nil, assert.AnError)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "invalid live schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 18, 0, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 20, 0, 0, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
				mocks.db.Live.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateLiveInput{
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
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
		StartAt:    jst.Date(2023, 7, 15, 17, 0, 0, 0),
		EndAt:      jst.Date(2023, 7, 15, 18, 30, 0, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	schedule := &entity.Schedule{
		ID:      "schedule-id",
		StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
		EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
	}
	livesIn := &database.ListLivesParams{
		ScheduleIDs: []string{"schedule-id"},
	}
	lives := entity.Lives{
		{
			ID:         "live-id",
			ScheduleID: "schedule-id",
			StartAt:    jst.Date(2023, 7, 15, 17, 0, 0, 0),
			EndAt:      jst.Date(2023, 7, 15, 18, 30, 0, 0),
		},
	}
	producerIn := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	producer := &uentity.Producer{
		AdminID:  "producer-id",
		Username: "&.農園",
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
		StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
		EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
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
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
				mocks.db.Live.EXPECT().Update(ctx, "live-id", params).Return(nil)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateLiveInput{},
			expectErr: exception.ErrInvalidArgument,
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
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(nil, assert.AnError)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to list lives",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, assert.AnError)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(nil, assert.AnError)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to unmatch products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(entity.Products{}, nil)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.db.Schedule.EXPECT().Get(gomock.Any(), "schedule-id").Return(schedule, nil)
				mocks.db.Live.EXPECT().List(gomock.Any(), livesIn).Return(lives, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
				mocks.db.Live.EXPECT().Update(ctx, "live-id", params).Return(assert.AnError)
			},
			input: &store.UpdateLiveInput{
				LiveID:     "live-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Delete(ctx, "live-id").Return(assert.AnError)
			},
			input: &store.DeleteLiveInput{
				LiveID: "live-id",
			},
			expectErr: exception.ErrInternal,
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
