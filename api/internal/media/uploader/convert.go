package uploader

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/and-period/furumaru/api/internal/media/entity"
	"go.uber.org/zap"
)

func (u *uploader) uploadConvetFile(ctx context.Context, event *entity.UploadEvent, reguration *entity.Regulation) (string, error) {
	if reguration.Convert == entity.ConvertTypeNone {
	}
	switch reguration.Convert {
	case entity.ConvertTypeJPEGToPNG:
		return u.convertJPEGToPNG(ctx, event)
	default:
		u.logger.Warn("Unsupported convert type", zap.String("key", event.Key), zap.Int32("convertType", int32(reguration.Convert)))
		return event.Key, nil // 変換できないファイルに対してはエラーにせず、元ファイルをそのまま利用する
	}
}

func (u *uploader) convertJPEGToPNG(ctx context.Context, event *entity.UploadEvent) (string, error) {
	f, err := u.tmp.Download(ctx, event.Key)
	if err != nil {
		return "", fmt.Errorf("uploader: failed to download file: %w", err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return "", fmt.Errorf("uploader: failed to decode image: %w", err)
	}
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		return "", fmt.Errorf("uploader: failed to encode image: %w", err)
	}
	const extension = ".png"
	key := u.replaceExt(event.Key, extension)
	if _, err := u.storage.Upload(ctx, key, buf); err != nil {
		return "", fmt.Errorf("uploader: failed to upload file: %w", err)
	}
	return key, nil
}

func (u *uploader) replaceExt(filePath, ext string) string {
	current := filepath.Ext(filePath)
	return strings.Replace(filePath, current, ext, 1)
}
