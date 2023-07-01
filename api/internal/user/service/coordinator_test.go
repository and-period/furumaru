package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListCoordinators(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	params := &database.ListCoordinatorsParams{
		Limit:  30,
		Offset: 0,
	}
	coordinators := entity.Coordinators{
		{
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
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *user.ListCoordinatorsInput
		expect      entity.Coordinators
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().List(gomock.Any(), params).Return(coordinators, nil)
				mocks.db.Coordinator.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListCoordinatorsInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      coordinators,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &user.ListCoordinatorsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Coordinator.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListCoordinatorsInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().List(gomock.Any(), params).Return(coordinators, nil)
				mocks.db.Coordinator.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &user.ListCoordinatorsInput{
				Limit:  30,
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
			actual, total, err := service.ListCoordinators(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetCoordinators(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	coordinators := entity.Coordinators{
		{
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
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetCoordinatorsInput
		expect    entity.Coordinators
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(coordinators, nil)
			},
			input: &user.MultiGetCoordinatorsInput{
				CoordinatorIDs: []string{"admin-id"},
			},
			expect:    coordinators,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetCoordinatorsInput{
				CoordinatorIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(nil, assert.AnError)
			},
			input: &user.MultiGetCoordinatorsInput{
				CoordinatorIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetCoordinators(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetCoordinator(t *testing.T) {
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
		input     *user.GetCoordinatorInput
		expect    *entity.Coordinator
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(coordinator, nil)
			},
			input: &user.GetCoordinatorInput{
				CoordinatorID: "admin-id",
			},
			expect:    coordinator,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetCoordinatorInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(nil, assert.AnError)
			},
			input: &user.GetCoordinatorInput{
				CoordinatorID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetCoordinator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateCoordinator(t *testing.T) {
	t.Parallel()

	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:   "product-type-id",
			Name: "じゃがいも",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateCoordinatorInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectCoordinator := &entity.Coordinator{
					Admin: entity.Admin{
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "スタッフ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "すたっふ",
						Email:         "test-admin@and-period.jp",
					},
					MarcheName:     "&.マルシェ",
					Username:       "&.農園",
					Profile:        "紹介文です。",
					ProductTypeIDs: []string{"product-type-id"},
					ThumbnailURL:   "https://and-period.jp/thumbnail.png",
					HeaderURL:      "https://and-period.jp/header.png",
					InstagramID:    "instgram-id",
					FacebookID:     "facebook-id",
					PhoneNumber:    "+819012345678",
					PostalCode:     "1000014",
					Prefecture:     codes.PrefectureValues["tokyo"],
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
				}
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(productTypes, nil)
				mocks.db.Coordinator.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, coordinator *entity.Coordinator, auth func(ctx context.Context) error) error {
						expectCoordinator.ID = coordinator.ID
						expectCoordinator.AdminID = coordinator.ID
						expectCoordinator.CognitoID = coordinator.CognitoID
						assert.Equal(t, expectCoordinator, coordinator)
						return nil
					})
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mocks.media.EXPECT().ResizeCoordinatorThumbnail(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mocks.media.EXPECT().ResizeCoordinatorHeader(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: nil,
		},
		{
			name: "success without notify register admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				MarcheName:    "&.マルシェ",
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
			input:     &user.CreateCoordinatorInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get product types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(nil, assert.AnError)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unmatch product types length",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(sentity.ProductTypes{}, nil)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				MarcheName:    "&.マルシェ",
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
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateCoordinator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateCoordinator(t *testing.T) {
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
	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:   "product-type-id",
			Name: "じゃがいも",
		},
	}
	params := &database.UpdateCoordinatorParams{
		Lastname:       "&.",
		Firstname:      "スタッフ",
		LastnameKana:   "あんどぴりおど",
		FirstnameKana:  "すたっふ",
		MarcheName:     "&.株式会社マルシェ",
		Username:       "&.農園",
		Profile:        "紹介文です。",
		ProductTypeIDs: []string{"product-type-id"},
		ThumbnailURL:   "https://and-period.jp/thumbnail.png",
		HeaderURL:      "https://and-period.jp/header.png",
		InstagramID:    "instagram-id",
		FacebookID:     "facebook-id",
		PhoneNumber:    "+819012345678",
		PostalCode:     "1000014",
		Prefecture:     codes.PrefectureValues["tokyo"],
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateCoordinatorInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := *params
				params.ThumbnailURL = "https://tmp.and-period.jp/thumbnail.png"
				params.HeaderURL = "https://tmp.and-period.jp/header.png"
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(productTypes, nil)
				mocks.db.Coordinator.EXPECT().Update(ctx, "coordinator-id", &params).Return(nil)
				mocks.media.EXPECT().ResizeCoordinatorThumbnail(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mocks.media.EXPECT().ResizeCoordinatorHeader(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.UpdateCoordinatorInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.株式会社マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:      "https://tmp.and-period.jp/header.png",
				InstagramID:    "instagram-id",
				FacebookID:     "facebook-id",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateCoordinatorInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &user.UpdateCoordinatorInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.株式会社マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instagram-id",
				FacebookID:     "facebook-id",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to multi get product types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(nil, assert.AnError)
			},
			input: &user.UpdateCoordinatorInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.株式会社マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instagram-id",
				FacebookID:     "facebook-id",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unmatch product types length",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(sentity.ProductTypes{}, nil)
			},
			input: &user.UpdateCoordinatorInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.株式会社マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instagram-id",
				FacebookID:     "facebook-id",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(productTypes, nil)
				mocks.db.Coordinator.EXPECT().Update(ctx, "coordinator-id", params).Return(assert.AnError)
			},
			input: &user.UpdateCoordinatorInput{
				CoordinatorID:  "coordinator-id",
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.株式会社マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instagram-id",
				FacebookID:     "facebook-id",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     codes.PrefectureValues["tokyo"],
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateCoordinator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateCoordinatorEmail(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	coordinator := &entity.Coordinator{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
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
	params := &cognito.AdminChangeEmailParams{
		Username: "cognito-id",
		Email:    "test-admin@and-period.jp",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateCoordinatorEmailInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "coordinator-id", "test-admin@and-period.jp").Return(nil)
			},
			input: &user.UpdateCoordinatorEmailInput{
				CoordinatorID: "coordinator-id",
				Email:         "test-admin@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateCoordinatorEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &user.UpdateCoordinatorEmailInput{
				CoordinatorID: "coordinator-id",
				Email:         "test-admin@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to admin change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(assert.AnError)
			},
			input: &user.UpdateCoordinatorEmailInput{
				CoordinatorID: "coordinator-id",
				Email:         "test-admin@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "coordinator-id", "test-admin@and-period.jp").Return(assert.AnError)
			},
			input: &user.UpdateCoordinatorEmailInput{
				CoordinatorID: "coordinator-id",
				Email:         "test-admin@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateCoordinatorEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateCoordinatorThumbnails(t *testing.T) {
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
		input     *user.UpdateCoordinatorThumbnailsInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().UpdateThumbnails(ctx, "coordinator-id", thumbnails).Return(nil)
			},
			input: &user.UpdateCoordinatorThumbnailsInput{
				CoordinatorID: "coordinator-id",
				Thumbnails:    thumbnails,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateCoordinatorThumbnailsInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update thumbnails",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().UpdateThumbnails(ctx, "coordinator-id", thumbnails).Return(assert.AnError)
			},
			input: &user.UpdateCoordinatorThumbnailsInput{
				CoordinatorID: "coordinator-id",
				Thumbnails:    thumbnails,
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateCoordinatorThumbnails(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateCoordinatorHeaders(t *testing.T) {
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
		input     *user.UpdateCoordinatorHeadersInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().UpdateHeaders(ctx, "coordinator-id", headers).Return(nil)
			},
			input: &user.UpdateCoordinatorHeadersInput{
				CoordinatorID: "coordinator-id",
				Headers:       headers,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateCoordinatorHeadersInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update headers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().UpdateHeaders(ctx, "coordinator-id", headers).Return(assert.AnError)
			},
			input: &user.UpdateCoordinatorHeadersInput{
				CoordinatorID: "coordinator-id",
				Headers:       headers,
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateCoordinatorHeaders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestResetCoordinatorPassword(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	coordinator := &entity.Coordinator{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Role:          entity.AdminRoleCoordinator,
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
		input     *user.ResetCoordinatorPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.adminAuth.EXPECT().
					AdminChangePassword(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *cognito.AdminChangePasswordParams) error {
						expect := &cognito.AdminChangePasswordParams{
							Username:  "cognito-id",
							Password:  params.Password, // ignore
							Permanent: true,
						}
						assert.Equal(t, params, expect)
						return nil
					})
				mocks.messenger.EXPECT().NotifyResetAdminPassword(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &user.ResetCoordinatorPasswordInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: nil,
		},
		{
			name: "success without notify",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyResetAdminPassword(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.ResetCoordinatorPasswordInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ResetCoordinatorPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &user.ResetCoordinatorPasswordInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to admin change password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &user.ResetCoordinatorPasswordInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResetCoordinatorPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteCoordinator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.DeleteCoordinatorInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Delete(ctx, "coordinator-id", gomock.Any()).Return(nil)
			},
			input: &user.DeleteCoordinatorInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.DeleteCoordinatorInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Delete(ctx, "coordinator-id", gomock.Any()).Return(assert.AnError)
			},
			input: &user.DeleteCoordinatorInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteCoordinator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestAggregateRelatedProducers(t *testing.T) {
	t.Parallel()
	total := map[string]int64{
		"coordinator-id": 6,
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.AggregateRealatedProducersInput
		expect    map[string]int64
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().AggregateByCoordinatorID(ctx, []string{"coordinator-id"}).Return(total, nil)
			},
			input: &user.AggregateRealatedProducersInput{
				CoordinatorIDs: []string{"coordinator-id"},
			},
			expect:    total,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.AggregateRealatedProducersInput{
				CoordinatorIDs: []string{""},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().AggregateByCoordinatorID(ctx, []string{"coordinator-id"}).Return(nil, assert.AnError)
			},
			input: &user.AggregateRealatedProducersInput{
				CoordinatorIDs: []string{"coordinator-id"},
			},
			expectErr: exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			res, err := service.AggregateRealatedProducers(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, res)
		}))
	}
}
