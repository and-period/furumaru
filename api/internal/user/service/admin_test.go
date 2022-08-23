package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *entity.Admin
		expect *Admin
	}{
		{
			name: "success",
			admin: &entity.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleAdministrator,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Admin{
				Admin: response.Admin{
					ID:            "admin-id",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAdmin(tt.admin))
		})
	}
}

func TestAdmin_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *Admin
		expect *response.Admin
	}{
		{
			name: "success",
			admin: &Admin{
				Admin: response.Admin{
					ID:            "admin-id",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
			expect: &response.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleAdministrator,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     1640962800,
				UpdatedAt:     1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admin.Response())
		})
	}
}

func TestAdmins(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins entity.Admins
		expect Admins
	}{
		{
			name: "success",
			admins: entity.Admins{
				{
					ID:            "admin-id01",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					ID:            "admin-id02",
					Role:          entity.AdminRoleCoordinator,
					Lastname:      "&.",
					Firstname:     "仲介者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "ちゅうかいしゃ",
					Email:         "test-admin02@and-period.jp",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Admins{
				{
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "仲介者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAdmins(tt.admins))
		})
	}
}

func TestAdmins_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins Admins
		expect map[string]*Admin
	}{
		{
			name: "success",
			admins: Admins{
				{
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "仲介者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: map[string]*Admin{
				"admin-id01": {
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				"admin-id02": {
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "仲介者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admins.Map())
		})
	}
}

func TestAdmins_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins Admins
		expect []*response.Admin
	}{
		{
			name: "success",
			admins: Admins{
				{
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "仲介者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: []*response.Admin{
				{
					ID:            "admin-id01",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				{
					ID:            "admin-id02",
					Role:          entity.AdminRoleCoordinator,
					Lastname:      "&.",
					Firstname:     "仲介者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "ちゅうかいしゃ",
					Email:         "test-admin02@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admins.Response())
		})
	}
}

func TestMultiGetAdmins(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	adminIDs := []string{
		"administrator-id",
		"coordinator-id",
		"producer-id",
	}
	auths := entity.AdminAuths{
		{AdminID: "administrator-id", Role: entity.AdminRoleAdministrator},
		{AdminID: "coordinator-id", Role: entity.AdminRoleCoordinator},
		{AdminID: "producer-id", Role: entity.AdminRoleProducer},
	}
	administrators := entity.Administrators{
		{
			ID:            "administrator-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			PhoneNumber:   "+819012345678",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
	coordinators := entity.Coordinators{
		{
			ID:               "coordinator-id",
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
	producers := entity.Producers{
		{
			ID:            "producer-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			Email:         "test-admin@and-period.jp",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			AddressLine1:  "永田町1-7-1",
			AddressLine2:  "",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
	admins := entity.Admins{
		{
			ID:            "administrator-id",
			Role:          entity.AdminRoleAdministrator,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            "coordinator-id",
			Role:          entity.AdminRoleCoordinator,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            "producer-id",
			Role:          entity.AdminRoleProducer,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetAdminsInput
		expect    entity.Admins
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, adminIDs).Return(auths, nil)
				mocks.db.Administrator.EXPECT().MultiGet(gomock.Any(), []string{"administrator-id"}).Return(administrators, nil)
				mocks.db.Coordinator.EXPECT().MultiGet(gomock.Any(), []string{"coordinator-id"}).Return(coordinators, nil)
				mocks.db.Producer.EXPECT().MultiGet(gomock.Any(), []string{"producer-id"}).Return(producers, nil)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: adminIDs,
			},
			expect:    admins,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetAdminsInput{
				AdminIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get admin auths",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, adminIDs).Return(nil, errmock)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: adminIDs,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to multi get administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, adminIDs).Return(auths, nil)
				mocks.db.Administrator.EXPECT().MultiGet(gomock.Any(), []string{"administrator-id"}).Return(nil, errmock)
				mocks.db.Coordinator.EXPECT().MultiGet(gomock.Any(), []string{"coordinator-id"}).Return(coordinators, nil)
				mocks.db.Producer.EXPECT().MultiGet(gomock.Any(), []string{"producer-id"}).Return(producers, nil)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: adminIDs,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to multi get coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, adminIDs).Return(auths, nil)
				mocks.db.Administrator.EXPECT().MultiGet(gomock.Any(), []string{"administrator-id"}).Return(administrators, nil)
				mocks.db.Coordinator.EXPECT().MultiGet(gomock.Any(), []string{"coordinator-id"}).Return(nil, errmock)
				mocks.db.Producer.EXPECT().MultiGet(gomock.Any(), []string{"producer-id"}).Return(producers, nil)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: adminIDs,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to multi get producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, adminIDs).Return(auths, nil)
				mocks.db.Administrator.EXPECT().MultiGet(gomock.Any(), []string{"administrator-id"}).Return(administrators, nil)
				mocks.db.Coordinator.EXPECT().MultiGet(gomock.Any(), []string{"coordinator-id"}).Return(coordinators, nil)
				mocks.db.Producer.EXPECT().MultiGet(gomock.Any(), []string{"producer-id"}).Return(nil, errmock)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: adminIDs,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to invalid role",
			setup: func(ctx context.Context, mocks *mocks) {
				auths := entity.AdminAuths{{AdminID: "admin-id", Role: entity.AdminRoleUnknown}}
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, adminIDs).Return(auths, nil)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: adminIDs,
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetAdmins(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetAdminDevices(t *testing.T) {
	t.Parallel()

	auths := entity.AdminAuths{
		{
			AdminID:   "admin-id",
			CognitoID: "username",
			Role:      entity.AdminRoleAdministrator,
			Device:    "instance-id",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetAdminDevicesInput
		expect    []string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, []string{"admin-id"}, "device").Return(auths, nil)
			},
			input: &user.MultiGetAdminDevicesInput{
				AdminIDs: []string{"admin-id"},
			},
			expect:    []string{"instance-id"},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetAdminDevicesInput{
				AdminIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get admin auths",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().MultiGet(ctx, []string{"admin-id"}, "device").Return(nil, errmock)
			},
			input: &user.MultiGetAdminDevicesInput{
				AdminIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetAdminDevices(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetAdmin(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	admin := &entity.Admin{
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		Email:         "test-admin@and-period.jp",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetAdminInput
		expect    *entity.Admin
		expectErr error
	}{
		{
			name: "success to administartor",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleAdministrator}
				administrator := &entity.Administrator{
					ID:            admin.ID,
					Lastname:      admin.Lastname,
					Firstname:     admin.Firstname,
					LastnameKana:  admin.LastnameKana,
					FirstnameKana: admin.FirstnameKana,
					Email:         admin.Email,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(administrator, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect: func() *entity.Admin {
				admin := *admin
				admin.Role = entity.AdminRoleAdministrator
				return &admin
			}(),
			expectErr: nil,
		},
		{
			name: "success to coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleCoordinator}
				coordinator := &entity.Coordinator{
					ID:            admin.ID,
					Lastname:      admin.Lastname,
					Firstname:     admin.Firstname,
					LastnameKana:  admin.LastnameKana,
					FirstnameKana: admin.FirstnameKana,
					Email:         admin.Email,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(coordinator, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect: func() *entity.Admin {
				admin := *admin
				admin.Role = entity.AdminRoleCoordinator
				return &admin
			}(),
			expectErr: nil,
		},
		{
			name: "success to producer",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleProducer}
				producer := &entity.Producer{
					ID:            admin.ID,
					Lastname:      admin.Lastname,
					Firstname:     admin.Firstname,
					LastnameKana:  admin.LastnameKana,
					FirstnameKana: admin.FirstnameKana,
					Email:         admin.Email,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(producer, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect: func() *entity.Admin {
				admin := *admin
				admin.Role = entity.AdminRoleProducer
				return &admin
			}(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetAdminInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admin auth",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleAdministrator}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleCoordinator}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleProducer}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unknown role",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleUnknown}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
