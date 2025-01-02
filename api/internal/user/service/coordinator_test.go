package service

import (
	"context"
	"testing"
	"time"

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
				Type:          entity.AdminTypeCoordinator,
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
			BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			CreatedAt:      now,
			UpdatedAt:      now,
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
			expectErr:   exception.ErrInternal,
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
			expectErr:   exception.ErrInternal,
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
				Type:          entity.AdminTypeCoordinator,
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
			BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			CreatedAt:      now,
			UpdatedAt:      now,
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
			name: "success with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().MultiGetWithDeleted(ctx, []string{"admin-id"}).Return(coordinators, nil)
			},
			input: &user.MultiGetCoordinatorsInput{
				CoordinatorIDs: []string{"admin-id"},
				WithDeleted:    true,
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
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get coordinators with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().MultiGetWithDeleted(ctx, []string{"admin-id"}).Return(nil, assert.AnError)
			},
			input: &user.MultiGetCoordinatorsInput{
				CoordinatorIDs: []string{"admin-id"},
				WithDeleted:    true,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
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
			Type:          entity.AdminTypeCoordinator,
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
		BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
		CreatedAt:      now,
		UpdatedAt:      now,
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
			name: "success with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().GetWithDeleted(ctx, "admin-id").Return(coordinator, nil)
			},
			input: &user.GetCoordinatorInput{
				CoordinatorID: "admin-id",
				WithDeleted:   true,
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
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get coordinator with deleted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().GetWithDeleted(ctx, "admin-id").Return(nil, assert.AnError)
			},
			input: &user.GetCoordinatorInput{
				CoordinatorID: "admin-id",
				WithDeleted:   true,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
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
						Type:          entity.AdminTypeCoordinator,
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
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
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
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
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
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				ThumbnailURL:   "",
				HeaderURL:      "",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
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
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInternal,
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
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "failed to new admin",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.CreateCoordinatorInput{
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				ThumbnailURL:   "https://and-period.jp/thumbnail.png",
				HeaderURL:      "https://and-period.jp/header.png",
				InstagramID:    "instgram-id",
				FacebookID:     "facebook-id",
				Email:          "test-admin@and-period.jp",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				PrefectureCode: -1,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:       "&.",
				Firstname:      "スタッフ",
				LastnameKana:   "あんどぴりおど",
				FirstnameKana:  "すたっふ",
				MarcheName:     "&.マルシェ",
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
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, _, err := service.CreateCoordinator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateCoordinator(t *testing.T) {
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
		PrefectureCode: 13,
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
		BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
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
				mocks.store.EXPECT().MultiGetProductTypes(ctx, productTypesIn).Return(productTypes, nil)
				mocks.db.Coordinator.EXPECT().Update(ctx, "coordinator-id", &params).Return(nil)
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
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
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
			name:  "failed to invalid prefecture code",
			setup: func(ctx context.Context, mocks *mocks) {},
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
				PrefectureCode: -1,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get product types",
			setup: func(ctx context.Context, mocks *mocks) {
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
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to unmatch product types length",
			setup: func(ctx context.Context, mocks *mocks) {
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
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
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
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expectErr: exception.ErrInternal,
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
			Type:          entity.AdminTypeCoordinator,
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
			expectErr: exception.ErrInternal,
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
			expectErr: exception.ErrInternal,
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
			expectErr: exception.ErrInternal,
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

func TestResetCoordinatorPassword(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	coordinator := &entity.Coordinator{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Role:          entity.AdminRoleCoordinator,
			Type:          entity.AdminTypeCoordinator,
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
			expectErr: exception.ErrInternal,
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
			expectErr: exception.ErrInternal,
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

func TestRemoveCoordinatorProductType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.RemoveCoordinatorProductTypeInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().RemoveProductTypeID(ctx, "product-type-id").Return(nil)
			},
			input: &user.RemoveCoordinatorProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.RemoveCoordinatorProductTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to remove product type id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().RemoveProductTypeID(ctx, "product-type-id").Return(assert.AnError)
			},
			input: &user.RemoveCoordinatorProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RemoveCoordinatorProductType(ctx, tt.input)
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
			expectErr: exception.ErrInternal,
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
			expectErr: exception.ErrInternal,
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
