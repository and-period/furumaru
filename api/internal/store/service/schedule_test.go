package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSchedule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateScheduleInput
		expect    *entity.Schedule
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, schedule *entity.Schedule) error {
						expect := &entity.Schedule{
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
			},
			input: &store.CreateScheduleInput{
				Title:        "タイトル",
				Description:  "説明",
				ThumbnailURL: "https://and-period.jp/thumbnail01.png",
				StartAt:      jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:        jst.Date(2022, 1, 3, 18, 30, 0, 0),
				Lives:        []*store.CreateScheduleLive{},
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
			_, err := service.CreateSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
