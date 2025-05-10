package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestVideo_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	videos := make(entity.Videos, 3)
	videos[0] = testVideo("video-id01", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[0].PublishedAt = now().AddDate(0, 0, -1)
	videos[1] = testVideo("video-id02", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[1].PublishedAt = now().AddDate(0, 0, -2)
	videos[2] = testVideo("video-id03", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[2].PublishedAt = now().AddDate(0, -1, 0)
	err = db.DB.Create(&videos).Error
	require.NoError(t, err)
	for _, video := range videos {
		err = db.DB.Create(&video.VideoProducts).Error
		require.NoError(t, err)
		err = db.DB.Create(&video.VideoExperiences).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.ListVideosParams
	}
	type want struct {
		videos entity.Videos
		err    error
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
				params: &database.ListVideosParams{
					Name:          "オンデマンド配信",
					CoordinatorID: "coordinator-id",
					Limit:         2,
					Offset:        1,
				},
			},
			want: want{
				videos: videos[1:],
				err:    nil,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.videos, actual)
		})
	}
}

func TestVideo_ListByProductID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	videos := make(entity.Videos, 3)
	videos[0] = testVideo("video-id01", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[0].PublishedAt = now().AddDate(0, 0, -1)
	videos[1] = testVideo("video-id02", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[1].PublishedAt = now().AddDate(0, 0, -2)
	videos[2] = testVideo("video-id03", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[2].PublishedAt = now().AddDate(0, -1, 0)
	err = db.DB.Create(&videos).Error
	require.NoError(t, err)
	for _, video := range videos {
		err = db.DB.Create(&video.VideoProducts).Error
		require.NoError(t, err)
		err = db.DB.Create(&video.VideoExperiences).Error
		require.NoError(t, err)
	}

	type args struct {
		productID string
	}
	type want struct {
		videos entity.Videos
		err    error
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
				productID: "product-id",
			},
			want: want{
				videos: videos,
				err:    nil,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			actual, err := db.ListByProductID(ctx, tt.args.productID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.videos, actual)
		})
	}
}

func TestVideo_ListByExperienceID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	videos := make(entity.Videos, 3)
	videos[0] = testVideo("video-id01", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[0].PublishedAt = now().AddDate(0, 0, -1)
	videos[1] = testVideo("video-id02", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[1].PublishedAt = now().AddDate(0, 0, -2)
	videos[2] = testVideo("video-id03", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[2].PublishedAt = now().AddDate(0, -1, 0)
	err = db.DB.Create(&videos).Error
	require.NoError(t, err)
	for _, video := range videos {
		err = db.DB.Create(&video.VideoProducts).Error
		require.NoError(t, err)
		err = db.DB.Create(&video.VideoExperiences).Error
		require.NoError(t, err)
	}

	type args struct {
		experienceID string
	}
	type want struct {
		videos entity.Videos
		err    error
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
				experienceID: "experience-id",
			},
			want: want{
				videos: videos,
				err:    nil,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			actual, err := db.ListByExperienceID(ctx, tt.args.experienceID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.videos, actual)
		})
	}
}

func TestVideo_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	videos := make(entity.Videos, 3)
	videos[0] = testVideo("video-id01", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[1] = testVideo("video-id02", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	videos[2] = testVideo("video-id03", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&videos).Error
	require.NoError(t, err)
	for _, video := range videos {
		err = db.DB.Create(&video.VideoProducts).Error
		require.NoError(t, err)
		err = db.DB.Create(&video.VideoExperiences).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.ListVideosParams
	}
	type want struct {
		total int64
		err   error
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
				params: &database.ListVideosParams{
					Name:          "オンデマンド配信",
					CoordinatorID: "coordinator-id",
					Limit:         1,
					Offset:        1,
				},
			},
			want: want{
				total: 3,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestVideo_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	v := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&v).Error
	require.NoError(t, err)
	err = db.DB.Create(&v.VideoProducts).Error
	require.NoError(t, err)
	err = db.DB.Create(&v.VideoExperiences).Error
	require.NoError(t, err)

	type args struct {
		videoID string
	}
	type want struct {
		video *entity.Video
		err   error
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
				videoID: "video-id",
			},
			want: want{
				video: v,
				err:   nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				videoID: "",
			},
			want: want{
				video: nil,
				err:   database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.videoID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.video, actual)
		})
	}
}

