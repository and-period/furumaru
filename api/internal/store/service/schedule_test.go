package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSchedule(t *testing.T) {
	t.Parallel()

	producersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id01", "producer-id02"},
	}
	producers := uentity.Producers{
		{
			AdminID: "producer-id01",
		},
		{
			AdminID: "producer-id02",
		},
	}
	productsIn := &store.MultiGetProductsInput{
		ProductIDs: []string{"product-id01"},
	}
	products := sentity.Products{
		{
			ID: "product-id01",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateScheduleInput
		expect    *sentity.Schedule
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), productsIn).Return(products, nil)
				mocks.db.Schedule.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, schedule *sentity.Schedule) error {
						expect := &sentity.Schedule{
							ID:           schedule.ID, // ignore
							Title:        "タイトル",
							Description:  "説明",
							ThumbnailURL: "https://and-period.jp/thumbnail01.png",
							StartAt:      jst.Date(2022, 1, 2, 18, 30, 0, 0),
							EndAt:        jst.Date(2022, 1, 3, 18, 30, 0, 0),
						}
						assert.Equal(t, expect, schedule)
						return nil
					})
				mocks.db.Live.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, lives sentity.Lives) error {
						expect := sentity.Lives{
							{
								ID:          lives[0].ID, //ignore
								ScheduleID:  "schedule-id",
								Title:       "配信タイトル",
								Description: "配信の説明",
								ProducerID:  "producer-id01",
								StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
								EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
								Canceled:    false,
								Recommends:  []string{"product-id01"},
							},
							{
								ID:          lives[1].ID, //ignore
								ScheduleID:  "schedule-id",
								Title:       "配信タイトル",
								Description: "配信の説明",
								ProducerID:  "producer-id02",
								StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
								EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
								Canceled:    false,
								Recommends:  []string{"product-id01"},
							},
						}
						assert.Equal(t, expect, lives)
						return nil
					})
			},
			input: &store.CreateScheduleInput{
				Title:        "タイトル",
				Description:  "説明",
				ThumbnailURL: "https://and-period.jp/thumbnail01.png",
				StartAt:      jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:        jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives: []*store.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id01",
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
						Recommends:  []string{"product-id01"},
					},
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id02",
						StartAt:     jst.Date(2022, 1, 2, 18, 30, 0, 0),
						EndAt:       jst.Date(2022, 1, 3, 18, 30, 0, 0),
						Recommends:  []string{"product-id01"},
					},
				},
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateScheduleInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &store.CreateScheduleInput{
				Title:        "タイトル",
				Description:  "説明",
				ThumbnailURL: "https://and-period.jp/thumbnail01.png",
				StartAt:      jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:        jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives:        []*store.CreateScheduleLive{},
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
