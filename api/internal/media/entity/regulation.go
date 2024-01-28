package entity

import (
	"errors"
	"fmt"
	"strings"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

var (
	ErrNotFoundReguration     = errors.New("entity: not found reguration")
	ErrTooLargeFileSize       = errors.New("entity: too large file size")
	ErrInvalidFileFormat      = errors.New("entity: invalid file format")
	ErrUnsupportedContentType = errors.New("entity: unsupported content type")
	ErrUnknownContentType     = errors.New("entity: unknown content type")
)

const (
	BroadcastLiveMP4Path          = "schedules/lives"              // ライブ配信中に使用する動画
	BroadcastArchivePath          = "schedules/archives"           // ライブ配信後のアーカイブ動画
	BroadcastArchiveMP4Path       = "schedules/archives/%s/mp4"    // ライブ配信後のアーカイブ動画(mp4)
	BroadcastArchiveHLSPath       = "schedules/archives/%s/hls"    // ライブ配信後のアーカイブ動画(hls)
	CoordinatorThumbnailPath      = "coordinators/thumbnail"       // コーディネータサムネイル画像
	CoordinatorHeaderPath         = "coordinators/header"          // コーディネータヘッダー画像
	CoordinatorPromotionVideoPath = "coordinators/promotion-video" // コーディネータ紹介映像
	CoordinatorBonusVideoPath     = "coordinators/bonus-video"     // コーディネータ購入特典映像
	ProducerThumbnailPath         = "producers/thumbnail"          // 生産者サムネイル画像
	ProducerHeaderPath            = "producers/header"             // 生産者ヘッダー画像
	ProducerPromotionVideoPath    = "producers/promotion-video"    // 生産者紹介映像
	ProducerBonusVideoPath        = "producers/bonus-video"        // 生産者購入特典映像
	UserThumbnailPath             = "users/thumbnail"              // 購入者サムネイル画像
	ProductMediaPath              = "products/media"               // 商品メディア
	ProductMediaImagePath         = "products/media/image"         // 商品メディア(画像)
	ProductMediaVideoPath         = "products/media/video"         // 商品メディア(映像)
	ProductTypeIconPath           = "product-types/icon"           // 品目アイコン
	ScheduleThumbnailPath         = "schedules/thumbnail"          // 開催スケジュールサムネイル
	ScheduleImagePath             = "schedules/image"              // 開催スケジュール蓋絵
	ScheduleOpeningVideoPath      = "schedules/opening-video"      // 開催スケジュールオープニング動画
)

// Regulation - ファイルアップロード制約
type Regulation struct {
	MaxSize int64            // ファイルサイズ上限
	Formats *set.Set[string] // ファイル形式
	dir     string           // 保管先ディレクトリPath
}

var (
	// ライブ配信関連
	BroadcastArchiveRegulation = &Regulation{
		MaxSize: 2 << 30, // 2GB
		Formats: set.New("video/mp4"),
		dir:     BroadcastArchiveMP4Path,
	}
	BroadcastLiveMP4Regulation = &Regulation{
		MaxSize: 2 << 30, // 2GB
		Formats: set.New("video/mp4"),
		dir:     BroadcastLiveMP4Path,
	}
	// コーディネータ関連
	CoordinatorThumbnailRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     CoordinatorThumbnailPath,
	}
	CoordinatorHeaderRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     CoordinatorHeaderPath,
	}
	CoordinatorPromotionVideoRegulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     CoordinatorPromotionVideoPath,
	}
	CoordinatorBonusVideoRegulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     CoordinatorBonusVideoPath,
	}
	// 生産者関連
	ProducerThumbnailRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     ProducerThumbnailPath,
	}
	ProducerHeaderRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     ProducerHeaderPath,
	}
	ProducerPromotionVideoRegulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     ProducerPromotionVideoPath,
	}
	ProducerBonusVideoRegulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     ProducerBonusVideoPath,
	}
	// 購入者関連
	UserThumbnailRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     UserThumbnailPath,
	}
	// 商品関連
	ProductMediaImageRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     ProductMediaImagePath,
	}
	ProductMediaVideoRegulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     ProductMediaVideoPath,
	}
	// 品目関連
	ProductTypeIconRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     ProductTypeIconPath,
	}
	// 開催スケジュール関連
	ScheduleThumbnailRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     ScheduleThumbnailPath,
	}
	ScheduleImageRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png"),
		dir:     ScheduleImagePath,
	}
	ScheduleOpeningVideoRegulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     ScheduleOpeningVideoPath,
	}
)

func (r *Regulation) FileGroup() string {
	return r.dir
}

func (r *Regulation) Validate(contentType string, size int64) error {
	if err := r.validateSize(size); err != nil {
		return err
	}
	return r.validateFormat(contentType)
}

func (r *Regulation) validateSize(size int64) error {
	if size > r.MaxSize {
		return fmt.Errorf("%w: size=%d", ErrTooLargeFileSize, size)
	}
	return nil
}

func (r *Regulation) validateFormat(contentType string) error {
	if !r.Formats.Contains(contentType) {
		return fmt.Errorf("%w: content type=%s", ErrInvalidFileFormat, contentType)
	}
	return nil
}

func (r *Regulation) GetObjectKey(contentType string, args ...interface{}) (string, error) {
	ext, err := r.GetFileExtension(contentType)
	if err != nil {
		return "", err
	}
	key := uuid.Base58Encode(uuid.New())
	dirname := fmt.Sprintf(r.dir, args...)
	filename := strings.Join([]string{key, ext}, ".")
	return strings.Join([]string{dirname, filename}, "/"), nil
}

// ref: https://developer.mozilla.org/ja/docs/Web/HTTP/Basics_of_HTTP/MIME_types
func (r *Regulation) GetFileExtension(contentType string) (string, error) {
	if !r.Formats.Contains(contentType) {
		return "", ErrUnsupportedContentType
	}
	switch contentType {
	// 画像タイプ
	case "image/jpeg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "video/mp4":
		return "mp4", nil
	default:
		return "", ErrUnknownContentType
	}
}
