package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListProducers(ctx context.Context, in *user.ListProducersInput) (entity.Producers, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListProducersParams{
		CoordinatorID: in.CoordinatorID,
		Name:          in.Name,
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
		return nil, 0, internalError(err)
	}
	return producers, total, nil
}

func (s *service) MultiGetProducers(ctx context.Context, in *user.MultiGetProducersInput) (entity.Producers, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	var (
		producers entity.Producers
		err       error
	)
	if in.WithDeleted {
		producers, err = s.db.Producer.MultiGetWithDeleted(ctx, in.ProducerIDs)
	} else {
		producers, err = s.db.Producer.MultiGet(ctx, in.ProducerIDs)
	}
	return producers, internalError(err)
}

func (s *service) GetProducer(ctx context.Context, in *user.GetProducerInput) (*entity.Producer, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	var (
		producer *entity.Producer
		err      error
	)
	if in.WithDeleted {
		producer, err = s.db.Producer.GetWithDeleted(ctx, in.ProducerID)
	} else {
		producer, err = s.db.Producer.Get(ctx, in.ProducerID)
	}
	return producer, internalError(err)
}

func (s *service) CreateProducer(ctx context.Context, in *user.CreateProducerInput) (*entity.Producer, error) {
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
	adminParams := &entity.NewAdminParams{
		CognitoID:     "", // 生産者は認証機能を持たせない
		Type:          entity.AdminTypeProducer,
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
		PrefectureCode:    in.PrefectureCode,
		City:              in.City,
		AddressLine1:      in.AddressLine1,
		AddressLine2:      in.AddressLine2,
	}
	producer, err := entity.NewProducer(params)
	if err != nil {
		return nil, fmt.Errorf("service: failed to new producer: %w: %s", exception.ErrInvalidArgument, err.Error())
	}
	auth := func(_ context.Context) error {
		return nil // 生産者は認証機能を持たないため何もしない
	}
	if err := s.db.Producer.Create(ctx, producer, auth); err != nil {
		return nil, internalError(err)
	}
	return producer, nil
}

func (s *service) UpdateProducer(ctx context.Context, in *user.UpdateProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if in.PrefectureCode > 0 {
		if _, err := codes.ToPrefectureJapanese(in.PrefectureCode); err != nil {
			return fmt.Errorf("service: invalid prefecture code: %w: %s", exception.ErrInvalidArgument, err.Error())
		}
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
		Email:             in.Email,
		PhoneNumber:       in.PhoneNumber,
		PostalCode:        in.PostalCode,
		PrefectureCode:    in.PrefectureCode,
		City:              in.City,
		AddressLine1:      in.AddressLine1,
		AddressLine2:      in.AddressLine2,
	}
	if err := s.db.Producer.Update(ctx, in.ProducerID, params); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) DeleteProducer(ctx context.Context, in *user.DeleteProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	auth := func(_ context.Context) error {
		return nil // 生産者は認証機能を持たないため何もしない
	}
	err := s.db.Producer.Delete(ctx, in.ProducerID, auth)
	return internalError(err)
}
