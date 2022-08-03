package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
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
		Orders: []*database.ListCoordinatorsOrder{
			{Key: entity.CoordinatorOrderByLastname, OrderByASC: true},
		},
	}
	coordinators := entity.Coordinators{
		{
			ID:               "admin-id",
			Lastname:         "&.",
			Firstname:        "スタッフ",
			LastnameKana:     "あんどぴりおど",
			FirstnameKana:    "すたっふ",
			StoreName:        "&.農園",
			ThumbnailURL:     "https://and-period.jp/thumbnail.png",
			HeaderURL:        "https://and-period.jp/header.png",
			TwitterAccount:   "twitter-account",
			InstagramAccount: "instagram-account",
			FacebookAccount:  "facebook-account",
			Email:            "test-admin@and-period.jp",
			PhoneNumber:      "+819012345678",
			PostalCode:       "1000014",
			Prefecture:       "東京都",
			City:             "千代田区",
			AddressLine1:     "永田町1-7-1",
			AddressLine2:     "",
			CreatedAt:        now,
			UpdatedAt:        now,
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
				Orders: []*user.ListCoordinatorsOrder{
					{Key: entity.CoordinatorOrderByLastname, OrderByASC: true},
				},
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
				mocks.db.Coordinator.EXPECT().List(gomock.Any(), params).Return(nil, errmock)
				mocks.db.Coordinator.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListCoordinatorsInput{
				Limit:  30,
				Offset: 0,
				Orders: []*user.ListCoordinatorsOrder{
					{Key: entity.CoordinatorOrderByLastname, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().List(gomock.Any(), params).Return(coordinators, nil)
				mocks.db.Coordinator.EXPECT().Count(gomock.Any(), params).Return(int64(0), errmock)
			},
			input: &user.ListCoordinatorsInput{
				Limit:  30,
				Offset: 0,
				Orders: []*user.ListCoordinatorsOrder{
					{Key: entity.CoordinatorOrderByLastname, OrderByASC: true},
				},
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
			ID:               "admin-id",
			Lastname:         "&.",
			Firstname:        "スタッフ",
			LastnameKana:     "あんどぴりおど",
			FirstnameKana:    "すたっふ",
			StoreName:        "&.農園",
			ThumbnailURL:     "https://and-period.jp/thumbnail.png",
			HeaderURL:        "https://and-period.jp/header.png",
			TwitterAccount:   "twitter-account",
			InstagramAccount: "instagram-account",
			FacebookAccount:  "facebook-account",
			Email:            "test-admin@and-period.jp",
			PhoneNumber:      "+819012345678",
			PostalCode:       "1000014",
			Prefecture:       "東京都",
			City:             "千代田区",
			AddressLine1:     "永田町1-7-1",
			AddressLine2:     "",
			CreatedAt:        now,
			UpdatedAt:        now,
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
				mocks.db.Coordinator.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(nil, errmock)
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
		ID:               "admin-id",
		Lastname:         "&.",
		Firstname:        "スタッフ",
		LastnameKana:     "あんどぴりおど",
		FirstnameKana:    "すたっふ",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-account",
		InstagramAccount: "instagram-account",
		FacebookAccount:  "facebook-account",
		Email:            "test-admin@and-period.jp",
		PhoneNumber:      "+819012345678",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		AddressLine1:     "永田町1-7-1",
		AddressLine2:     "",
		CreatedAt:        now,
		UpdatedAt:        now,
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
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
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

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateCoordinatorInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectAuth := &entity.AdminAuth{
					Role: entity.AdminRoleCoordinator,
				}
				expectCoordinator := &entity.Coordinator{
					Lastname:         "&.",
					Firstname:        "スタッフ",
					LastnameKana:     "あんどぴりおど",
					FirstnameKana:    "すたっふ",
					CompanyName:      "&.",
					StoreName:        "&.農園",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instgram-id",
					FacebookAccount:  "facebook-id",
					Email:            "test-admin@and-period.jp",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					AddressLine1:     "永田町1-7-1",
					AddressLine2:     "",
				}
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.db.Coordinator.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, auth *entity.AdminAuth, coordinator *entity.Coordinator) error {
						expectAuth.AdminID, expectAuth.CognitoID = auth.AdminID, auth.CognitoID
						assert.Equal(t, expectAuth, auth)
						expectCoordinator.ID = coordinator.ID
						assert.Equal(t, expectCoordinator, coordinator)
						return nil
					})
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:         "&.",
				Firstname:        "スタッフ",
				LastnameKana:     "あんどぴりおど",
				FirstnameKana:    "すたっふ",
				CompanyName:      "&.",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instgram-id",
				FacebookAccount:  "facebook-id",
				Email:            "test-admin@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
			},
			expectErr: nil,
		},
		{
			name: "success without notify register admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.db.Coordinator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:         "&.",
				Firstname:        "スタッフ",
				LastnameKana:     "あんどぴりおど",
				FirstnameKana:    "すたっふ",
				CompanyName:      "&.",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instgram-id",
				FacebookAccount:  "facebook-id",
				Email:            "test-admin@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
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
			name: "failed to create admin auth",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(errmock)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:         "&.",
				Firstname:        "スタッフ",
				LastnameKana:     "あんどぴりおど",
				FirstnameKana:    "すたっふ",
				CompanyName:      "&.",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instgram-id",
				FacebookAccount:  "facebook-id",
				Email:            "test-admin@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.db.Coordinator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateCoordinatorInput{
				Lastname:         "&.",
				Firstname:        "スタッフ",
				LastnameKana:     "あんどぴりおど",
				FirstnameKana:    "すたっふ",
				CompanyName:      "&.",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instgram-id",
				FacebookAccount:  "facebook-id",
				Email:            "test-admin@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
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

	params := &database.UpdateCoordinatorParams{
		Lastname:         "&.",
		Firstname:        "スタッフ",
		LastnameKana:     "あんどぴりおど",
		FirstnameKana:    "すたっふ",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		PhoneNumber:      "+819012345678",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		AddressLine1:     "永田町1-7-1",
		AddressLine2:     "",
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
				mocks.db.Coordinator.EXPECT().Update(ctx, "coordinator-id", params).Return(nil)
			},
			input: &user.UpdateCoordinatorInput{
				CoordinatorID:    "coordinator-id",
				Lastname:         "&.",
				Firstname:        "スタッフ",
				LastnameKana:     "あんどぴりおど",
				FirstnameKana:    "すたっふ",
				CompanyName:      "&.株式会社",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instagram-id",
				FacebookAccount:  "facebook-id",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
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
			name: "failed to update coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Coordinator.EXPECT().Update(ctx, "coordinator-id", params).Return(errmock)
			},
			input: &user.UpdateCoordinatorInput{
				CoordinatorID:    "coordinator-id",
				Lastname:         "&.",
				Firstname:        "スタッフ",
				LastnameKana:     "あんどぴりおど",
				FirstnameKana:    "すたっふ",
				CompanyName:      "&.株式会社",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instagram-id",
				FacebookAccount:  "facebook-id",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
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

	auth := &entity.AdminAuth{
		CognitoID: "cognito-id",
		Role:      entity.AdminRoleCoordinator,
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
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Coordinator.EXPECT().UpdateEmail(ctx, "coordinator-id", "test-admin@and-period.jp").Return(nil)
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
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(nil, errmock)
			},
			input: &user.UpdateCoordinatorEmailInput{
				CoordinatorID: "coordinator-id",
				Email:         "test-admin@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "invalid coordinator role",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{CognitoID: "cognito-id", Role: entity.AdminRoleAdministrator}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
			},
			input: &user.UpdateCoordinatorEmailInput{
				CoordinatorID: "coordinator-id",
				Email:         "test-admin@and-period.jp",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to admin change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(errmock)
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
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Coordinator.EXPECT().UpdateEmail(ctx, "coordinator-id", "test-admin@and-period.jp").Return(errmock)
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

func TestResetCoordinatorPassword(t *testing.T) {
	t.Parallel()

	auth := &entity.AdminAuth{
		CognitoID: "cognito-id",
		Role:      entity.AdminRoleCoordinator,
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
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
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
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyResetAdminPassword(gomock.Any(), gomock.Any()).Return(errmock)
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
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(nil, errmock)
			},
			input: &user.ResetCoordinatorPasswordInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "invalid coordinator role",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{CognitoID: "cognito-id", Role: entity.AdminRoleAdministrator}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
			},
			input: &user.ResetCoordinatorPasswordInput{
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to admin change password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "coordinator-id", "cognito_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(errmock)
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
