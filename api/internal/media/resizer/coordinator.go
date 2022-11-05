package resizer

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media/entity"
	"go.uber.org/zap"
)

func (r *resizer) coordinatorThumbnail(ctx context.Context, payload *entity.ResizerPayload) error {
	url := payload.URLs[0]
	if url == "" {
		return errRequiredMediaURL
	}
	buf, err := r.storage.Download(ctx, payload.URLs[0])
	if err != nil {
		return err
	}
	resizedImages, err := r.resizeImages(url, buf)
	if err != nil {
		return err
	}
	images, err := r.uploadImages(ctx, url, resizedImages)
	if err != nil {
		return err
	}
	// データの更新処理
	r.logger.Info("Success upload", zap.Any("payload", payload), zap.Any("images", images))
	return nil
}

func (r *resizer) coordinatorHeader(ctx context.Context, payload *entity.ResizerPayload) error {
	url := payload.URLs[0]
	if url == "" {
		return errRequiredMediaURL
	}
	buf, err := r.storage.Download(ctx, payload.URLs[0])
	if err != nil {
		return err
	}
	resizedImages, err := r.resizeImages(url, buf)
	if err != nil {
		return err
	}
	images, err := r.uploadImages(ctx, url, resizedImages)
	if err != nil {
		return err
	}
	// データの更新処理
	r.logger.Info("Success upload", zap.Any("payload", payload), zap.Any("images", images))
	return nil
}
