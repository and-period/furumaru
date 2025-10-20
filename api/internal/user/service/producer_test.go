package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
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
				Type:          entity.AdminTypeProducer,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:        "admin-id",
			CoordinatorID:  "coordinator-id",
			Username:       "&.農園",
			ThumbnailURL:   "https://and-period.jp/thumbnail.png",
			HeaderURL:      "https://and-period.jp/header.png",
			InstagramID:    "instagram-account",
			FacebookID:     "facebook-account",
			PhoneNumber:    "+819012345678",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			CreatedAt:      now,
			UpdatedAt:      now,
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
			expectErr:   exception.ErrInvalidArgument,
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
			expectErr:   exception.ErrInternal,
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
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
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
				Type:          entity.AdminTypeProducer,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:        "admin-id",
			CoordinatorID:  "coordinator-id",
			Username:       "&.農園",
			ThumbnailURL:   "https://and-period.jp/thumbnail.png",
			HeaderURL:      "https://and-period.jp/header.png",
			InstagramID:    "instagram-account",
			FacebookID:     "facebook-account",
			PhoneNumber:    "+819012345678",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			CreatedAt:      now,
			UpdatedAt:      now,
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
			name: "success with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().MultiGetWithDeleted(ctx, []string{"admin-id"}).Return(producers, nil)
			},
			input: &user.MultiGetProducersInput{
				ProducerIDs: []string{"admin-id"},
				WithDeleted: true,
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
			expectErr: exception.ErrInvalidArgument,
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
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get producers with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().MultiGetWithDeleted(ctx, []string{"admin-id"}).Return(nil, assert.AnError)
			},
			input: &user.MultiGetProducersInput{
				ProducerIDs: []string{"admin-id"},
				WithDeleted: true,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
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
			Type:          entity.AdminTypeProducer,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:        "admin-id",
		CoordinatorID:  "coordinator-id",
		Username:       "&.農園",
		ThumbnailURL:   "https://and-period.jp/thumbnail.png",
		HeaderURL:      "https://and-period.jp/header.png",
		InstagramID:    "instagram-account",
		FacebookID:     "facebook-account",
		PhoneNumber:    "+819012345678",
		PostalCode:     "1000014",
		Prefecture:     "東京都",
		PrefectureCode: 13,
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
		CreatedAt:      now,
		UpdatedAt:      now,
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
			name: "success with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().GetWithDeleted(ctx, "admin-id").Return(producer, nil)
			},
			input: &user.GetProducerInput{
				ProducerID:  "admin-id",
				WithDeleted: true,
			},
			expect:    producer,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetProducerInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
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
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get producer with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().GetWithDeleted(ctx, "admin-id").Return(nil, assert.AnError)
			},
			input: &user.GetProducerInput{
				ProducerID:  "admin-id",
				WithDeleted: true,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
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
			Type:          entity.AdminTypeProducer,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:        "admin-id",
		PhoneNumber:    "+819012345678",
		Username:       "&.農園",
		ThumbnailURL:   "https://and-period.jp/thumbnail.png",
		HeaderURL:      "https://and-period.jp/header.png",
		InstagramID:    "instagram-account",
		FacebookID:     "facebook-account",
		PostalCode:     "1000014",
		Prefecture:     "東京都",
		PrefectureCode: 13,
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	shop := &entity.Shop{
		ID:             "shop-id",
		CoordinatorID:  "coordinator-id",
		ProducerIDs:    []string{},
		ProductTypeIDs: []string{"product-type-id"},
		Name:           "&.株式会社マルシェ",
		Activated:      true,
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
						Type:          entity.AdminTypeProducer,
						GroupIDs:      []string{"group-id"},
						Lastname:      "&.",
						Firstname:     "スタッフ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "すたっふ",
						Email:         "test-admin@and-period.jp",
					},
					CoordinatorID:  "coordinator-id",
					Username:       "&.農園",
					Profile:        "紹介文です。",
					ThumbnailURL:   "https://and-period.jp/thumbnail.png",
					HeaderURL:      "https://and-period.jp/header.png",
					InstagramID:    "instgram-id",
					FacebookID:     "facebook-id",
					PhoneNumber:    "+819012345678",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
				}
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Shop.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(shop, nil)
				mocks.db.Producer.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, producer *entity.Producer, shopID string, auth func(ctx context.Context) error) error {
						expectProducer.ID = producer.ID
						expectProducer.AdminID = producer.ID
						assert.Equal(t, expectProducer, producer)
						return nil
					})
				mocks.store.EXPECT().RelateShopProducer(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateProducerInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, exception.ErrNotFound)
			},
			input: &user.CreateProducerInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Shop.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to new admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Shop.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(shop, nil)
			},
			input: &user.CreateProducerInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 100,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Shop.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(shop, nil)
				mocks.db.Producer.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateProducerInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProducer(t *testing.T) {
	t.Parallel()

	params := &database.UpdateProducerParams{
		Lastname:       "&.",
		Firstname:      "スタッフ",
		LastnameKana:   "あんどぴりおど",
		FirstnameKana:  "すたっふ",
		Username:       "&.農園",
		ThumbnailURL:   "https://and-period.jp/thumbnail.png",
		HeaderURL:      "https://and-period.jp/header.png",
		InstagramID:    "instagram-account",
		FacebookID:     "facebook-account",
		Email:          "test-admin@and-period.jp",
		PhoneNumber:    "+819012345678",
		PostalCode:     "1000014",
		PrefectureCode: 13,
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
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
				mocks.db.Producer.EXPECT().Update(ctx, "producer-id", &params).Return(nil)
			},
			input: &user.UpdateProducerInput{
				ProducerID:     "producer-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:      "https://tmp.and-period.jp/header.png",
				InstagramID:    "instagram-account",
				FacebookID:     "facebook-account",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateProducerInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid prefecture code",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.UpdateProducerInput{
				ProducerID:     "producer-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instagram-account",
				FacebookID:     "facebook-account",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 100,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Update(ctx, "producer-id", params).Return(assert.AnError)
			},
			input: &user.UpdateProducerInput{
				ProducerID:     "producer-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instagram-account",
				FacebookID:     "facebook-account",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteProducer(t *testing.T) {
	t.Parallel()

	shopsIn := &store.ListShopsInput{
		ProducerIDs: []string{"producer-id"},
		NoLimit:     true,
	}
	shops := sentity.Shops{
		{
			ID:             "shop-id",
			CoordinatorID:  "coordinator-id",
			ProducerIDs:    []string{"producer-id"},
			ProductTypeIDs: []string{"product-type-id"},
			Name:           "&.株式会社マルシェ",
			Activated:      true,
		},
	}
	deleteIn := &store.UnrelateShopProducerInput{
		ShopID:     "shop-id",
		ProducerID: "producer-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.DeleteProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().ListShops(ctx, shopsIn).Return(shops, int64(1), nil)
				mocks.store.EXPECT().UnrelateShopProducer(ctx, deleteIn).Return(nil)
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list shops",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().ListShops(ctx, shopsIn).Return(nil, int64(0), assert.AnError)
			},
			input: &user.DeleteProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to unrelate shop producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().ListShops(ctx, shopsIn).Return(shops, int64(1), nil)
				mocks.store.EXPECT().UnrelateShopProducer(ctx, deleteIn).Return(assert.AnError)
			},
			input: &user.DeleteProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().ListShops(ctx, shopsIn).Return(shops, int64(1), nil)
				mocks.store.EXPECT().UnrelateShopProducer(ctx, deleteIn).Return(nil)
				mocks.db.Producer.EXPECT().Delete(ctx, "producer-id", gomock.Any()).Return(assert.AnError)
			},
			input: &user.DeleteProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
