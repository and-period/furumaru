package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
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

func TestListSchedules(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	params := &database.ListSchedulesParams{
		Limit:  20,
		Offset: 0,
	}
	schedules := entity.Schedules{
		{
			ID:              "schedule-id",
			CoordinatorID:   "coordinator-id",
			ShippingID:      "shipping-id",
			Status:          entity.ScheduleStatusLive,
			Title:           "&.マルシェ",
			Description:     "&.マルシェの開催内容です。",
			ThumbnailURL:    "https://and-period.jp/thumbnail.png",
			ImageURL:        "https://and-period.jp/image.png",
			OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
			Public:          true,
			Approved:        true,
			ApprovedAdminID: "admin-id",
			StartAt:         now.AddDate(0, -1, 0),
			EndAt:           now.AddDate(0, 1, 0),
			CreatedAt:       now,
			UpdatedAt:       now,
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
	coordinator := &uentity.Coordinator{
		AdminID:  "coordinator-id",
		Username: "&.コーディネータ",
	}
	shipping := &entity.Shipping{
		ID:   "shipping-id",
		Name: "デフォルト配送設定",
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
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), "shipping-id").Return(shipping, nil)
				mocks.db.Schedule.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, schedule *entity.Schedule) error {
						expect := &entity.Schedule{
							ID:              schedule.ID, // ignore
							CoordinatorID:   "coordinator-id",
							ShippingID:      "shipping-id",
							Title:           "タイトル",
							Description:     "説明",
							ThumbnailURL:    "https://and-period.jp/thumbnail.png",
							ImageURL:        "https://and-period.jp/image.png",
							OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
							Public:          true,
							StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
							EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
						}
						assert.Equal(t, expect, schedule)
						return nil
					})
				mocks.media.EXPECT().ResizeScheduleThumbnail(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID:   "coordinator-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
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
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), "shipping-id").Return(shipping, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID:   "coordinator-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, exception.ErrNotFound)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), "shipping-id").Return(shipping, nil)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID:   "coordinator-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), "shipping-id").Return(nil, assert.AnError)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID:   "coordinator-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to not found shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), "shipping-id").Return(nil, exception.ErrNotFound)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID:   "coordinator-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.db.Shipping.EXPECT().Get(gomock.Any(), "shipping-id").Return(shipping, nil)
				mocks.db.Schedule.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateScheduleInput{
				CoordinatorID:   "coordinator-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
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

func TestUpdateSchedule(t *testing.T) {
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
	params := &database.UpdateScheduleParams{
		ShippingID:      "shipping-id",
		Title:           "タイトル",
		Description:     "説明",
		ThumbnailURL:    "https://and-period.jp/thumbnail.png",
		ImageURL:        "https://and-period.jp/image.png",
		OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
		Public:          true,
		StartAt:         now.AddDate(0, -1, 0),
		EndAt:           now.AddDate(0, 1, 0),
	}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.UpdateScheduleInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := *params
				params.ThumbnailURL = "https://tmp.and-period.jp/thumbnail.png"
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.db.Shipping.EXPECT().Get(ctx, "shipping-id").Return(&entity.Shipping{}, nil)
				mocks.db.Schedule.EXPECT().Update(ctx, "schedule-id", &params).Return(nil)
				mocks.media.EXPECT().ResizeScheduleThumbnail(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &store.UpdateScheduleInput{
				ScheduleID:      "schedule-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://tmp.and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.UpdateScheduleInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &store.UpdateScheduleInput{
				ScheduleID:      "schedule-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
			},
			expect: exception.ErrUnknown,
		},
		{
			name: "not found shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.db.Shipping.EXPECT().Get(ctx, "shipping-id").Return(nil, exception.ErrNotFound)
			},
			input: &store.UpdateScheduleInput{
				ScheduleID:      "schedule-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
			},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.db.Shipping.EXPECT().Get(ctx, "shipping-id").Return(nil, assert.AnError)
			},
			input: &store.UpdateScheduleInput{
				ScheduleID:      "schedule-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
			},
			expect: exception.ErrUnknown,
		},
		{
			name: "failed to update schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.db.Shipping.EXPECT().Get(ctx, "shipping-id").Return(&entity.Shipping{}, nil)
				mocks.db.Schedule.EXPECT().Update(ctx, "schedule-id", params).Return(assert.AnError)
			},
			input: &store.UpdateScheduleInput{
				ScheduleID:      "schedule-id",
				ShippingID:      "shipping-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
			},
			expect: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestUpdateScheduleThumbnails(t *testing.T) {
	t.Parallel()

	thumbnails := common.Images{
		{
			Size: common.ImageSizeSmall,
			URL:  "https://and-period.jp/thumbnail_240.png",
		},
		{
			Size: common.ImageSizeMedium,
			URL:  "https://and-period.jp/thumbnail_675.png",
		},
		{
			Size: common.ImageSizeLarge,
			URL:  "https://and-period.jp/thumbnail_900.png",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateScheduleThumbnailsInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpdateThumbnails(ctx, "schedule-id", thumbnails).Return(nil)
			},
			input: &store.UpdateScheduleThumbnailsInput{
				ScheduleID: "schedule-id",
				Thumbnails: thumbnails,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateScheduleThumbnailsInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update thumbnails",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpdateThumbnails(ctx, "schedule-id", thumbnails).Return(assert.AnError)
			},
			input: &store.UpdateScheduleThumbnailsInput{
				ScheduleID: "schedule-id",
				Thumbnails: thumbnails,
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateScheduleThumbnails(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
