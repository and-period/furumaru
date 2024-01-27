package entity

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
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
	BroadcastArchiveRegulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     BroadcastArchiveMP4Path,
	}
	BroadcastLiveMP4Regulation = &Regulation{
		MaxSize: 200 << 20, // 200MB
		Formats: set.New("video/mp4"),
		dir:     BroadcastLiveMP4Path,
	}
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
	UserThumbnailRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     UserThumbnailPath,
	}
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
	ProductTypeIconRegulation = &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     ProductTypeIconPath,
	}
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

func FindByObjectKey(key string) (*Regulation, error) {
	dir := path.Dir(key)
	if strings.HasPrefix(dir, BroadcastArchivePath) {
		return BroadcastArchiveRegulation, nil
	}
	switch dir {
	case BroadcastLiveMP4Path:
		return BroadcastLiveMP4Regulation, nil
	case CoordinatorThumbnailPath:
		return CoordinatorThumbnailRegulation, nil
	case CoordinatorHeaderPath:
		return CoordinatorHeaderRegulation, nil
	case CoordinatorPromotionVideoPath:
		return CoordinatorPromotionVideoRegulation, nil
	case CoordinatorBonusVideoPath:
		return CoordinatorBonusVideoRegulation, nil
	case ProducerThumbnailPath:
		return ProducerThumbnailRegulation, nil
	case ProducerHeaderPath:
		return ProducerHeaderRegulation, nil
	case ProducerPromotionVideoPath:
		return ProducerPromotionVideoRegulation, nil
	case ProducerBonusVideoPath:
		return ProducerBonusVideoRegulation, nil
	case UserThumbnailPath:
		return UserThumbnailRegulation, nil
	case ProductMediaImagePath:
		return ProductMediaImageRegulation, nil
	case ProductMediaVideoPath:
		return ProductMediaVideoRegulation, nil
	case ProductTypeIconPath:
		return ProductTypeIconRegulation, nil
	case ScheduleThumbnailPath:
		return ScheduleThumbnailRegulation, nil
	case ScheduleImagePath:
		return ScheduleImageRegulation, nil
	case ScheduleOpeningVideoPath:
		return ScheduleOpeningVideoRegulation, nil
	default:
		return nil, ErrNotFoundReguration
	}
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
	if r.Formats.Contains(contentType) {
		return fmt.Errorf("%w: content type=%s", ErrInvalidFileFormat, contentType)
	}
	return nil
}

// Deprecated: Use to Validate
func (r *Regulation) ValidateV1(file io.Reader, header *multipart.FileHeader) error {
	if file == nil || header == nil {
		return fmt.Errorf("entity: file and header is required: %w", ErrInvalidFileFormat)
	}
	if !r.validateV1Size(header) {
		return fmt.Errorf("%w: size=%d", ErrTooLargeFileSize, header.Size)
	}
	return r.validateV1Format(file)
}

func (r *Regulation) validateV1Size(header *multipart.FileHeader) bool {
	return header.Size <= r.MaxSize
}

func (r *Regulation) validateV1Format(file io.Reader) error {
	if r.Formats.Len() == 0 {
		return nil
	}
	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(buf)
	if r.Formats.Contains(contentType) {
		return nil
	}
	return fmt.Errorf("%w: content type=%s", ErrInvalidFileFormat, contentType)
}

func (r *Regulation) GenerateFilePath(header *multipart.FileHeader, args ...interface{}) string {
	key := uuid.Base58Encode(uuid.New())
	extension := strings.ToLower(filepath.Ext(header.Filename))
	dirname := fmt.Sprintf(r.dir, args...)
	filename := strings.Join([]string{key, extension}, "")
	return strings.Join([]string{dirname, filename}, "/")
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
