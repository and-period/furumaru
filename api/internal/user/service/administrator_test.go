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

func TestListAdministrators(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	params := &database.ListAdministratorsParams{
		Limit:  30,
		Offset: 0,
	}
	administrators := entity.Administrators{
		{
			Admin: entity.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleAdministrator,
				Type:          entity.AdminTypeAdministrator,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:     "admin-id",
			PhoneNumber: "+819012345678",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *user.ListAdministratorsInput
		expect      entity.Administrators
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().List(gomock.Any(), params).Return(administrators, nil)
				mocks.db.Administrator.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListAdministratorsInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      administrators,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &user.ListAdministratorsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Administrator.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListAdministratorsInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().List(gomock.Any(), params).Return(administrators, nil)
				mocks.db.Administrator.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &user.ListAdministratorsInput{
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
			actual, total, err := service.ListAdministrators(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetAdministrators(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	administrators := entity.Administrators{
		{
			Admin: entity.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleAdministrator,
				Type:          entity.AdminTypeAdministrator,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:     "admin-id",
			PhoneNumber: "+819012345678",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetAdministratorsInput
		expect    entity.Administrators
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(administrators, nil)
			},
			input: &user.MultiGetAdministratorsInput{
				AdministratorIDs: []string{"admin-id"},
			},
			expect:    administrators,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetAdministratorsInput{
				AdministratorIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(nil, assert.AnError)
			},
			input: &user.MultiGetAdministratorsInput{
				AdministratorIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetAdministrators(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetAdministrator(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	administrator := &entity.Administrator{
		Admin: entity.Admin{
			ID:            "admin-id",
			Role:          entity.AdminRoleAdministrator,
			Type:          entity.AdminTypeAdministrator,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:     "admin-id",
		PhoneNumber: "+819012345678",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetAdministratorInput
		expect    *entity.Administrator
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(administrator, nil)
			},
			input: &user.GetAdministratorInput{
				AdministratorID: "admin-id",
			},
			expect:    administrator,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetAdministratorInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(nil, assert.AnError)
			},
			input: &user.GetAdministratorInput{
				AdministratorID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetAdministrator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateAdministrator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateAdministratorInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectAdministrator := &entity.Administrator{
					Admin: entity.Admin{
						Role:          entity.AdminRoleAdministrator,
						Type:          entity.AdminTypeAdministrator,
						Lastname:      "&.",
						Firstname:     "スタッフ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "すたっふ",
						Email:         "test-admin@and-period.jp",
					},
					PhoneNumber: "+819012345678",
				}
				mocks.db.Administrator.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, administrator *entity.Administrator, auth func(ctx context.Context) error) error {
						expectAdministrator.ID = administrator.ID
						expectAdministrator.AdminID = administrator.AdminID
						expectAdministrator.CognitoID = administrator.CognitoID
						assert.Equal(t, expectAdministrator, administrator)
						return nil
					})
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &user.CreateAdministratorInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expectErr: nil,
		},
		{
			name: "success without notify register admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateAdministratorInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateAdministratorInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateAdministratorInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateAdministrator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateAdministrator(t *testing.T) {
	t.Parallel()

	params := &database.UpdateAdministratorParams{
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		PhoneNumber:   "+819012345678",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateAdministratorInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Update(ctx, "administrator-id", params).Return(nil)
			},
			input: &user.UpdateAdministratorInput{
				AdministratorID: "administrator-id",
				Lastname:        "&.",
				Firstname:       "スタッフ",
				LastnameKana:    "あんどぴりおど",
				FirstnameKana:   "すたっふ",
				PhoneNumber:     "+819012345678",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateAdministratorInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Update(ctx, "administrator-id", params).Return(assert.AnError)
			},
			input: &user.UpdateAdministratorInput{
				AdministratorID: "administrator-id",
				Lastname:        "&.",
				Firstname:       "スタッフ",
				LastnameKana:    "あんどぴりおど",
				FirstnameKana:   "すたっふ",
				PhoneNumber:     "+819012345678",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateAdministrator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateAdministratorEmail(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	administrator := &entity.Administrator{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Role:          entity.AdminRoleAdministrator,
			Type:          entity.AdminTypeAdministrator,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:     "admin-id",
		PhoneNumber: "+819012345678",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	params := &cognito.AdminChangeEmailParams{
		Username: "cognito-id",
		Email:    "test-admin@and-period.jp",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateAdministratorEmailInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(administrator, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "administrator-id", "test-admin@and-period.jp").Return(nil)
			},
			input: &user.UpdateAdministratorEmailInput{
				AdministratorID: "administrator-id",
				Email:           "test-admin@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateAdministratorEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(nil, assert.AnError)
			},
			input: &user.UpdateAdministratorEmailInput{
				AdministratorID: "administrator-id",
				Email:           "test-admin@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to admin change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(administrator, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(assert.AnError)
			},
			input: &user.UpdateAdministratorEmailInput{
				AdministratorID: "administrator-id",
				Email:           "test-admin@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to update email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(administrator, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "administrator-id", "test-admin@and-period.jp").Return(assert.AnError)
			},
			input: &user.UpdateAdministratorEmailInput{
				AdministratorID: "administrator-id",
				Email:           "test-admin@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateAdministratorEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestResetAdministratorPassword(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	administrator := &entity.Administrator{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Role:          entity.AdminRoleAdministrator,
			Type:          entity.AdminTypeAdministrator,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:     "admin-id",
		PhoneNumber: "+819012345678",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ResetAdministratorPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(administrator, nil)
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
			input: &user.ResetAdministratorPasswordInput{
				AdministratorID: "administrator-id",
			},
			expectErr: nil,
		},
		{
			name: "success without notify",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(administrator, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyResetAdminPassword(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.ResetAdministratorPasswordInput{
				AdministratorID: "administrator-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ResetAdministratorPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(nil, assert.AnError)
			},
			input: &user.ResetAdministratorPasswordInput{
				AdministratorID: "administrator-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to admin change password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Get(ctx, "administrator-id").Return(administrator, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &user.ResetAdministratorPasswordInput{
				AdministratorID: "administrator-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResetAdministratorPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteAdministrator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.DeleteAdministratorInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Delete(ctx, "administrator-id", gomock.Any()).Return(nil)
			},
			input: &user.DeleteAdministratorInput{
				AdministratorID: "administrator-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.DeleteAdministratorInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().Delete(ctx, "administrator-id", gomock.Any()).Return(assert.AnError)
			},
			input: &user.DeleteAdministratorInput{
				AdministratorID: "administrator-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteAdministrator(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
