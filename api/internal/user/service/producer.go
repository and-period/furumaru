package service

import (
	"context"
	"errors"
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
		CoordinatorID: in.CoordinatorID,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
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
	adminParams := &entity.NewAdminParams{
		CognitoID:     cognitoID,
		Role:          entity.AdminRoleProducer,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		Email:         in.Email,
	}
	params := &entity.NewProducerParams{
		Admin:        entity.NewAdmin(adminParams),
		StoreName:    in.StoreName,
		ThumbnailURL: in.ThumbnailURL,
		HeaderURL:    in.HeaderURL,
		PhoneNumber:  in.PhoneNumber,
		PostalCode:   in.PostalCode,
		Prefecture:   in.Prefecture,
		City:         in.City,
		AddressLine1: in.AddressLine1,
		AddressLine2: in.AddressLine2,
	}
	producer := entity.NewProducer(params)
	auth := s.createCognitoAdmin(cognitoID, in.Email, password)
	if err := s.db.Producer.Create(ctx, producer, auth); err != nil {
		return nil, exception.InternalError(err)
	}
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
	producer, err := s.db.Producer.Get(ctx, in.ProducerID)
	if err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.AdminChangeEmailParams{
		Username: producer.CognitoID,
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
	producer, err := s.db.Producer.Get(ctx, in.ProducerID)
	if err != nil {
		return exception.InternalError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminChangePasswordParams{
		Username:  producer.CognitoID,
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

func (s *service) RelatedProducer(ctx context.Context, in *user.RelatedProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	producer, err := s.db.Producer.Get(ctx, in.ProducerID, "coordinator_id")
	if err != nil {
		return exception.InternalError(err)
	}
	if producer.CoordinatorID != "" {
		return fmt.Errorf("api: this producer is related: %w", exception.ErrFailedPrecondition)
	}
	_, err = s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid coordinator id: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return exception.InternalError(err)
	}
	err = s.db.Producer.UpdateRelationship(ctx, in.ProducerID, in.CoordinatorID)
	return exception.InternalError(err)
}

func (s *service) UnrelatedProducer(ctx context.Context, in *user.UnrelatedProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Producer.UpdateRelationship(ctx, in.ProducerID, "")
	return exception.InternalError(err)
}

func (s *service) DeleteProducer(ctx context.Context, in *user.DeleteProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Producer.Delete(ctx, in.ProducerID, s.deleteCognitoAdmin(in.ProducerID))
	return exception.InternalError(err)
}
