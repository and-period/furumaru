package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/random"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListProducers(ctx context.Context, in *user.ListProducersInput) (entity.Producers, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListProducersParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		producers entity.Producers
		total     int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producers, err = s.db.Producer.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Producer.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return producers, total, nil
}

func (s *service) MultiGetProducers(ctx context.Context, in *user.MultiGetProducersInput) (entity.Producers, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	producers, err := s.db.Producer.MultiGet(ctx, in.ProducerIDs)
	return producers, exception.InternalError(err)
}

func (s *service) GetProducer(ctx context.Context, in *user.GetProducerInput) (*entity.Producer, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	producer, err := s.db.Producer.Get(ctx, in.ProducerID)
	return producer, exception.InternalError(err)
}

func (s *service) CreateProducer(ctx context.Context, in *user.CreateProducerInput) (*entity.Producer, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	cognitoID := uuid.Base58Encode(uuid.New())
	password := random.NewStrings(size)
	if err := s.createCognitoAdmin(ctx, cognitoID, in.Email, password); err != nil {
		return nil, exception.InternalError(err)
	}
	adminParams := &entity.NewAdminParams{
		CognitoID:     cognitoID,
		Role:          entity.AdminRoleProducer,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		Email:         in.Email,
	}
	admin := entity.NewAdmin(adminParams)
	params := &entity.NewProducerParams{
		Admin:         admin,
		CoordinatorID: in.CoordinatorID,
		StoreName:     in.StoreName,
		ThumbnailURL:  in.ThumbnailURL,
		HeaderURL:     in.HeaderURL,
		PhoneNumber:   in.PhoneNumber,
		PostalCode:    in.PostalCode,
		Prefecture:    in.Prefecture,
		City:          in.City,
		AddressLine1:  in.AddressLine1,
		AddressLine2:  in.AddressLine2,
	}
	producer := entity.NewProducer(params)
	if err := s.db.Producer.Create(ctx, admin, producer); err != nil {
		return nil, exception.InternalError(err)
	}
	producer.Admin = *admin
	s.logger.Debug("Create producer", zap.String("producerId", producer.ID), zap.String("password", password))
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyRegisterAdmin(context.Background(), producer.ID, password)
		if err != nil {
			s.logger.Warn("Failed to notify register admin", zap.String("producerId", producer.ID), zap.Error(err))
		}
	}()
	return producer, nil
}

func (s *service) UpdateProducer(ctx context.Context, in *user.UpdateProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateProducerParams{
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		StoreName:     in.StoreName,
		ThumbnailURL:  in.ThumbnailURL,
		HeaderURL:     in.HeaderURL,
		PhoneNumber:   in.PhoneNumber,
		PostalCode:    in.PostalCode,
		Prefecture:    in.Prefecture,
		City:          in.City,
		AddressLine1:  in.AddressLine1,
		AddressLine2:  in.AddressLine2,
	}
	err := s.db.Producer.Update(ctx, in.ProducerID, params)
	return exception.InternalError(err)
}

func (s *service) UpdateProducerEmail(ctx context.Context, in *user.UpdateProducerEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	auth, err := s.db.Admin.Get(ctx, in.ProducerID, "cognito_id", "role")
	if err != nil {
		return exception.InternalError(err)
	}
	if auth.Role != entity.AdminRoleProducer {
		return fmt.Errorf("api: this admin role is not producer: %w", exception.ErrFailedPrecondition)
	}
	params := &cognito.AdminChangeEmailParams{
		Username: auth.CognitoID,
		Email:    in.Email,
	}
	if err := s.adminAuth.AdminChangeEmail(ctx, params); err != nil {
		return exception.InternalError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, in.ProducerID, in.Email)
	return exception.InternalError(err)
}

func (s *service) ResetProducerPassword(ctx context.Context, in *user.ResetProducerPasswordInput) error {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	auth, err := s.db.Admin.Get(ctx, in.ProducerID, "cognito_id", "role")
	if err != nil {
		return exception.InternalError(err)
	}
	if auth.Role != entity.AdminRoleProducer {
		return fmt.Errorf("api: this admin role is not producer: %w", exception.ErrFailedPrecondition)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminChangePasswordParams{
		Username:  auth.CognitoID,
		Password:  password,
		Permanent: true,
	}
	if err := s.adminAuth.AdminChangePassword(ctx, params); err != nil {
		return exception.InternalError(err)
	}
	s.logger.Debug("Reset producer password",
		zap.String("producerId", in.ProducerID), zap.String("password", password),
	)
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyResetAdminPassword(context.Background(), in.ProducerID, password)
		if err != nil {
			s.logger.Warn("Failed to notify reset admin password",
				zap.String("producerId", in.ProducerID), zap.Error(err),
			)
		}
	}()
	return nil
}
