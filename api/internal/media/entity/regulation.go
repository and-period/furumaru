package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
	BroadcastArchiveTextPath      = "schedules/archives/%s/text"   // ライブ配信後のアーカイブ音声文字起こし
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
	ScheduleThumbnailPath         = "schedules/thumbnail"          // 開催スケジュールサムネイル画像
	ScheduleImagePath             = "schedules/image"              // 開催スケジュール蓋絵
	ScheduleOpeningVideoPath      = "schedules/opening-video"      // 開催スケジュールオープニング動画
	ExperienceMediaImagePath      = "experiences/media/image"      // 体験メディア(画像)
	ExperienceMediaVideoPath      = "experiences/media/video"      // 体験メディア(映像)
	ExperiencePromotionVideoPath  = "experiences/promotion-video"  // 体験紹介映像
	VideoThumbnailPath            = "videos/thumbnail"             // オンデマンド配信サムネイル画像
	VideoMP4Path                  = "videos/mp4"                   // オンデマンド配信動画(mp4)
	SpotThumbnailPath             = "spots/thumbnail"              // スポットサムネイル画像
)

const defaultCacheTTL = 14 * 24 * time.Hour // 2週間

// ConversionType - ファイルの変換種別
type ConversionType int32

const (
	ConversionTypeNone      ConversionType = iota // 変換不要
	ConversionTypeJPEGToPNG                       // 画像の変換(JPEG -> PNG)
)

// Regulation - ファイルアップロード制約
type Regulation struct {
	MaxSize        int64            // ファイルサイズ上限
	Formats        *set.Set[string] // ファイル形式
	ConversionType ConversionType   // ファイル変換が必要な場合の変換種別
	CacheTTL       time.Duration    // キャッシュの有効期限
	dir            string           // 保管先ディレクトリPath
}

var (
	// ライブ配信関連
	BroadcastArchiveRegulation = &Regulation{
		MaxSize:  3 << 30, // 3GB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      BroadcastArchiveMP4Path,
	}
	BroadcastLiveMP4Regulation = &Regulation{
		MaxSize:  3 << 30, // 3GB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      BroadcastLiveMP4Path,
	}
	BroadcastArchiveTextRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("text/vtt"),
		CacheTTL: defaultCacheTTL,
		dir:      BroadcastArchiveTextPath,
	}
	// コーディネータ関連
	CoordinatorThumbnailRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      CoordinatorThumbnailPath,
	}
	CoordinatorHeaderRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      CoordinatorHeaderPath,
	}
	CoordinatorPromotionVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      CoordinatorPromotionVideoPath,
	}
	CoordinatorBonusVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      CoordinatorBonusVideoPath,
	}
	// 生産者関連
	ProducerThumbnailRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      ProducerThumbnailPath,
	}
	ProducerHeaderRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      ProducerHeaderPath,
	}
	ProducerPromotionVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      ProducerPromotionVideoPath,
	}
	ProducerBonusVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      ProducerBonusVideoPath,
	}
	// 購入者関連
	UserThumbnailRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      UserThumbnailPath,
	}
	// 商品関連
	ProductMediaImageRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      ProductMediaImagePath,
	}
	ProductMediaVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      ProductMediaVideoPath,
	}
	// 品目関連
	ProductTypeIconRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      ProductTypeIconPath,
	}
	// 開催スケジュール関連
	ScheduleThumbnailRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      ScheduleThumbnailPath,
	}
	ScheduleImageRegulation = &Regulation{
		MaxSize:        10 << 20, // 10MB
		Formats:        set.New("image/png", "image/jpeg"),
		ConversionType: ConversionTypeJPEGToPNG, // MediaLiveの仕様に合わせてPNG形式に変換
		CacheTTL:       defaultCacheTTL,
		dir:            ScheduleImagePath,
	}
	ScheduleOpeningVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      ScheduleOpeningVideoPath,
	}
	// 体験関連
	ExperienceMediaImageRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      ExperienceMediaImagePath,
	}
	ExperienceMediaVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      ExperienceMediaVideoPath,
	}
	ExperiencePromotionVideoRegulation = &Regulation{
		MaxSize:  200 << 20, // 200MB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      ExperiencePromotionVideoPath,
	}
	// オンデマンド配信関連
	VideoThumbnailRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      VideoThumbnailPath,
	}
	VideoMP4Regulation = &Regulation{
		MaxSize:  3 << 30, // 3GB
		Formats:  set.New("video/mp4"),
		CacheTTL: defaultCacheTTL,
		dir:      VideoMP4Path,
	}
	// スポット関連
	SpotThumbnailRegulation = &Regulation{
		MaxSize:  10 << 20, // 10MB
		Formats:  set.New("image/png", "image/jpeg"),
		CacheTTL: defaultCacheTTL,
		dir:      SpotThumbnailPath,
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
	case "text/vtt":
		return "vtt", nil
	default:
		return "", ErrUnknownContentType
	}
}

func (r *Regulation) ShouldConvert(contentType string) bool {
	switch r.ConversionType {
	case ConversionTypeJPEGToPNG:
		return contentType == "image/jpeg"
	case ConversionTypeNone:
		return false
	default:
		return false
	}
}
