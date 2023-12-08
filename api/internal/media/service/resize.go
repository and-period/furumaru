package service

import (
	"context"
	"encoding/json"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) ResizeCoordinatorThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeCoordinatorThumbnail)
}

func (s *service) ResizeCoordinatorHeader(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeCoordinatorHeader)
}

func (s *service) ResizeProducerThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeProducerThumbnail)
}

func (s *service) ResizeProducerHeader(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeProducerHeader)
}

func (s *service) ResizeUserThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeUserThumbnail)
}

func (s *service) ResizeProductMedia(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeProductMedia)
}

func (s *service) ResizeProductTypeIcon(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeProductTypeIcon)
}

func (s *service) ResizeScheduleThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	return s.sendResizeMessage(ctx, in, entity.FileTypeScheduleThumbnail)
}

func (s *service) sendResizeMessage(ctx context.Context, in *media.ResizeFileInput, fileType entity.FileType) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	payload := &entity.ResizerPayload{
		TargetID: in.TargetID,
		FileType: fileType,
		URLs:     in.URLs,
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return internalError(err)
	}
	_, err = s.producer.SendMessage(ctx, buf)
	return internalError(err)
}
