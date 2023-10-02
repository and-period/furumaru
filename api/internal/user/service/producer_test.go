package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListProducers(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	params := &database.ListProducersParams{
		Limit:  30,
		Offset: 0,
	}
	producers := entity.Producers{
		{
			Admin: entity.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleProducer,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:       "admin-id",
			CoordinatorID: "coordinator-id",
			Username:      "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			Thumbnails:    common.Images{},
			HeaderURL:     "https://and-period.jp/header.png",
			Headers:       common.Images{},
			InstagramID:   "instagram-account",
			FacebookID:    "facebook-account",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    codes.PrefectureValues["tokyo"],
			City:          "千代田区",
			AddressLine1:  "永田町1-7-1",
			AddressLine2:  "",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *user.ListProducersInput
		expect      entity.Producers
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().List(gomock.Any(), params).Return(producers, nil)
				mocks.db.Producer.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListProducersInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      producers,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &user.ListProducersInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   user.ErrInvalidArgument,
		},
		{
			name: "failed to list producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Producer.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListProducersInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   user.ErrInternal,
		},
		{
			name: "failed to count producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().List(gomock.Any(), params).Return(producers, nil)
				mocks.db.Producer.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &user.ListProducersInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListProducers(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetProducers(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	producers := entity.Producers{
		{
			Admin: entity.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleProducer,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:       "admin-id",
			CoordinatorID: "coordinator-id",
			Username:      "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			Thumbnails:    common.Images{},
			HeaderURL:     "https://and-period.jp/header.png",
			Headers:       common.Images{},
			InstagramID:   "instagram-account",
			FacebookID:    "facebook-account",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    codes.PrefectureValues["tokyo"],
			City:          "千代田区",
			AddressLine1:  "永田町1-7-1",
			AddressLine2:  "",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetProducersInput
		expect    entity.Producers
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(producers, nil)
			},
			input: &user.MultiGetProducersInput{
				ProducerIDs: []string{"admin-id"},
			},
			expect:    producers,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetProducersInput{
				ProducerIDs: []string{""},
			},
			expect:    nil,
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to multi get producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(nil, assert.AnError)
			},
			input: &user.MultiGetProducersInput{
				ProducerIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetProducers(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetProducer(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	producer := &entity.Producer{
		Admin: entity.Admin{
			ID:            "admin-id",
			Role:          entity.AdminRoleProducer,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:       "admin-id",
		CoordinatorID: "coordinator-id",
		Username:      "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		Thumbnails:    common.Images{},
		HeaderURL:     "https://and-period.jp/header.png",
		Headers:       common.Images{},
		InstagramID:   "instagram-account",
		FacebookID:    "facebook-account",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    codes.PrefectureValues["tokyo"],
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetProducerInput
		expect    *entity.Producer
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(producer, nil)
			},
			input: &user.GetProducerInput{
				ProducerID: "admin-id",
			},
			expect:    producer,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetProducerInput{},
			expect:    nil,
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(nil, assert.AnError)
			},
			input: &user.GetProducerInput{
				ProducerID: "admin-id",
			},
			expect:    nil,
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateProducer(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	coordinator := &entity.Coordinator{
		Admin: entity.Admin{
			ID:            "admin-id",
			Role:          entity.AdminRoleCoordinator,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:      "admin-id",
		PhoneNumber:  "+819012345678",
		Username:     "&.農園",
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		Thumbnails:   common.Images{},
		HeaderURL:    "https://and-period.jp/header.png",
		Headers:      common.Images{},
		InstagramID:  "instagram-account",
		FacebookID:   "facebook-account",
		PostalCode:   "1000014",
		Prefecture:   codes.PrefectureValues["tokyo"],
		City:         "千代田区",
		AddressLine1: "永田町1-7-1",
		AddressLine2: "",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectProducer := &entity.Producer{
					Admin: entity.Admin{
						Role:          entity.AdminRoleProducer,
						Lastname:      "&.",
						Firstname:     "スタッフ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "すたっふ",
						Email:         "test-admin@and-period.jp",
					},
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					Profile:       "紹介文です。",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					InstagramID:   "instgram-id",
					FacebookID:    "facebook-id",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    codes.PrefectureValues["tokyo"],
					City:          "千代田区",
					AddressLine1:  "永田町1-7-1",
					AddressLine2:  "",
				}
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, producer *entity.Producer, auth func(ctx context.Context) error) error {
						expectProducer.ID = producer.ID
						expectProducer.AdminID = producer.ID
						expectProducer.CognitoID = producer.CognitoID
						assert.Equal(t, expectProducer, producer)
						return nil
					})
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(nil)
				mocks.media.EXPECT().ResizeProducerThumbnail(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mocks.media.EXPECT().ResizeProducerHeader(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID: "coordinator-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				Profile:       "紹介文です。",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				InstagramID:   "instgram-id",
				FacebookID:    "facebook-id",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: nil,
		},
		{
			name: "success without notify register admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID: "coordinator-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				ThumbnailURL:  "",
				HeaderURL:     "",
				InstagramID:   "instgram-id",
				FacebookID:    "facebook-id",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateProducerInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID: "coordinator-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				InstagramID:   "instgram-id",
				FacebookID:    "facebook-id",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: user.ErrInternal,
		},
		{
			name: "not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, user.ErrNotFound)
			},
			input: &user.CreateProducerInput{
				CoordinatorID: "coordinator-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				InstagramID:   "instgram-id",
				FacebookID:    "facebook-id",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID: "coordinator-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				InstagramID:   "instgram-id",
				FacebookID:    "facebook-id",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProducer(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	producer := &entity.Producer{
		Admin: entity.Admin{
			ID:            "admin-id",
			Role:          entity.AdminRoleProducer,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:       "admin-id",
		CoordinatorID: "coordinator-id",
		Username:      "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		Thumbnails:    common.Images{},
		HeaderURL:     "https://and-period.jp/header.png",
		Headers:       common.Images{},
		InstagramID:   "instagram-account",
		FacebookID:    "facebook-account",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    codes.PrefectureValues["tokyo"],
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	params := &database.UpdateProducerParams{
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		Username:      "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		InstagramID:   "instagram-account",
		FacebookID:    "facebook-account",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    codes.PrefectureValues["tokyo"],
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := *params
				params.ThumbnailURL = "https://tmp.and-period.jp/thumbnail.png"
				params.HeaderURL = "https://tmp.and-period.jp/header.png"
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.db.Producer.EXPECT().Update(ctx, "producer-id", &params).Return(nil)
				mocks.media.EXPECT().ResizeProducerThumbnail(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mocks.media.EXPECT().ResizeProducerHeader(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.UpdateProducerInput{
				ProducerID:    "producer-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				ThumbnailURL:  "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:     "https://tmp.and-period.jp/header.png",
				InstagramID:   "instagram-account",
				FacebookID:    "facebook-account",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateProducerInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(nil, assert.AnError)
			},
			input: &user.UpdateProducerInput{
				ProducerID:    "producer-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				InstagramID:   "instagram-account",
				FacebookID:    "facebook-account",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: user.ErrInternal,
		},
		{
			name: "failed to update producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.db.Producer.EXPECT().Update(ctx, "producer-id", params).Return(assert.AnError)
			},
			input: &user.UpdateProducerInput{
				ProducerID:    "producer-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				InstagramID:   "instagram-account",
				FacebookID:    "facebook-account",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    codes.PrefectureValues["tokyo"],
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProducerEmail(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	producer := &entity.Producer{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Role:          entity.AdminRoleProducer,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:       "admin-id",
		CoordinatorID: "coordinator-id",
		Username:      "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		Thumbnails:    common.Images{},
		HeaderURL:     "https://and-period.jp/header.png",
		Headers:       common.Images{},
		InstagramID:   "instagram-account",
		FacebookID:    "facebook-account",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    codes.PrefectureValues["tokyo"],
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	params := &cognito.AdminChangeEmailParams{
		Username: "cognito-id",
		Email:    "test-admin@and-period.jp",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateProducerEmailInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "producer-id", "test-admin@and-period.jp").Return(nil)
			},
			input: &user.UpdateProducerEmailInput{
				ProducerID: "producer-id",
				Email:      "test-admin@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateProducerEmailInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(nil, assert.AnError)
			},
			input: &user.UpdateProducerEmailInput{
				ProducerID: "producer-id",
				Email:      "test-admin@and-period.jp",
			},
			expectErr: user.ErrInternal,
		},
		{
			name: "failed to admin change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(assert.AnError)
			},
			input: &user.UpdateProducerEmailInput{
				ProducerID: "producer-id",
				Email:      "test-admin@and-period.jp",
			},
			expectErr: user.ErrInternal,
		},
		{
			name: "failed to update email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "producer-id", "test-admin@and-period.jp").Return(assert.AnError)
			},
			input: &user.UpdateProducerEmailInput{
				ProducerID: "producer-id",
				Email:      "test-admin@and-period.jp",
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProducerEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProducerThumbnails(t *testing.T) {
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
		input     *user.UpdateProducerThumbnailsInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateThumbnails(ctx, "producer-id", thumbnails).Return(nil)
			},
			input: &user.UpdateProducerThumbnailsInput{
				ProducerID: "producer-id",
				Thumbnails: thumbnails,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateProducerThumbnailsInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to update thumbnails",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateThumbnails(ctx, "producer-id", thumbnails).Return(assert.AnError)
			},
			input: &user.UpdateProducerThumbnailsInput{
				ProducerID: "producer-id",
				Thumbnails: thumbnails,
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProducerThumbnails(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProducerHeaders(t *testing.T) {
	t.Parallel()

	headers := common.Images{
		{
			Size: common.ImageSizeSmall,
			URL:  "https://and-period.jp/header_240.png",
		},
		{
			Size: common.ImageSizeMedium,
			URL:  "https://and-period.jp/header_675.png",
		},
		{
			Size: common.ImageSizeLarge,
			URL:  "https://and-period.jp/header_900.png",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateProducerHeadersInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateHeaders(ctx, "producer-id", headers).Return(nil)
			},
			input: &user.UpdateProducerHeadersInput{
				ProducerID: "producer-id",
				Headers:    headers,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateProducerHeadersInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to update headers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateHeaders(ctx, "producer-id", headers).Return(assert.AnError)
			},
			input: &user.UpdateProducerHeadersInput{
				ProducerID: "producer-id",
				Headers:    headers,
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProducerHeaders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestResetProducerPassword(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	producer := &entity.Producer{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Role:          entity.AdminRoleProducer,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:       "admin-id",
		CoordinatorID: "coordinator-id",
		Username:      "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		Thumbnails:    common.Images{},
		HeaderURL:     "https://and-period.jp/header.png",
		Headers:       common.Images{},
		InstagramID:   "instagram-account",
		FacebookID:    "facebook-account",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    codes.PrefectureValues["tokyo"],
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ResetProducerPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().
					AdminChangePassword(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *cognito.AdminChangePasswordParams) error {
						expect := &cognito.AdminChangePasswordParams{
							Username:  "cognito-id",
							Password:  params.Password, // ignore
							Permanent: true,
						}
						assert.Equal(t, expect, params)
						return nil
					})
				mocks.messenger.EXPECT().NotifyResetAdminPassword(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &user.ResetProducerPasswordInput{
				ProducerID: "producer-id",
			},
			expectErr: nil,
		},
		{
			name: "success without notify",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyResetAdminPassword(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.ResetProducerPasswordInput{
				ProducerID: "producer-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ResetProducerPasswordInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(nil, assert.AnError)
			},
			input: &user.ResetProducerPasswordInput{
				ProducerID: "producer-id",
			},
			expectErr: user.ErrInternal,
		},
		{
			name: "failed to admin change password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &user.ResetProducerPasswordInput{
				ProducerID: "producer-id",
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResetProducerPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestRelateProducers(t *testing.T) {
	t.Parallel()

	producers := entity.Producers{{
		AdminID: "producer-id",
	}}
	coordinator := &entity.Coordinator{
		AdminID: "coordinator-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.RelateProducersInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().MultiGet(ctx, []string{"producer-id"}).Return(producers, nil)
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "coordinator-id", "producer-id").Return(nil)
			},
			input: &user.RelateProducersInput{
				CoordinatorID: "coordinator-id",
				ProducerIDs:   []string{"producer-id"},
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.RelateProducersInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, user.ErrNotFound)
			},
			input: &user.RelateProducersInput{
				CoordinatorID: "coordinator-id",
				ProducerIDs:   []string{"producer-id"},
			},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &user.RelateProducersInput{
				CoordinatorID: "coordinator-id",
				ProducerIDs:   []string{"producer-id"},
			},
			expectErr: user.ErrInternal,
		},
		{
			name: "failed to multi get producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().MultiGet(ctx, []string{"producer-id"}).Return(nil, assert.AnError)
			},
			input: &user.RelateProducersInput{
				CoordinatorID: "coordinator-id",
				ProducerIDs:   []string{"producer-id"},
			},
			expectErr: user.ErrInternal,
		},
		{
			name: "failed to contain invalid producers",
			setup: func(ctx context.Context, mocks *mocks) {
				producers := entity.Producers{{
					AdminID:       "producer-id",
					CoordinatorID: "coordinator-id",
				}}
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().MultiGet(ctx, []string{"producer-id"}).Return(producers, nil)
			},
			input: &user.RelateProducersInput{
				CoordinatorID: "coordinator-id",
				ProducerIDs:   []string{"producer-id"},
			},
			expectErr: user.ErrFailedPrecondition,
		},
		{
			name: "failed to update relationship",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().MultiGet(ctx, []string{"producer-id"}).Return(producers, nil)
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "coordinator-id", "producer-id").Return(assert.AnError)
			},
			input: &user.RelateProducersInput{
				CoordinatorID: "coordinator-id",
				ProducerIDs:   []string{"producer-id"},
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RelateProducers(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUnrelateProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UnrelateProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "", "producer-id").Return(nil)
			},
			input: &user.UnrelateProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UnrelateProducerInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to update relationship",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "", "producer-id").Return(assert.AnError)
			},
			input: &user.UnrelateProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: user.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UnrelateProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteProducer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.DeleteProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Delete(ctx, "producer-id", gomock.Any()).Return(nil)
			},
			input: &user.DeleteProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.DeleteProducerInput{},
			expectErr: user.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Delete(ctx, "producer-id", gomock.Any()).Return(assert.AnError)
			},
			input: &user.DeleteProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: user.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
