package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) UploadCoordinatorThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorThumbnailPath)
}

func (s *service) UploadCoordinatorHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorHeaderPath)
}

func (s *service) UploadCoordinatorPromotionVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorPromotionVideoPath)
}

func (s *service) UploadCoordinatorBonusVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorBonusVideoPath)
}

func (s *service) UploadProducerThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerThumbnailPath)
}

func (s *service) UploadProducerHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerHeaderPath)
}

func (s *service) UploadProducerPromotionVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerPromotionVideoPath)
}

func (s *service) UploadProducerBonusVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerBonusVideoPath)
}

func (s *service) UploadUserThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.UserThumbnailPath)
}

func (s *service) UploadProductMedia(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProductMediaPath)
}

func (s *service) UploadProductTypeIcon(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProductTypeIconPath)
}

func (s *service) UploadScheduleThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ScheduleThumbnailPath)
}

func (s *service) UploadScheduleImage(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ScheduleImagePath)
}

func (s *service) UploadScheduleOpeningVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ScheduleOpeningVideoPath)
}

func (s *service) uploadFile(ctx context.Context, in *media.UploadFileInput, prefix string) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	u, err := s.parseURL(in, prefix)
	if err != nil {
		return "", fmt.Errorf("%s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	var url string
	switch u.Host {
	case s.tmpURL().Host:
		url, err = s.uploadPermanentFile(ctx, u)
	case s.storageURL().Host:
		url, err = s.downloadFile(ctx, u)
	default:
		return "", fmt.Errorf("service: unknown storage host. host=%s: %w", u.Host, exception.ErrInvalidArgument)
	}
	return url, internalError(err)
}

func (s *service) parseURL(in *media.UploadFileInput, prefix string) (*url.URL, error) {
	u, err := url.Parse(in.URL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errParseURL, err.Error())
	}
	if !strings.Contains(u.Path, prefix) {
		return nil, fmt.Errorf("%w. url=%s", errInvalidURL, in.URL)
	}
	return u, nil
}

func (s *service) uploadPermanentFile(ctx context.Context, u *url.URL) (string, error) {
	file, err := s.tmp.Download(ctx, u.String())
	if err != nil {
		return "", err
	}
	path := strings.TrimPrefix(u.Path, "/") // url.URLから取得したPathは / から始まるため
	return s.storage.Upload(ctx, path, file)
}

func (s *service) downloadFile(ctx context.Context, u *url.URL) (string, error) {
	url := u.String()
	if _, err := s.storage.Download(ctx, url); err != nil {
		return "", err
	}
	return url, nil
}
