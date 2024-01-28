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

func (s *service) GetUploadEvent(ctx context.Context, in *media.GetUploadEventInput) (*entity.UploadEvent, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	u, err := url.Parse(in.UploadURL)
	if err != nil {
		return nil, fmt.Errorf("service: failed to parse url: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	event := &entity.UploadEvent{Key: strings.TrimPrefix(u.Path, "/")}
	if err := s.cache.Get(ctx, event); err != nil {
		return nil, internalError(err)
	}
	return event, nil
}

// Deprecated
func (s *service) UploadCoordinatorThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorThumbnailPath)
}

// Deprecated
func (s *service) UploadCoordinatorHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorHeaderPath)
}

// Deprecated
func (s *service) UploadCoordinatorPromotionVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorPromotionVideoPath)
}

// Deprecated
func (s *service) UploadCoordinatorBonusVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorBonusVideoPath)
}

// Deprecated
func (s *service) UploadProducerThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerThumbnailPath)
}

// Deprecated
func (s *service) UploadProducerHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerHeaderPath)
}

// Deprecated
func (s *service) UploadProducerPromotionVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerPromotionVideoPath)
}

// Deprecated
func (s *service) UploadProducerBonusVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerBonusVideoPath)
}

// Deprecated
func (s *service) UploadUserThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.UserThumbnailPath)
}

// Deprecated
func (s *service) UploadProductMedia(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProductMediaPath)
}

// Deprecated
func (s *service) UploadProductTypeIcon(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProductTypeIconPath)
}

// Deprecated
func (s *service) UploadScheduleThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ScheduleThumbnailPath)
}

// Deprecated
func (s *service) UploadScheduleImage(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ScheduleImagePath)
}

// Deprecated
func (s *service) UploadScheduleOpeningVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ScheduleOpeningVideoPath)
}

// Deprecated
func (s *service) uploadFile(ctx context.Context, in *media.UploadFileInput, prefix string) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	u, err := s.parseURL(in, prefix)
	if err != nil {
		return "", fmt.Errorf("%s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	var url string
	switch {
	case s.tmp.IsMyHost(u.String()):
		url, err = s.uploadPermanentFile(ctx, u)
	case s.storage.IsMyHost(u.String()):
		url, err = s.downloadFile(ctx, u)
	default:
		return "", fmt.Errorf("service: unknown storage host. host=%s: %w", u.Host, exception.ErrInvalidArgument)
	}
	return url, internalError(err)
}

// Deprecated
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

// Deprecated
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
