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

func (s *service) GetCoordinatorThumbnailUploadURL(_ context.Context, in *media.GenerateUploadURLInput) (string, error) {
	return s.generateUploadURL(in, entity.CoordinatorThumbnailRegulation)
}

// Deprecated
func (s *service) GenerateCoordinatorThumbnail(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.CoordinatorThumbnailRegulation)
}

func (s *service) GenerateCoordinatorHeader(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.CoordinatorHeaderRegulation)
}

func (s *service) GenerateCoordinatorPromotionVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.CoordinatorPromotionVideoRegulation)
}

func (s *service) GenerateCoordinatorBonusVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.CoordinatorBonusVideoRegulation)
}

func (s *service) GenerateProducerThumbnail(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ProducerThumbnailRegulation)
}

func (s *service) GenerateProducerHeader(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ProducerHeaderRegulation)
}

func (s *service) GenerateProducerPromotionVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ProducerPromotionVideoRegulation)
}

func (s *service) GenerateUserThumbnail(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.UserThumbnailRegulation)
}

func (s *service) GenerateProducerBonusVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ProducerBonusVideoRegulation)
}

func (s *service) GenerateProductMediaImage(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ProductMediaImageRegulation)
}

func (s *service) GenerateProductMediaVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ProductMediaVideoRegulation)
}

func (s *service) GenerateProductTypeIcon(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ProductTypeIconRegulation)
}

func (s *service) GenerateScheduleThumbnail(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ScheduleThumbnailRegulation)
}

func (s *service) GenerateScheduleImage(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ScheduleImageRegulation)
}

func (s *service) GenerateScheduleOpeningVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	return s.generateFile(ctx, in, entity.ScheduleOpeningVideoRegulation)
}

func (s *service) generateUploadURL(in *media.GenerateUploadURLInput, reg *entity.Regulation) (string, error) {
	const expiresIn = 10 * time.Minute
	key, err := reg.GetObjectKey(in.FileType)
	if err != nil {
		return "", fmt.Errorf("service: failed to get object key: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	url, err := s.tmp.GeneratePresignUploadURI(key, expiresIn)
	return url, internalError(err)
}

func (s *service) generateFile(
	ctx context.Context, in *media.GenerateFileInput, reg *entity.Regulation,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	var buf bytes.Buffer
	teeReader := io.TeeReader(in.File, &buf)
	if err := reg.ValidateV1(teeReader, in.Header); err != nil {
		return "", fmt.Errorf("%w: %s", exception.ErrInvalidArgument, err.Error())
	}
	path := reg.GenerateFilePath(in.Header)
	url, err := s.tmp.Upload(ctx, path, &buf)
	return url, internalError(err)
}
