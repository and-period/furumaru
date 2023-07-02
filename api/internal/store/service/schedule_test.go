package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSchedules(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	params := &database.ListSchedulesParams{
		Limit:  20,
		Offset: 0,
	}
	schedules := entity.Schedules{
		{
			ID:                   "schedule-id",
			CoordinatorID:        "coordinator-id",
			ShippingID:           "shipping-id",
			Status:               entity.ScheduleStatusLive,
			Title:                "&.マルシェ",
			Description:          "&.マルシェの開催内容です。",
			ThumbnailURL:         "https://and-period.jp/thumbnail.png",
			OpeningVideoURL:      "https://and-period.jp/opening-video.mp4",
			IntermissionVideoURL: "https://and-period.jp/intermission-video.mp4",
			Public:               true,
			Approved:             true,
			ApprovedAdminID:      "admin-id",
			StartAt:              now.AddDate(0, -1, 0),
			EndAt:                now.AddDate(0, 1, 0),
			CreatedAt:            now,
			UpdatedAt:            now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListSchedulesInput
		expect      entity.Schedules
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(gomock.Any(), params).Return(schedules, nil)
				mocks.db.Schedule.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListSchedulesInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      schedules,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListSchedulesInput{
				Limit:  0,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list schedules",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Schedule.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListSchedulesInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count schedules",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(gomock.Any(), params).Return(schedules, nil)
				mocks.db.Schedule.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListSchedulesInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListSchedules(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	now := time.Date(2022, time.January, 3, 18, 30, 0, 0, time.UTC)
	schedule := &entity.Schedule{
		ID:            "schedule-id",
		CoordinatorID: "coordinator-id",
		ShippingID:    "shipping-id",
		Title:         "タイトル",
		Description:   "説明",
		ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
		StartAt:       now,
		EndAt:         now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetScheduleInput
		expect    *entity.Schedule
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
			},
			input: &store.GetScheduleInput{
				ScheduleID: "schedule-id",
			},
			expect:    schedule,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetScheduleInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &store.GetScheduleInput{
				ScheduleID: "schedule-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateSchedule(t *testing.T) {
	t.Parallel()

	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	producersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id01"},
	}
	shippingIn := "shipping-id"

	coordinator := &uentity.Coordinator{
		AdminID: "coordinator-id",
	}
	producers := uentity.Producers{
		{AdminID: "producer-id01"},
	}
	shipping := &entity.Shipping{
		ID: "shipping-id",
	}
	products := entity.Products{
		{ID: "product-id"},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateScheduleInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
				mocks.db.Schedule.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, s *entity.Schedule, ls entity.Lives, ps entity.LiveProducts) error {
						eschedule := &entity.Schedule{
							ID:            s.ID, // ignore
							CoordinatorID: "coordinator-id",
							ShippingID:    "shipping-id",
							Title:         "タイトル",
							Description:   "説明",
							ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
							StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
							EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
						}
						assert.Equal(t, eschedule, s)
						require.Len(t, ls, 1)
						elives := entity.Lives{{
							ID:          ls[0].ID, // ignore
							ScheduleID:  s.ID,
							ProducerID:  "producer-id01",
							Title:       "配信タイトル",
							Description: "配信の説明",
							StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
							EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
						}}
						assert.Equal(t, elives, ls)
						eproducts := entity.LiveProducts{
							{
								LiveID:    ls[0].ID,
								ProductID: "product-id",
							},
						}
						assert.Equal(t, eproducts, ps)
						return nil
					})
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: nil,
		},
		{
			name:      "success",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateScheduleInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, exception.ErrNotFound)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(nil, assert.AnError)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unmatch producers length",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(uentity.Producers{}, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(nil, assert.AnError)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to not found shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(nil, exception.ErrNotFound)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrNotFound,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(nil, assert.AnError)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unmatch products length",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(entity.Products{}, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(products, nil)
				mocks.db.Schedule.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "タイトル",
				Description:   "説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:         jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						ProductIDs:  []string{"product-id"},
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
					},
				},
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, _, err := service.CreateSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
