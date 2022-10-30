package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"go.uber.org/zap"
)

func (s *service) UploadCoordinatorThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorThumbnailPath)
}

func (s *service) UploadCoordinatorHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.CoordinatorHeaderPath)
}

func (s *service) UploadProducerThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerThumbnailPath)
}

func (s *service) UploadProducerHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProducerHeaderPath)
}

func (s *service) UploadProductMedia(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProductMediaPath)
}

func (s *service) UploadProductTypeIcon(ctx context.Context, in *media.UploadFileInput) (string, error) {
	return s.uploadFile(ctx, in, entity.ProductTypeIconPath)
}

func (s *service) uploadFile(ctx context.Context, in *media.UploadFileInput, prefix string) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", exception.InternalError(err)
	}
	u, err := s.parseURL(in, prefix)
	if err != nil {
		return "", fmt.Errorf("%s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	// TODO: remove
	s.logger.Debug("upload file",
		zap.Any("input", in), zap.Any("url", u),
		zap.Any("storage", s.storageURL()), zap.String("storageHost", s.storageURL().Host),
		zap.Any("tmp", s.tmpURL()), zap.String("tmpHost", s.tmpURL().Host),
		zap.Bool("storageMatch", s.storageURL().Host == u.Host),
		zap.Bool("tmpMatch", s.tmpURL().Host == u.Host),
	)
	var path string
	switch u.Host {
	case s.tmpURL().Host:
		path, err = s.uploadPermanentFile(ctx, u)
	case s.storageURL().Host:
		path, err = s.downloadFile(ctx, u)
	default:
		err = fmt.Errorf("service: unknown storage host. host=%s: %w", u.Host, exception.ErrInvalidArgument)
	}
	return path, exception.InternalError(err)
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
	file, err := s.tmp.Download(ctx, u.Path)
	if err != nil {
		return "", err
	}
	return s.storage.Upload(ctx, u.Path, file)
}

func (s *service) downloadFile(ctx context.Context, u *url.URL) (string, error) {
	if _, err := s.storage.Download(ctx, u.Path); err != nil {
		return "", err
	}
	return u.String(), nil
}
