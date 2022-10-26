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
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:       "admin-id",
			CoordinatorID: "coordinator-id",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
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
				mocks.db.Producer.EXPECT().List(gomock.Any(), params).Return(nil, errmock)
				mocks.db.Producer.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListProducersInput{
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().List(gomock.Any(), params).Return(producers, nil)
				mocks.db.Producer.EXPECT().Count(gomock.Any(), params).Return(int64(0), errmock)
			},
			input: &user.ListProducersInput{
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
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			AdminID:       "admin-id",
			CoordinatorID: "coordinator-id",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(nil, errmock)
			},
			input: &user.MultiGetProducersInput{
				ProducerIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:       "admin-id",
		CoordinatorID: "coordinator-id",
		StoreName:     "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetProducerInput{
				ProducerID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
					StoreName:    "&.農園",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
					HeaderURL:    "https://and-period.jp/header.png",
					PhoneNumber:  "+819012345678",
					PostalCode:   "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
				}
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
			},
			input: &user.CreateProducerInput{
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
			},
			expectErr: nil,
		},
		{
			name: "success without notify register admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateProducerInput{
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
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateProducerInput{
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
			},
			expectErr: exception.ErrUnknown,
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

	params := &database.UpdateProducerParams{
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		StoreName:     "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
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
				mocks.db.Producer.EXPECT().Update(ctx, "producer-id", params).Return(nil)
			},
			input: &user.UpdateProducerInput{
				ProducerID:    "producer-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Update(ctx, "producer-id", params).Return(errmock)
			},
			input: &user.UpdateProducerInput{
				ProducerID:    "producer-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
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
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:       "admin-id",
		CoordinatorID: "coordinator-id",
		StoreName:     "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(nil, errmock)
			},
			input: &user.UpdateProducerEmailInput{
				ProducerID: "producer-id",
				Email:      "test-admin@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to admin change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(errmock)
			},
			input: &user.UpdateProducerEmailInput{
				ProducerID: "producer-id",
				Email:      "test-admin@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangeEmail(ctx, params).Return(nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "producer-id", "test-admin@and-period.jp").Return(errmock)
			},
			input: &user.UpdateProducerEmailInput{
				ProducerID: "producer-id",
				Email:      "test-admin@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
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

func TestResetProducerPassword(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	producer := &entity.Producer{
		Admin: entity.Admin{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Role:          entity.AdminRoleProducer,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
		},
		AdminID:       "admin-id",
		CoordinatorID: "coordinator-id",
		StoreName:     "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
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
						assert.Equal(t, params, expect)
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
				mocks.messenger.EXPECT().NotifyResetAdminPassword(gomock.Any(), gomock.Any()).Return(errmock)
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by admin id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(nil, errmock)
			},
			input: &user.ResetProducerPasswordInput{
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to admin change password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id").Return(producer, nil)
				mocks.adminAuth.EXPECT().AdminChangePassword(ctx, gomock.Any()).Return(errmock)
			},
			input: &user.ResetProducerPasswordInput{
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrUnknown,
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

func TestRelatedProducer(t *testing.T) {
	t.Parallel()

	producer := &entity.Producer{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &entity.Coordinator{
		AdminID: "coordinator-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.RelatedProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id", "coordinator_id").Return(&entity.Producer{}, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "producer-id", "coordinator-id").Return(nil)
			},
			input: &user.RelatedProducerInput{
				ProducerID:    "producer-id",
				CoordinatorID: "coordinator-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.RelatedProducerInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to already related",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id", "coordinator_id").Return(producer, nil)
			},
			input: &user.RelatedProducerInput{
				ProducerID:    "producer-id",
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id", "coordinator_id").Return(nil, errmock)
			},
			input: &user.RelatedProducerInput{
				ProducerID:    "producer-id",
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to not found coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id", "coordinator_id").Return(&entity.Producer{}, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, exception.ErrNotFound)
			},
			input: &user.RelatedProducerInput{
				ProducerID:    "producer-id",
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id", "coordinator_id").Return(&entity.Producer{}, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(nil, errmock)
			},
			input: &user.RelatedProducerInput{
				ProducerID:    "producer-id",
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update relationship",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Get(ctx, "producer-id", "coordinator_id").Return(&entity.Producer{}, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "coordinator-id").Return(coordinator, nil)
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "producer-id", "coordinator-id").Return(errmock)
			},
			input: &user.RelatedProducerInput{
				ProducerID:    "producer-id",
				CoordinatorID: "coordinator-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RelatedProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUnrelatedProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UnrelatedProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "producer-id", "").Return(nil)
			},
			input: &user.UnrelatedProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UnrelatedProducerInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update relationship",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().UpdateRelationship(ctx, "producer-id", "").Return(errmock)
			},
			input: &user.UnrelatedProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UnrelatedProducer(ctx, tt.input)
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Producer.EXPECT().Delete(ctx, "producer-id", gomock.Any()).Return(errmock)
			},
			input: &user.DeleteProducerInput{
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrUnknown,
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
