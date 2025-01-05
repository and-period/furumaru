package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExperience(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewExperience(nil))
}

func TestExperience_List(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 2)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)
	internal := make(internalExperiences, 3)
	internal[0] = testExperience("experience-id01", "experience-type-id01", "coordinator-id", "producer-id", 1, now())
	internal[0].StartAt = now().AddDate(0, 0, -1)
	internal[1] = testExperience("experience-id02", "experience-type-id02", "coordinator-id", "producer-id", 2, now())
	internal[1].StartAt = now().AddDate(0, 0, -2)
	internal[2] = testExperience("experience-id03", "experience-type-id02", "coordinator-id", "producer-id", 3, now())
	internal[2].StartAt = now().AddDate(0, -1, 0)
	err = db.DB.Table(experienceTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ExperienceRevision).Error
		require.NoError(t, err)
	}
	experiences, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		params *database.ListExperiencesParams
	}
	type want struct {
		experiences entity.Experiences
		err         error
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
				params: &database.ListExperiencesParams{
					Name:   "収穫",
					Limit:  2,
					Offset: 1,
				},
			},
			want: want{
				experiences: experiences[1:],
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.experiences, actual)
		})
	}
}

func TestExperience_ListByGeolocation(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 2)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)
	internal := make(internalExperiences, 3)
	internal[0] = testExperience("experience-id01", "experience-type-id01", "coordinator-id", "producer-id", 1, now())
	internal[0].StartAt = now().AddDate(0, 0, -1)
	internal[1] = testExperience("experience-id02", "experience-type-id02", "coordinator-id", "producer-id", 2, now())
	internal[1].StartAt = now().AddDate(0, 0, -2)
	internal[2] = testExperience("experience-id03", "experience-type-id02", "coordinator-id", "producer-id", 3, now())
	internal[2].StartAt = now().AddDate(0, -1, 0)
	err = db.DB.Table(experienceTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ExperienceRevision).Error
		require.NoError(t, err)
	}
	experiences, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		params *database.ListExperiencesByGeolocationParams
	}
	type want struct {
		experiences entity.Experiences
		err         error
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
				params: &database.ListExperiencesByGeolocationParams{
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					Latitude:      35.276833,
					Longitude:     136.251739,
					Radius:        1000,
				},
			},
			want: want{
				experiences: experiences,
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			actual, err := db.ListByGeolocation(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.experiences, actual)
		})
	}
}

func TestExperience_Count(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 2)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)
	experiences := make(internalExperiences, 3)
	experiences[0] = testExperience("experience-id01", "experience-type-id01", "coordinator-id", "producer-id", 1, now())
	experiences[1] = testExperience("experience-id02", "experience-type-id02", "coordinator-id", "producer-id", 2, now())
	experiences[2] = testExperience("experience-id03", "experience-type-id02", "coordinator-id", "producer-id", 3, now())
	err = db.DB.Table(experienceTable).Create(&experiences).Error
	require.NoError(t, err)
	for i := range experiences {
		err := db.DB.Create(&experiences[i].ExperienceRevision).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.ListExperiencesParams
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
				params: &database.ListExperiencesParams{
					Name:   "収穫",
					Limit:  2,
					Offset: 1,
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

			db := &experience{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestExperience_MultiGet(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 2)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)
	internal := make(internalExperiences, 3)
	internal[0] = testExperience("experience-id01", "experience-type-id01", "coordinator-id", "producer-id", 1, now())
	internal[0].StartAt = now().AddDate(0, 0, -1)
	internal[1] = testExperience("experience-id02", "experience-type-id02", "coordinator-id", "producer-id", 2, now())
	internal[1].StartAt = now().AddDate(0, 0, -2)
	internal[2] = testExperience("experience-id03", "experience-type-id02", "coordinator-id", "producer-id", 3, now())
	internal[2].StartAt = now().AddDate(0, -1, 0)
	err = db.DB.Table(experienceTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ExperienceRevision).Error
		require.NoError(t, err)
	}
	experiences, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		experienceIDs []string
	}
	type want struct {
		experiences entity.Experiences
		err         error
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
				experienceIDs: []string{"experience-id01", "experience-id02"},
			},
			want: want{
				experiences: experiences[:2],
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.experienceIDs)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.experiences, actual)
		})
	}
}

func TestExperience_MultiGetByRevision(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 2)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)
	internal := make(internalExperiences, 3)
	internal[0] = testExperience("experience-id01", "experience-type-id01", "coordinator-id", "producer-id", 1, now())
	internal[0].StartAt = now().AddDate(0, 0, -1)
	internal[1] = testExperience("experience-id02", "experience-type-id02", "coordinator-id", "producer-id", 2, now())
	internal[1].StartAt = now().AddDate(0, 0, -2)
	internal[2] = testExperience("experience-id03", "experience-type-id02", "coordinator-id", "producer-id", 3, now())
	internal[2].StartAt = now().AddDate(0, -1, 0)
	err = db.DB.Table(experienceTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ExperienceRevision).Error
		require.NoError(t, err)
	}
	experiences, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		revisionIDs []int64
	}
	type want struct {
		experiences entity.Experiences
		err         error
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
				revisionIDs: []int64{1, 2, 3},
			},
			want: want{
				experiences: experiences,
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			actual, err := db.MultiGetByRevision(ctx, tt.args.revisionIDs)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.experiences, actual)
		})
	}
}

