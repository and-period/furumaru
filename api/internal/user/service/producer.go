package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
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
	params := &entity.NewProducerParams{
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		StoreName:     in.StoreName,
		ThumbnailURL:  in.ThumbnailURL,
		HeaderURL:     in.HeaderURL,
		Email:         in.Email,
		PhoneNumber:   in.PhoneNumber,
		PostalCode:    in.PostalCode,
		Prefecture:    in.Prefecture,
		City:          in.City,
		AddressLine1:  in.AddressLine1,
		AddressLine2:  in.AddressLine2,
	}
	producer := entity.NewProducer(params)
	auth := entity.NewAdminAuth(producer.ID, cognitoID, entity.AdminRoleProducer)
	if err := s.db.Producer.Create(ctx, auth, producer); err != nil {
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
