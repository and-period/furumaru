package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) GetCoordinatorThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorThumbnailRegulation)
}

func (s *service) GetCoordinatorHeaderUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorHeaderRegulation)
}

func (s *service) GetCoordinatorPromotionVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorPromotionVideoRegulation)
}

func (s *service) GetCoordinatorBonusVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.CoordinatorBonusVideoRegulation)
}

func (s *service) GetProducerThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerThumbnailRegulation)
}

func (s *service) GetProducerHeaderUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerHeaderRegulation)
}

func (s *service) GetProducerPromotionVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerPromotionVideoRegulation)
}

func (s *service) GetProducerBonusVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ProducerBonusVideoRegulation)
}

func (s *service) GetUserThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.UserThumbnailRegulation)
}

func (s *service) GetProductMediaImageUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ProductMediaImageRegulation)
}

func (s *service) GetProductMediaVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ProductMediaVideoRegulation)
}

func (s *service) GetProductTypeIconUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ProductTypeIconRegulation)
}

func (s *service) GetScheduleThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ScheduleThumbnailRegulation)
}

func (s *service) GetScheduleImageUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ScheduleImageRegulation)
}

func (s *service) GetScheduleOpeningVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(ctx, in, entity.ScheduleOpeningVideoRegulation)
}

func (s *service) generateUploadURL(ctx context.Context, in *media.GenerateUploadURLInput, reg *entity.Regulation) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	key, err := reg.GetObjectKey(in.FileType)
	if err != nil {
		return "", fmt.Errorf("service: failed to get object key: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	const expiresIn = 10 * time.Minute
	url, err := s.tmp.GeneratePresignUploadURI(key, expiresIn)
	if err != nil {
		return "", internalError(err)
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
		return "", internalError(err)
	}
	return url, nil
}

func (s *service) generateFile(
	ctx context.Context, in *media.GenerateFileInput, reg *entity.Regulation,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	var buf bytes.Buffer
	teeReader := io.TeeReader(in.File, &buf)
	//nolint:staticcheck
	if err := reg.ValidateV1(teeReader, in.Header); err != nil {
		return "", fmt.Errorf("%w: %s", exception.ErrInvalidArgument, err.Error())
	}
	path := reg.GenerateFilePath(in.Header)
	url, err := s.tmp.Upload(ctx, path, &buf)
	return url, internalError(err)
}