func TestExperience_Get(t *testing.T) {
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

	typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
	err = db.DB.Create(&typ).Error
	require.NoError(t, err)
	internal := testExperience("experience-id", "experience-type-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&internal).Error
	require.NoError(t, err)
	err = db.DB.Create(&internal.ExperienceRevision).Error
	require.NoError(t, err)
	e, err := internal.entity()

	type args struct {
		experienceID string
	}
	type want struct {
		experience *entity.Experience
		err        error
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
				experience: e,
				err:        nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				experienceID: "",
			},
			want: want{
				experience: nil,
				err:        database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.experienceID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.experience, actual)
		})
	}
}

func TestExperience_Create(t *testing.T) {
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

	typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
	err = db.DB.Create(&typ).Error
	require.NoError(t, err)
	internal := testExperience("experience-id", "experience-type-id", "coordinator-id", "producer-id", 1, now())
	e, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		experience *entity.Experience
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
				experience: e,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				e := testExperience("experience-id", "experience-type-id", "coordinator-id", "producer-id", 1, now())
				err := db.DB.Table(experienceTable).Create(&e).Error
				require.NoError(t, err)
				err = db.DB.Create(&e.ExperienceRevision).Error
				require.NoError(t, err)
			},
			args: args{
				experience: e,
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

			err := delete(ctx, experienceRevisionTable, experienceTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			err = db.Create(ctx, tt.args.experience)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperience_Update(t *testing.T) {
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

	typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
	err = db.DB.Create(&typ).Error
	require.NoError(t, err)

	type args struct {
		experienceID string
		params       *database.UpdateExperienceParams
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
				e := testExperience("experience-id", "experience-type-id", "coordinator-id", "producer-id", 1, now())
				err := db.DB.Table(experienceTable).Create(&e).Error
				require.NoError(t, err)
				err = db.DB.Create(&e.ExperienceRevision).Error
				require.NoError(t, err)
			},
			args: args{
				experienceID: "experience-id",
				params: &database.UpdateExperienceParams{
					TypeID:                "experience-type-id",
					Title:                 "じゃがいも収穫",
					Description:           "じゃがいもを収穫する体験です。",
					Public:                true,
					SoldOut:               true,
					Media:                 entity.MultiExperienceMedia{},
					PriceAdult:            1000,
					PriceJuniorHighSchool: 800,
					PriceElementarySchool: 600,
					PricePreschool:        400,
					PriceSenior:           500,
					RecommendedPoints: []string{
						"じゃがいもを収穫する",
						"じゃがいもを食べる",
						"じゃがいもを持ち帰る",
					},
					Duration:           60,
					Direction:          "彦根駅から徒歩10分",
					BusinessOpenTime:   "1000",
					BusinessCloseTime:  "1800",
					HostPrefectureCode: 25,
					HostCity:           "彦根市",
					StartAt:            now().AddDate(0, -1, 0),
					EndAt:              now().AddDate(0, 1, 0),
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

			err := delete(ctx, experienceRevisionTable, experienceTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			err = db.Update(ctx, tt.args.experienceID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperience_Delete(t *testing.T) {
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

	typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
	err = db.DB.Create(&typ).Error
	require.NoError(t, err)

	type args struct {
		experienceID string
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
				e := testExperience("experience-id", "experience-type-id", "coordinator-id", "producer-id", 1, now())
				err := db.DB.Table(experienceTable).Create(&e).Error
				require.NoError(t, err)
				err = db.DB.Create(&e.ExperienceRevision).Error
				require.NoError(t, err)
			},
			args: args{
				experienceID: "experience-id",
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

			err := delete(ctx, experienceRevisionTable, experienceTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experience{db: db, now: now}
			err = db.Delete(ctx, tt.args.experienceID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testExperience(experienceID, typeID, coordinatorID, producerID string, revisionID int64, now time.Time) *internalExperience {
	experience := &entity.Experience{
		ID:            experienceID,
		CoordinatorID: coordinatorID,
		ProducerID:    producerID,
		TypeID:        typeID,
		Title:         "じゃがいも収穫",
		Description:   "じゃがいもを収穫する体験です。",
		Public:        true,
		SoldOut:       false,
		Status:        entity.ExperienceStatusAccepting,
		ThumbnailURL:  "http://example.com/thumbnail01.png",
		Media: entity.MultiExperienceMedia{
			{
				URL:         "http://example.com/thumbnail01.png",
				IsThumbnail: true,
			},
			{
				URL:         "http://example.com/thumbnail02.png",
				IsThumbnail: false,
			},
		},
		RecommendedPoints: []string{
			"じゃがいもを収穫する",
			"じゃがいもを食べる",
			"じゃがいもを持ち帰る",
		},
		PromotionVideoURL:  "http://example.com/promotion.mp4",
		Duration:           60,
		Direction:          "彦根駅から徒歩10分",
		BusinessOpenTime:   "1000",
		BusinessCloseTime:  "1800",
		HostPostalCode:     "5220061",
		HostPrefecture:     "滋賀県",
		HostPrefectureCode: 25,
		HostCity:           "彦根市",
		HostAddressLine1:   "金亀町１−１",
		HostAddressLine2:   "",
		HostLongitude:      136.251739,
		HostLatitude:       35.276833,
		StartAt:            now.AddDate(0, -1, 0),
		EndAt:              now.AddDate(0, 1, 0),
		ExperienceRevision: *testExperienceRevision(revisionID, experienceID, now),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	internal, _ := newInternalExperience(experience)
	return internal
}

func testExperienceRevision(revisionID int64, experienceID string, now time.Time) *entity.ExperienceRevision {
	return &entity.ExperienceRevision{
		ID:                    revisionID,
		ExperienceID:          experienceID,
		PriceAdult:            1000,
		PriceJuniorHighSchool: 800,
		PriceElementarySchool: 600,
		PricePreschool:        400,
		PriceSenior:           200,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}
