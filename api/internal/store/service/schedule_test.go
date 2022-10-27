package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
							Canceled:      false,
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
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, errmock)
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
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(nil, errmock)
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
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), shippingIn).Return(nil, errmock)
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
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{"product-id"}).Return(nil, errmock)
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
				mocks.db.Schedule.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errmock)
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
