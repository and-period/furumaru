package service

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) GetUploadEvent(ctx context.Context, in *media.GetUploadEventInput) (*entity.UploadEvent, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	event := &entity.UploadEvent{Key: in.Key}
	if err := s.cache.Get(ctx, event); err != nil {
		return nil, internalError(err)
	}
	return event, nil
}

/**
 * ライブ配信関連
 */
func (s *service) GetBroadcastArchiveMP4UploadURL(
	ctx context.Context, in *media.GenerateBroadcastArchiveMP4UploadInput,
) (*entity.UploadEvent, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusDisabled {
		return nil, fmt.Errorf("service: this broadcast is not disabled: %w", exception.ErrFailedPrecondition)
	}
	return s.generateUploadURL(ctx, &in.GenerateUploadURLInput, entity.BroadcastArchiveRegulation, in.ScheduleID)
}

func (s *service) GetBroadcastLiveMP4UploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.BroadcastLiveMP4Regulation)
}

/**
 * コーディネータ関連
 */
func (s *service) GetCoordinatorThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorThumbnailRegulation)
}

func (s *service) GetCoordinatorHeaderUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorHeaderRegulation)
}

func (s *service) GetCoordinatorPromotionVideoUploadURL(
	ctx context.Context, in *media.GenerateUploadURLInput,
) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorPromotionVideoRegulation)
}

func (s *service) GetCoordinatorBonusVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorBonusVideoRegulation)
}

/**
 * 生産者関連
 */
func (s *service) GetProducerThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerThumbnailRegulation)
}

func (s *service) GetProducerHeaderUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerHeaderRegulation)
}

func (s *service) GetProducerPromotionVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerPromotionVideoRegulation)
}

func (s *service) GetProducerBonusVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerBonusVideoRegulation)
}

/**
 * 購入者関連
 */
func (s *service) GetUserThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.UserThumbnailRegulation)
}

/**
 * 商品関連
 */
func (s *service) GetProductMediaImageUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ProductMediaImageRegulation)
}

func (s *service) GetProductMediaVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ProductMediaVideoRegulation)
}

/**
 * 品目関連
 */
func (s *service) GetProductTypeIconUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ProductTypeIconRegulation)
}

/**
 * 開始スケジュール関連
 */
func (s *service) GetScheduleThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ScheduleThumbnailRegulation)
}

func (s *service) GetScheduleImageUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ScheduleImageRegulation)
}

func (s *service) GetScheduleOpeningVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (*entity.UploadEvent, error) {
	return s.generateUploadURL(ctx, in, entity.ScheduleOpeningVideoRegulation)
}

/**
 * private
 */
func (s *service) generateUploadURL(
	ctx context.Context,
	in *media.GenerateUploadURLInput,
	reg *entity.Regulation,
	keyArgs ...interface{},
) (*entity.UploadEvent, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	key, err := reg.GetObjectKey(in.FileType, keyArgs...)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get object key: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	const expiresIn = 10 * time.Minute
	url, err := s.tmp.GeneratePresignUploadURI(key, expiresIn)
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.UploadEventParams{
		Key:       key,
		FileGroup: reg.FileGroup(),
		FileType:  in.FileType,
		UploadURL: url,
		Now:       s.now(),
		TTL:       s.uploadEventTTL,
	}
	event := entity.NewUploadEvent(params)
	if err := s.cache.Insert(ctx, event); err != nil {
		return nil, internalError(err)
	}
	return event, nil
}
