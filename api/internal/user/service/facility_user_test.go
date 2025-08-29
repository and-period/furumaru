package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateFacilityUser(t *testing.T) {
	t.Parallel()

	lastCheckInAt := jst.Date(2025, 8, 27, 12, 0, 0, 0)

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateFacilityUserInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectUser := &entity.User{
					Type:       entity.UserTypeFacilityUser,
					Registered: false,
					FacilityUser: entity.FacilityUser{
						ExternalID:    "external-id",
						ProducerID:    "producer-id",
						Lastname:      "田中",
						Firstname:     "太郎",
						LastnameKana:  "たなか",
						FirstnameKana: "たろう",
						ProviderType:  entity.UserAuthProviderTypeLINE,
						Email:         "test@example.com",
						PhoneNumber:   "+819012345678",
						LastCheckInAt: lastCheckInAt,
					},
				}
				mocks.db.FacilityUser.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, user *entity.User) error {
						expectUser.ID = user.ID
						expectUser.FacilityUser.UserID = user.ID
						assert.Equal(t, expectUser, user)
						return nil
					})
			},
			input: &user.CreateFacilityUserInput{
				ProducerID:    "producer-id",
				ProviderType:  entity.UserAuthProviderTypeLINE,
				ProviderID:    "external-id",
				Lastname:      "田中",
				Firstname:     "太郎",
				LastnameKana:  "たなか",
				FirstnameKana: "たろう",
				Email:         "test@example.com",
				PhoneNumber:   "+819012345678",
				LastCheckInAt: lastCheckInAt,
			},
			expectErr: nil,
		},
		{
			name:  "validation error - missing required fields",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.CreateFacilityUserInput{
				ProducerID:   "",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "external-id",
				Email:        "test@example.com",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "validation error - invalid last checkin at (future date)",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.CreateFacilityUserInput{
				ProducerID:    "producer-id",
				ProviderType:  entity.UserAuthProviderTypeLINE,
				ProviderID:    "external-id",
				Lastname:      "田中",
				Firstname:     "太郎",
				LastnameKana:  "たなか",
				FirstnameKana: "たろう",
				Email:         "test@example.com",
				PhoneNumber:   "+819012345678",
				LastCheckInAt: jst.Now().Add(time.Hour), // future date
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "database error",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.FacilityUser.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateFacilityUserInput{
				ProducerID:    "producer-id",
				ProviderType:  entity.UserAuthProviderTypeLINE,
				ProviderID:    "external-id",
				Lastname:      "田中",
				Firstname:     "太郎",
				LastnameKana:  "たなか",
				FirstnameKana: "たろう",
				Email:         "test@example.com",
				PhoneNumber:   "+819012345678",
				LastCheckInAt: lastCheckInAt,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateFacilityUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestGetFacilityUser(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	fuser := &entity.FacilityUser{
		UserID:        "user-id",
		ExternalID:    "external-id",
		ProducerID:    "producer-id",
		ProviderType:  entity.UserAuthProviderTypeLINE,
		LastCheckInAt: now,
	}
	expectUser := &entity.User{
		ID:         "user-id",
		Type:       entity.UserTypeFacilityUser,
		Registered: false,
		CreatedAt:  now,
		UpdatedAt:  now,
		FacilityUser: entity.FacilityUser{
			UserID:        "user-id",
			ExternalID:    "external-id",
			ProducerID:    "producer-id",
			Lastname:      "田中",
			Firstname:     "太郎",
			LastnameKana:  "たなか",
			FirstnameKana: "たろう",
			ProviderType:  entity.UserAuthProviderTypeLINE,
			Email:         "test@example.com",
			PhoneNumber:   "+819012345678",
			LastCheckInAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *user.GetFacilityUserInput
		expect      *entity.User
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.FacilityUser.EXPECT().
					GetByExternalID(ctx, entity.UserAuthProviderTypeLINE, "external-id", "producer-id").
					Return(fuser, nil)
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(expectUser, nil)
			},
			input: &user.GetFacilityUserInput{
				ProducerID:   "producer-id",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "external-id",
			},
			expect:    expectUser,
			expectErr: nil,
		},
		{
			name:  "validation error - missing producer id",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.GetFacilityUserInput{
				ProducerID:   "",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "external-id",
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "validation error - missing provider id",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.GetFacilityUserInput{
				ProducerID:   "producer-id",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "",
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "facility user not found",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.FacilityUser.EXPECT().
					GetByExternalID(ctx, entity.UserAuthProviderTypeLINE, "external-id", "producer-id").
					Return(nil, exception.ErrNotFound)
			},
			input: &user.GetFacilityUserInput{
				ProducerID:   "producer-id",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "external-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "user not found",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.FacilityUser.EXPECT().
					GetByExternalID(ctx, entity.UserAuthProviderTypeLINE, "external-id", "producer-id").
					Return(fuser, nil)
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(nil, exception.ErrNotFound)
			},
			input: &user.GetFacilityUserInput{
				ProducerID:   "producer-id",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "external-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "database error on facility user query",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.FacilityUser.EXPECT().
					GetByExternalID(ctx, entity.UserAuthProviderTypeLINE, "external-id", "producer-id").
					Return(nil, assert.AnError)
			},
			input: &user.GetFacilityUserInput{
				ProducerID:   "producer-id",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "external-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "database error on user query",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.FacilityUser.EXPECT().
					GetByExternalID(ctx, entity.UserAuthProviderTypeLINE, "external-id", "producer-id").
					Return(fuser, nil)
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(nil, assert.AnError)
			},
			input: &user.GetFacilityUserInput{
				ProducerID:   "producer-id",
				ProviderType: entity.UserAuthProviderTypeLINE,
				ProviderID:   "external-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetFacilityUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
