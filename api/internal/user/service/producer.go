package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
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
		return nil, 0, internalError(err)
	}
	params := &database.ListProducersParams{
		CoordinatorID: in.CoordinatorID,
		Username:      in.Username,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
		OnlyUnrelated: in.OnlyUnrelated,
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
		return nil, 0, internalError(err)
	}
	return producers, total, nil
}

func (s *service) MultiGetProducers(ctx context.Context, in *user.MultiGetProducersInput) (entity.Producers, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	producers, err := s.db.Producer.MultiGet(ctx, in.ProducerIDs)
	return producers, internalError(err)
}

func (s *service) GetProducer(ctx context.Context, in *user.GetProducerInput) (*entity.Producer, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	producer, err := s.db.Producer.Get(ctx, in.ProducerID)
	return producer, internalError(err)
}

func (s *service) CreateProducer(ctx context.Context, in *user.CreateProducerInput) (*entity.Producer, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	_, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid coordinator id: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
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
		Admin:             entity.NewAdmin(adminParams),
		CoordinatorID:     in.CoordinatorID,
		PhoneNumber:       in.PhoneNumber,
		Username:          in.Username,
		Profile:           in.Profile,
		ThumbnailURL:      in.ThumbnailURL,
		HeaderURL:         in.HeaderURL,
		PromotionVideoURL: in.PromotionVideoURL,
		BonusVideoURL:     in.BonusVideoURL,
		InstagramID:       in.InstagramID,
		FacebookID:        in.FacebookID,
		PostalCode:        in.PostalCode,
		Prefecture:        in.Prefecture,
		City:              in.City,
		AddressLine1:      in.AddressLine1,
		AddressLine2:      in.AddressLine2,
	}
	producer := entity.NewProducer(params)
	auth := s.createCognitoAdmin(cognitoID, in.Email, password)
	if err := s.db.Producer.Create(ctx, producer, auth); err != nil {
		return nil, internalError(err)
	}
	s.logger.Debug("Create producer", zap.String("producerId", producer.ID), zap.String("password", password))
	s.waitGroup.Add(2)
	go func() {
		defer s.waitGroup.Done()
		s.resizeProducer(context.Background(), producer.ID, in.ThumbnailURL, in.HeaderURL)
	}()
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
		return internalError(err)
	}
	producer, err := s.db.Producer.Get(ctx, in.ProducerID)
	if err != nil {
		return internalError(err)
	}
	params := &database.UpdateProducerParams{
		Lastname:          in.Lastname,
		Firstname:         in.Firstname,
		LastnameKana:      in.LastnameKana,
		FirstnameKana:     in.FirstnameKana,
		Username:          in.Username,
		Profile:           in.Profile,
		ThumbnailURL:      in.ThumbnailURL,
		HeaderURL:         in.HeaderURL,
		PromotionVideoURL: in.PromotionVideoURL,
		BonusVideoURL:     in.BonusVideoURL,
		InstagramID:       in.InstagramID,
		FacebookID:        in.FacebookID,
		PhoneNumber:       in.PhoneNumber,
		PostalCode:        in.PostalCode,
		Prefecture:        in.Prefecture,
		City:              in.City,
		AddressLine1:      in.AddressLine1,
		AddressLine2:      in.AddressLine2,
	}
	if err := s.db.Producer.Update(ctx, in.ProducerID, params); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		var thumbnailURL, headerURL string
		if producer.ThumbnailURL != in.ThumbnailURL {
			thumbnailURL = in.ThumbnailURL
		}
		if producer.HeaderURL != in.HeaderURL {
			headerURL = in.HeaderURL
		}
		s.resizeProducer(context.Background(), producer.ID, thumbnailURL, headerURL)
	}()
	return nil
}

func (s *service) UpdateProducerEmail(ctx context.Context, in *user.UpdateProducerEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	producer, err := s.db.Producer.Get(ctx, in.ProducerID)
	if err != nil {
		return internalError(err)
	}
	params := &cognito.AdminChangeEmailParams{
		Username: producer.CognitoID,
		Email:    in.Email,
	}
	if err := s.adminAuth.AdminChangeEmail(ctx, params); err != nil {
		return internalError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, in.ProducerID, in.Email)
	return internalError(err)
}

func (s *service) UpdateProducerThumbnails(ctx context.Context, in *user.UpdateProducerThumbnailsInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Producer.UpdateThumbnails(ctx, in.ProducerID, in.Thumbnails)
	return internalError(err)
}

func (s *service) UpdateProducerHeaders(ctx context.Context, in *user.UpdateProducerHeadersInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Producer.UpdateHeaders(ctx, in.ProducerID, in.Headers)
	return internalError(err)
}

func (s *service) ResetProducerPassword(ctx context.Context, in *user.ResetProducerPasswordInput) error {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	producer, err := s.db.Producer.Get(ctx, in.ProducerID)
	if err != nil {
		return internalError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminChangePasswordParams{
		Username:  producer.CognitoID,
		Password:  password,
		Permanent: true,
	}
	if err := s.adminAuth.AdminChangePassword(ctx, params); err != nil {
		return internalError(err)
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

func (s *service) RelateProducers(ctx context.Context, in *user.RelateProducersInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	_, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid coordinator id: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	producers, err := s.db.Producer.MultiGet(ctx, in.ProducerIDs)
	if err != nil {
		return internalError(err)
	}
	producers = producers.Unrelated()
	if len(producers) != len(in.ProducerIDs) {
		return fmt.Errorf("api: contains invalid producers: %w", exception.ErrFailedPrecondition)
	}
	err = s.db.Producer.UpdateRelationship(ctx, in.CoordinatorID, in.ProducerIDs...)
	return internalError(err)
}

func (s *service) UnrelateProducer(ctx context.Context, in *user.UnrelateProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Producer.UpdateRelationship(ctx, "", in.ProducerID)
	return internalError(err)
}

func (s *service) DeleteProducer(ctx context.Context, in *user.DeleteProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Producer.Delete(ctx, in.ProducerID, s.deleteCognitoAdmin(in.ProducerID))
	return internalError(err)
}

func (s *service) resizeProducer(ctx context.Context, producerID, thumbnailURL, headerURL string) {
	s.waitGroup.Add(2)
	go func() {
		defer s.waitGroup.Done()
		if thumbnailURL == "" {
			return
		}
		in := &media.ResizeFileInput{
			TargetID: producerID,
			URLs:     []string{thumbnailURL},
		}
		if err := s.media.ResizeProducerThumbnail(ctx, in); err != nil {
			s.logger.Error("Failed to resize producer thumbnail",
				zap.String("producerId", producerID), zap.Error(err),
			)
		}
	}()
	go func() {
		defer s.waitGroup.Done()
		if headerURL == "" {
			return
		}
		in := &media.ResizeFileInput{
			TargetID: producerID,
			URLs:     []string{headerURL},
		}
		if err := s.media.ResizeProducerHeader(ctx, in); err != nil {
			s.logger.Error("Failed to resize producer header",
				zap.String("producerId", producerID), zap.Error(err),
			)
		}
	}()
}
