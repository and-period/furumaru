package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
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
		ShopID: "shop-id",
		Limit:  20,
		Offset: 0,
	}
	schedules := entity.Schedules{
		{
			ID:              "schedule-id",
			ShopID:          "shop-id",
			CoordinatorID:   "coordinator-id",
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
				ShopID: "shop-id",
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
				ShopID: "shop-id",
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count schedules",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(gomock.Any(), params).Return(schedules, nil)
				mocks.db.Schedule.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListSchedulesInput{
				ShopID: "shop-id",
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListSchedules(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetSchedules(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 7, 1, 18, 30, 0, 0)
	schedules := entity.Schedules{
		{
			ID:              "schedule-id",
			CoordinatorID:   "coordinator-id",
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
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetSchedulesInput
		expect    entity.Schedules
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().MultiGet(gomock.Any(), []string{"schedule-id"}).Return(schedules, nil)
			},
			input: &store.MultiGetSchedulesInput{
				ScheduleIDs: []string{"schedule-id"},
			},
			expect:    schedules,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetSchedulesInput{
				ScheduleIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list schedules",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().MultiGet(gomock.Any(), []string{"schedule-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetSchedulesInput{
				ScheduleIDs: []string{"schedule-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetSchedules(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	now := time.Date(2022, time.January, 3, 18, 30, 0, 0, time.UTC)
	schedule := &entity.Schedule{
		ID:            "schedule-id",
		CoordinatorID: "coordinator-id",
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
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateSchedule(t *testing.T) {
	t.Parallel()

	shop := &entity.Shop{
		ID:            "shop-id",
		Name:          "じゃがいも農園",
		CoordinatorID: "coordinator-id",
		ProducerIDs:   []string{"producer-id"},
	}
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		AdminID:  "coordinator-id",
		Username: "&.コーディネータ",
		Admin:    uentity.Admin{ID: "coordinator-id"},
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
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(coordinator, nil)
				mocks.db.Schedule.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, schedule *entity.Schedule) error {
						expect := &entity.Schedule{
							ID:              schedule.ID, // ignore
							ShopID:          "shop-id",
							CoordinatorID:   "coordinator-id",
							Title:           "タイトル",
							Description:     "説明",
							ThumbnailURL:    "https://and-period.jp/thumbnail.png",
							ImageURL:        "https://and-period.jp/image.png",
							OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
							Public:          true,
							Approved:        true,
							StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
							EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
						}
						assert.Equal(t, expect, schedule)
						return nil
					})
				mocks.media.EXPECT().CreateBroadcast(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				mocks.messenger.EXPECT().ReserveStartLive(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateScheduleInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
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
			name: "failed to get shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(nil, assert.AnError)
			},
			input: &store.CreateScheduleInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to not found shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(nil, database.ErrNotFound)
			},
			input: &store.CreateScheduleInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
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
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(nil, assert.AnError)
			},
			input: &store.CreateScheduleInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(nil, exception.ErrNotFound)
			},
			input: &store.CreateScheduleInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
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
			name: "unmatch coordinator id",
			setup: func(ctx context.Context, mocks *mocks) {
				shop := &entity.Shop{ID: "shop-id"}
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(coordinator, nil)
			},
			input: &store.CreateScheduleInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
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
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(coordinator, nil)
				mocks.db.Schedule.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateScheduleInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://ane-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 1, 2, 18, 30, 0, 0),
				EndAt:           jst.Date(2022, 1, 3, 18, 30, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
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
		Title:         "タイトル",
		Description:   "説明",
		ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
		StartAt:       now,
		EndAt:         now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	params := &database.UpdateScheduleParams{
		Title:           "タイトル",
		Description:     "説明",
		ThumbnailURL:    "https://and-period.jp/thumbnail.png",
		ImageURL:        "https://and-period.jp/image.png",
		OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
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
				mocks.db.Schedule.EXPECT().Update(ctx, "schedule-id", &params).Return(nil)
				mocks.messenger.EXPECT().ReserveStartLive(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &store.UpdateScheduleInput{
				ScheduleID:      "schedule-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://tmp.and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
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
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to update schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.db.Schedule.EXPECT().Update(ctx, "schedule-id", params).Return(assert.AnError)
			},
			input: &store.UpdateScheduleInput{
				ScheduleID:      "schedule-id",
				Title:           "タイトル",
				Description:     "説明",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
			},
			expect: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestDeleteSchedule(t *testing.T) {
	t.Parallel()

	now := time.Date(2022, time.January, 3, 18, 30, 0, 0, time.UTC)
	schedule := &entity.Schedule{
		ID:            "schedule-id",
		CoordinatorID: "coordinator-id",
		Title:         "タイトル",
		Description:   "説明",
		ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
		StartAt:       now,
		EndAt:         now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	broadcastIn := &media.GetBroadcastByScheduleIDInput{
		ScheduleID: "schedule-id",
	}
	broadcast := &mentity.Broadcast{
		ID:         "broadcast-id",
		ScheduleID: "schedule-id",
		Status:     mentity.BroadcastStatusDisabled,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.DeleteScheduleInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.media.EXPECT().GetBroadcastByScheduleID(ctx, broadcastIn).Return(broadcast, nil)
				mocks.db.Schedule.EXPECT().Delete(ctx, "schedule-id").Return(nil)
			},
			input: &store.DeleteScheduleInput{
				ScheduleID: "schedule-id",
			},
			expect: nil,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &store.DeleteScheduleInput{
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.media.EXPECT().GetBroadcastByScheduleID(ctx, broadcastIn).Return(nil, assert.AnError)
			},
			input: &store.DeleteScheduleInput{
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "broadcast is not disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &mentity.Broadcast{Status: mentity.BroadcastStatusActive}
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.media.EXPECT().GetBroadcastByScheduleID(ctx, broadcastIn).Return(broadcast, nil)
			},
			input: &store.DeleteScheduleInput{
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to delete schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Get(ctx, "schedule-id").Return(schedule, nil)
				mocks.media.EXPECT().GetBroadcastByScheduleID(ctx, broadcastIn).Return(broadcast, nil)
				mocks.db.Schedule.EXPECT().Delete(ctx, "schedule-id").Return(assert.AnError)
			},
			input: &store.DeleteScheduleInput{
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestApproveSchedule(t *testing.T) {
	t.Parallel()

	adminIn := &user.GetAdministratorInput{
		AdministratorID: "admin-id",
	}
	admin := &uentity.Administrator{
		AdminID: "admin-id",
	}
	params := &database.ApproveScheduleParams{
		Approved:        true,
		ApprovedAdminID: "admin-id",
	}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.ApproveScheduleInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdministrator(ctx, adminIn).Return(admin, nil)
				mocks.db.Schedule.EXPECT().Approve(ctx, "schedule-id", params).Return(nil)
			},
			input: &store.ApproveScheduleInput{
				ScheduleID: "schedule-id",
				AdminID:    "admin-id",
				Approved:   true,
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.ApproveScheduleInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdministrator(ctx, adminIn).Return(nil, assert.AnError)
			},
			input: &store.ApproveScheduleInput{
				ScheduleID: "schedule-id",
				AdminID:    "admin-id",
				Approved:   true,
			},
			expect: exception.ErrInternal,
		},
		{
			name: "not found administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdministrator(ctx, adminIn).Return(nil, exception.ErrNotFound)
			},
			input: &store.ApproveScheduleInput{
				ScheduleID: "schedule-id",
				AdminID:    "admin-id",
				Approved:   true,
			},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to approve schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdministrator(ctx, adminIn).Return(admin, nil)
				mocks.db.Schedule.EXPECT().Approve(ctx, "schedule-id", params).Return(assert.AnError)
			},
			input: &store.ApproveScheduleInput{
				ScheduleID: "schedule-id",
				AdminID:    "admin-id",
				Approved:   true,
			},
			expect: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ApproveSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestPublishSchedule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.PublishScheduleInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Publish(ctx, "schedule-id", true).Return(nil)
			},
			input: &store.PublishScheduleInput{
				ScheduleID: "schedule-id",
				Public:     true,
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.PublishScheduleInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to publish schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().Publish(ctx, "schedule-id", true).Return(assert.AnError)
			},
			input: &store.PublishScheduleInput{
				ScheduleID: "schedule-id",
				Public:     true,
			},
			expect: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.PublishSchedule(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}