func TestVideo_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		video *entity.Video
	}
	type want struct {
		err error
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
				video: testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				v := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
				err = db.DB.Create(&v).Error
				require.NoError(t, err)
				err = db.DB.Create(&v.VideoProducts).Error
				require.NoError(t, err)
				err = db.DB.Create(&v.VideoExperiences).Error
				require.NoError(t, err)
			},
			args: args{
				video: testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, videoProductTable, videoExperienceTable, videoTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			err = db.Create(ctx, tt.args.video)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestVideo_Update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		videoID string
		params  *database.UpdateVideoParams
	}
	type want struct {
		err error
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
				v := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
				err = db.DB.Create(&v).Error
				require.NoError(t, err)
				err = db.DB.Create(&v.VideoProducts).Error
				require.NoError(t, err)
				err = db.DB.Create(&v.VideoExperiences).Error
				require.NoError(t, err)
			},
			args: args{
				videoID: "video-id",
				params: &database.UpdateVideoParams{
					Title:         "オンデマンド配信",
					Description:   "オンデマンド配信の説明",
					ProductIDs:    []string{"product-id"},
					ExperienceIDs: []string{"experience-id"},
					ThumbnailURL:  "https://example.com/thumbnail.jpg",
					VideoURL:      "https://example.com/video.mp4",
					Public:        true,
					Limited:       true,
					PublishedAt:   now().AddDate(0, 0, -1),
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, videoProductTable, videoExperienceTable, videoTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			err = db.Update(ctx, tt.args.videoID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestVideo_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		videoID string
	}
	type want struct {
		err error
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
				v := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
				err = db.DB.Create(&v).Error
				require.NoError(t, err)
				err = db.DB.Create(&v.VideoProducts).Error
				require.NoError(t, err)
				err = db.DB.Create(&v.VideoExperiences).Error
				require.NoError(t, err)
			},
			args: args{
				videoID: "video-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, videoProductTable, videoExperienceTable, videoTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &video{db: db, now: now}
			err = db.Delete(ctx, tt.args.videoID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testVideo(videoID, coordinatorID string, productIDs, experienceIDs []string, now time.Time) *entity.Video {
	products := make(entity.VideoProducts, len(productIDs))
	for i, productID := range productIDs {
		products[i] = testVideoProduct(videoID, productID, int64(i), now)
	}
	experiences := make(entity.VideoExperiences, len(experienceIDs))
	for i, experienceID := range experienceIDs {
		experiences[i] = testVideoExperience(videoID, experienceID, int64(i), now)
	}
	return &entity.Video{
		ID:               videoID,
		CoordinatorID:    coordinatorID,
		ProductIDs:       productIDs,
		ExperienceIDs:    experienceIDs,
		Title:            "オンデマンド配信",
		Description:      "オンデマンド配信の説明",
		Status:           entity.VideoStatusPublished,
		ThumbnailURL:     "https://example.com/thumbnail.jpg",
		VideoURL:         "https://example.com/video.mp4",
		Public:           true,
		Limited:          false,
		VideoProducts:    products,
		VideoExperiences: experiences,
		PublishedAt:      now.AddDate(0, 0, -1),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func testVideoProduct(videoID, productID string, priority int64, now time.Time) *entity.VideoProduct {
	return &entity.VideoProduct{
		VideoID:   videoID,
		ProductID: productID,
		Priority:  priority,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func testVideoExperience(videoID, experienceID string, priority int64, now time.Time) *entity.VideoExperience {
	return &entity.VideoExperience{
		VideoID:      videoID,
		ExperienceID: experienceID,
		Priority:     priority,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
