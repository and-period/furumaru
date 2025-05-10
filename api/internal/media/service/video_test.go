package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListVideos(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	params := &database.ListVideosParams{
		Name:                  "じゃがいも収穫",
		CoordinatorID:         "coordinator-id",
		OnlyPublished:         true,
		OnlyDisplayProduct:    false,
		OnlyDisplayExperience: false,
		ExcludeLimited:        false,
		Limit:                 20,
		Offset:                0,
	}
	videos := entity.Videos{
		{
			ID:            "video-id",
			CoordinatorID: "coordinator-id",
			ProductIDs:    []string{"product-id"},
			ExperienceIDs: []string{"experience-id"},
			Title:         "じゃがいも収穫",
			Description:   "じゃがいも収穫の説明",
			Status:        entity.VideoStatusPublished,
			ThumbnailURL:  "https://example.com/thumbnail.jpg",
			VideoURL:      "https://example.com/video.mp4",
			Public:        true,
			Limited:       false,
			VideoProducts: []*entity.VideoProduct{{
				VideoID:   "video-id",
				ProductID: "product-id",
				Priority:  1,
				CreatedAt: now,
				UpdatedAt: now,
			}},
			VideoExperiences: []*entity.VideoExperience{{
				VideoID:      "video-id",
				ExperienceID: "experience-id",
				Priority:     1,
				CreatedAt:    now,
				UpdatedAt:    now,
			}},
			PublishedAt: now.AddDate(0, 0, -1),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *media.ListVideosInput
		expect      entity.Videos
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().List(gomock.Any(), params).Return(videos, nil)
				mocks.db.Video.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &media.ListVideosInput{
				Name:                  "じゃがいも収穫",
				CoordinatorID:         "coordinator-id",
				OnlyPublished:         true,
				OnlyDisplayProduct:    false,
				OnlyDisplayExperience: false,
				ExcludeLimited:        false,
				Limit:                 20,
				Offset:                0,
				NoLimit:               false,
			},
			expect:      videos,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &media.ListVideosInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list videos",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Video.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &media.ListVideosInput{
				Name:                  "じゃがいも収穫",
				CoordinatorID:         "coordinator-id",
				OnlyPublished:         true,
				OnlyDisplayProduct:    false,
				OnlyDisplayExperience: false,
				ExcludeLimited:        false,
				Limit:                 20,
				Offset:                0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count videos",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().List(gomock.Any(), params).Return(videos, nil)
				mocks.db.Video.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &media.ListVideosInput{
				Name:                  "じゃがいも収穫",
				CoordinatorID:         "coordinator-id",
				OnlyPublished:         true,
				OnlyDisplayProduct:    false,
				OnlyDisplayExperience: false,
				ExcludeLimited:        false,
				Limit:                 20,
				Offset:                0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListVideos(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestListProductVideos(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	videos := entity.Videos{
		{
			ID:            "video-id",
			CoordinatorID: "coordinator-id",
			ProductIDs:    []string{"product-id"},
			ExperienceIDs: []string{"experience-id"},
			Title:         "じゃがいも収穫",
			Description:   "じゃがいも収穫の説明",
			Status:        entity.VideoStatusPublished,
			ThumbnailURL:  "https://example.com/thumbnail.jpg",
			VideoURL:      "https://example.com/video.mp4",
			Public:        true,
			Limited:       false,
			VideoProducts: []*entity.VideoProduct{{
				VideoID:   "video-id",
				ProductID: "product-id",
				Priority:  1,
				CreatedAt: now,
				UpdatedAt: now,
			}},
			VideoExperiences: []*entity.VideoExperience{{
				VideoID:      "video-id",
				ExperienceID: "experience-id",
				Priority:     1,
				CreatedAt:    now,
				UpdatedAt:    now,
			}},
			PublishedAt: now.AddDate(0, 0, -1),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.ListProductVideosInput
		expect    entity.Videos
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().ListByProductID(gomock.Any(), "product-id").Return(videos, nil)
			},
			input: &media.ListProductVideosInput{
				ProductID: "product-id",
			},
			expect:    videos,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.ListProductVideosInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list videos",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().ListByProductID(gomock.Any(), "product-id").Return(nil, assert.AnError)
			},
			input: &media.ListProductVideosInput{
				ProductID: "product-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListProductVideos(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestListExperienceVideos(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	videos := entity.Videos{
		{
			ID:            "video-id",
			CoordinatorID: "coordinator-id",
			ProductIDs:    []string{"product-id"},
			ExperienceIDs: []string{"experience-id"},
			Title:         "じゃがいも収穫",
			Description:   "じゃがいも収穫の説明",
			Status:        entity.VideoStatusPublished,
			ThumbnailURL:  "https://example.com/thumbnail.jpg",
			VideoURL:      "https://example.com/video.mp4",
			Public:        true,
			Limited:       false,
			VideoProducts: []*entity.VideoProduct{{
				VideoID:   "video-id",
				ProductID: "product-id",
				Priority:  1,
				CreatedAt: now,
				UpdatedAt: now,
			}},
			VideoExperiences: []*entity.VideoExperience{{
				VideoID:      "video-id",
				ExperienceID: "experience-id",
				Priority:     1,
				CreatedAt:    now,
				UpdatedAt:    now,
			}},
			PublishedAt: now.AddDate(0, 0, -1),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.ListExperienceVideosInput
		expect    entity.Videos
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().ListByExperienceID(gomock.Any(), "experience-id").Return(videos, nil)
			},
			input: &media.ListExperienceVideosInput{
				ExperienceID: "experience-id",
			},
			expect:    videos,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.ListExperienceVideosInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list videos",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().ListByExperienceID(gomock.Any(), "experience-id").Return(nil, assert.AnError)
			},
			input: &media.ListExperienceVideosInput{
				ExperienceID: "experience-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListExperienceVideos(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetVideo(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	video := &entity.Video{
		ID:            "video-id",
		CoordinatorID: "coordinator-id",
		ProductIDs:    []string{"product-id"},
		ExperienceIDs: []string{"experience-id"},
		Title:         "じゃがいも収穫",
		Description:   "じゃがいも収穫の説明",
		Status:        entity.VideoStatusPublished,
		ThumbnailURL:  "https://example.com/thumbnail.jpg",
		VideoURL:      "https://example.com/video.mp4",
		Public:        true,
		Limited:       false,
		VideoProducts: []*entity.VideoProduct{{
			VideoID:   "video-id",
			ProductID: "product-id",
			Priority:  1,
			CreatedAt: now,
			UpdatedAt: now,
		}},
		VideoExperiences: []*entity.VideoExperience{{
			VideoID:      "video-id",
			ExperienceID: "experience-id",
			Priority:     1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
		PublishedAt: now.AddDate(0, 0, -1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.GetVideoInput
		expect    *entity.Video
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(gomock.Any(), "video-id").Return(video, nil)
			},
			input: &media.GetVideoInput{
				VideoID: "video-id",
			},
			expect:    video,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.GetVideoInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get video",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(gomock.Any(), "video-id").Return(nil, assert.AnError)
			},
			input: &media.GetVideoInput{
				VideoID: "video-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateVideo(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		AdminID: "coordinator-id",
	}
	productIn := &store.MultiGetProductsInput{
		ProductIDs: []string{"product-id"},
	}
	products := sentity.Products{
		{ID: "product-id"},
	}
	experienceIn := &store.MultiGetExperiencesInput{
		ExperienceIDs: []string{"experience-id"},
	}
	experiences := sentity.Experiences{
		{ID: "experience-id"},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.CreateVideoInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
				mocks.db.Video.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, video *entity.Video) error {
						expect := &entity.Video{
							ID:                video.ID, // ignore
							CoordinatorID:     "coordinator-id",
							ProductIDs:        []string{"product-id"},
							ExperienceIDs:     []string{"experience-id"},
							Title:             "オンデマンド配信",
							Description:       "オンデマンド配信の説明",
							Status:            entity.VideoStatusUnknown,
							ThumbnailURL:      "https://example.com/thumbnail.jpg",
							VideoURL:          "https://example.com/video.mp4",
							Public:            true,
							Limited:           false,
							DisplayProduct:    true,
							DisplayExperience: true,
							PublishedAt:       now,
							VideoProducts: []*entity.VideoProduct{{
								VideoID:   video.ID,
								ProductID: "product-id",
								Priority:  1,
							}},
							VideoExperiences: []*entity.VideoExperience{{
								VideoID:      video.ID,
								ExperienceID: "experience-id",
								Priority:     1,
							}},
						}
						assert.Equal(t, expect, video)
						return nil
					})
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.CreateVideoInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, exception.ErrNotFound)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "unmatch products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(sentity.Products{}, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "unmatch experiences",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(sentity.Experiences{}, nil)
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get experiences",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(nil, assert.AnError)
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create video",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
				mocks.db.Video.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &media.CreateVideoInput{
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				CoordinatorID:     "coordinator-id",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateVideo(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	productIn := &store.MultiGetProductsInput{
		ProductIDs: []string{"product-id"},
	}
	products := sentity.Products{
		{ID: "product-id"},
	}
	experienceIn := &store.MultiGetExperiencesInput{
		ExperienceIDs: []string{"experience-id"},
	}
	experiences := sentity.Experiences{
		{ID: "experience-id"},
	}
	params := &database.UpdateVideoParams{
		Title:             "オンデマンド配信",
		Description:       "オンデマンド配信の説明",
		ProductIDs:        []string{"product-id"},
		ExperienceIDs:     []string{"experience-id"},
		ThumbnailURL:      "https://example.com/thumbnail.jpg",
		VideoURL:          "https://example.com/video.mp4",
		Public:            true,
		Limited:           false,
		DisplayProduct:    true,
		DisplayExperience: true,
		PublishedAt:       now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UpdateVideoInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
				mocks.db.Video.EXPECT().Update(ctx, "video-id", params).Return(nil)
			},
			input: &media.UpdateVideoInput{
				VideoID:           "video-id",
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UpdateVideoInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "unmatch products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(sentity.Products{}, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
			},
			input: &media.UpdateVideoInput{
				VideoID:           "video-id",
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "unmatch experiences",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(sentity.Experiences{}, nil)
			},
			input: &media.UpdateVideoInput{
				VideoID:           "video-id",
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
			},
			input: &media.UpdateVideoInput{
				VideoID:           "video-id",
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get experiences",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(nil, assert.AnError)
			},
			input: &media.UpdateVideoInput{
				VideoID:           "video-id",
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to update video",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productIn).Return(products, nil)
				mocks.store.EXPECT().MultiGetExperiences(gomock.Any(), experienceIn).Return(experiences, nil)
				mocks.db.Video.EXPECT().Update(ctx, "video-id", params).Return(assert.AnError)
			},
			input: &media.UpdateVideoInput{
				VideoID:           "video-id",
				Title:             "オンデマンド配信",
				Description:       "オンデマンド配信の説明",
				ProductIDs:        []string{"product-id"},
				ExperienceIDs:     []string{"experience-id"},
				ThumbnailURL:      "https://example.com/thumbnail.jpg",
				VideoURL:          "https://example.com/video.mp4",
				Public:            true,
				Limited:           false,
				DisplayProduct:    true,
				DisplayExperience: true,
				PublishedAt:       now,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteVideo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.DeleteVideoInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Delete(ctx, "video-id").Return(nil)
			},
			input: &media.DeleteVideoInput{
				VideoID: "video-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.DeleteVideoInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete video",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Delete(ctx, "video-id").Return(assert.AnError)
			},
			input: &media.DeleteVideoInput{
				VideoID: "video-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
