package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
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
			ID:            "admin-id",
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
				mocks.db.Administrator.EXPECT().List(gomock.Any(), params).Return(nil, errmock)
				mocks.db.Administrator.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListAdministratorsInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Administrator.EXPECT().List(gomock.Any(), params).Return(administrators, nil)
				mocks.db.Administrator.EXPECT().Count(gomock.Any(), params).Return(int64(0), errmock)
			},
			input: &user.ListAdministratorsInput{
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
			ID:            "admin-id",
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
				mocks.db.Administrator.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(nil, errmock)
			},
			input: &user.MultiGetAdministratorsInput{
				AdministratorIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		Email:         "test-admin@and-period.jp",
		PhoneNumber:   "+819012345678",
		CreatedAt:     now,
		UpdatedAt:     now,
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
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdministratorInput{
				AdministratorID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
				expectAuth := &entity.AdminAuth{
					Role: entity.AdminRoleAdministrator,
				}
				expectAdmin := &entity.Administrator{
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
				}
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.db.Administrator.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, auth *entity.AdminAuth, admin *entity.Administrator) error {
						expectAuth.AdminID, expectAuth.CognitoID = auth.AdminID, auth.CognitoID
						assert.Equal(t, expectAuth, auth)
						expectAdmin.ID = admin.ID
						assert.Equal(t, expectAdmin, admin)
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
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.db.Administrator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(errmock)
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
			name: "failed to create admin aith",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(errmock)
			},
			input: &user.CreateAdministratorInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.db.Administrator.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateAdministratorInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expectErr: exception.ErrUnknown,
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
