package entity

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	set "github.com/and-period/furumaru/api/pkg/set/v2"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

var (
	ErrTooLargeFileSize  = errors.New("entity: too large file size")
	ErrInvalidFileFormat = errors.New("entity: invalid file format")
)

const (
	CoordinatorThumbnailPath = "coordinators/thumbnail" // コーディネータサムネイル画像
	CoordinatorHeaderPath    = "coordinators/header"    // コーディネータヘッダー画像
	ProducerThumbnailPath    = "producers/thumbnail"    // 生産者サムネイル画像
	ProducerHeaderPath       = "producers/header"       // 生産者ヘッダー画像
	ProductMediaPath         = "products/media"         // 商品メディア
	ProductMediaImagePath    = "products/media/image"   // 商品メディア(画像)
	ProductMediaVideoPath    = "products/media/video"   // 商品メディア(映像)
	ProductTypeIconPath      = "product-types/icon"     // 品目アイコン
)

// Regulation - ファイルアップロード制約
type Regulation struct {
	MaxSize int64            // ファイルサイズ上限
	Formats *set.Set[string] // ファイル形式
	dir     string           // 保管先ディレクトリPath
}

var (
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
)

func (r *Regulation) Validate(file io.Reader, header *multipart.FileHeader) error {
	if file == nil || header == nil {
		return fmt.Errorf("entity: file and header is required: %w", ErrInvalidFileFormat)
	}
	if !r.validateSize(header) {
		return fmt.Errorf("%w: size=%d", ErrTooLargeFileSize, header.Size)
	}
	return r.validateFormat(file)
}

func (r *Regulation) validateSize(header *multipart.FileHeader) bool {
	return header.Size <= r.MaxSize
}

func (r *Regulation) validateFormat(file io.Reader) error {
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

func (r *Regulation) GenerateFilePath(header *multipart.FileHeader) string {
	key := uuid.Base58Encode(uuid.New())
	extension := filepath.Ext(header.Filename)
	filename := strings.Join([]string{key, extension}, "")
	return strings.Join([]string{r.dir, filename}, "/")
}
