package entity

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegulation_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		regulation *Regulation
		input      func(t *testing.T) (io.Reader, *multipart.FileHeader)
		expect     error
	}{
		// CoordinatorThumbnail
		{
			name:       "success coordinator thumbnail",
			regulation: CoordinatorThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for coordinator thumbnail",
			regulation: CoordinatorThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testImageFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for coordinator thumbnail",
			regulation: CoordinatorThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for coordinator thumbnail",
			regulation: CoordinatorThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// CoordinatorHeader
		{
			name:       "success coordinator header",
			regulation: CoordinatorHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for coordinator header",
			regulation: CoordinatorHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testImageFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for coordinator header",
			regulation: CoordinatorHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for coordinator header",
			regulation: CoordinatorHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// CoordinatorPromotionVideo
		{
			name:       "success coordinator promotion video",
			regulation: CoordinatorPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for coordinator promotion video",
			regulation: CoordinatorPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testVideoFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for coordinator promotion video",
			regulation: CoordinatorPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for coordinator promotion video",
			regulation: CoordinatorPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// CoordinatorBonusVideo
		{
			name:       "success coordinator bonus video",
			regulation: CoordinatorBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for coordinator bonus video",
			regulation: CoordinatorBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testVideoFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for coordinator bonus video",
			regulation: CoordinatorBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for coordinator bonus video",
			regulation: CoordinatorBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// ProducerThumbnail
		{
			name:       "success producer thumbnail",
			regulation: ProducerThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for producer thumbnail",
			regulation: ProducerThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testImageFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for producer thumbnail",
			regulation: ProducerThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for producer thumbnail",
			regulation: ProducerThumbnailRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// ProducerHeader
		{
			name:       "success producer header",
			regulation: ProducerHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for producer header",
			regulation: ProducerHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testImageFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for producer header",
			regulation: ProducerHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for producer header",
			regulation: ProducerHeaderRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// ProducerPromotionVideo
		{
			name:       "success producer promotion video",
			regulation: ProducerPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for producer promotion video",
			regulation: ProducerPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testVideoFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for producer promotion video",
			regulation: ProducerPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for producer promotion video",
			regulation: ProducerPromotionVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// ProducerBonusVideo
		{
			name:       "success producer bonus video",
			regulation: ProducerBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for producer bonus video",
			regulation: ProducerBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testVideoFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for producer bonus video",
			regulation: ProducerBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for producer bonus video",
			regulation: ProducerBonusVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// ProductMediaImage
		{
			name:       "success product media image",
			regulation: ProductMediaImageRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for product media image",
			regulation: ProductMediaImageRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testImageFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for product media image",
			regulation: ProductMediaImageRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for product media image",
			regulation: ProductMediaImageRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// ProductMediaVideo
		{
			name:       "success product media video",
			regulation: ProductMediaVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for product media video",
			regulation: ProductMediaVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testVideoFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for product media video",
			regulation: ProductMediaVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for product media video",
			regulation: ProductMediaVideoRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
		// ProductTypeIcon
		{
			name:       "success product type icon",
			regulation: ProductTypeIconRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testImageFile(t)
			},
			expect: nil,
		},
		{
			name:       "required for product type icon",
			regulation: ProductTypeIconRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				_, header := testImageFile(t)
				return nil, header
			},
			expect: ErrInvalidFileFormat,
		},
		{
			name:       "invalid size for product type icon",
			regulation: ProductTypeIconRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return file, header
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name:       "invalid format for product type icon",
			regulation: ProductTypeIconRegulation,
			input: func(t *testing.T) (io.Reader, *multipart.FileHeader) {
				return testVideoFile(t)
			},
			expect: ErrInvalidFileFormat,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, tt.regulation.Validate(tt.input(t)), tt.expect)
		})
	}
}

func TestRegulation_GenerateFilePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		regulation *Regulation
		header     *multipart.FileHeader
		expect     string
	}{
		{
			name:       "coordinator thumbnail",
			regulation: CoordinatorThumbnailRegulation,
			header:     &multipart.FileHeader{Filename: "and-period.png"},
			expect:     "coordinators/thumbnail/[a-zA-Z0-9]+.png",
		},
		{
			name:       "coordinator header",
			regulation: CoordinatorHeaderRegulation,
			header:     &multipart.FileHeader{Filename: "and-period.png"},
			expect:     "coordinators/header/[a-zA-Z0-9]+.png",
		},
		{
			name:       "producer thumbnail",
			regulation: ProducerThumbnailRegulation,
			header:     &multipart.FileHeader{Filename: "and-period.png"},
			expect:     "producers/thumbnail/[a-zA-Z0-9]+.png",
		},
		{
			name:       "producer header",
			regulation: ProducerHeaderRegulation,
			header:     &multipart.FileHeader{Filename: "and-period.png"},
			expect:     "producers/header/[a-zA-Z0-9]+.png",
		},
		{
			name:       "product media image",
			regulation: ProductMediaImageRegulation,
			header:     &multipart.FileHeader{Filename: "and-period.png"},
			expect:     "products/media/image/[a-zA-Z0-9]+.png",
		},
		{
			name:       "product media video",
			regulation: ProductMediaVideoRegulation,
			header:     &multipart.FileHeader{Filename: "and-period.mp4"},
			expect:     "products/media/video/[a-zA-Z0-9]+.mp4",
		},
		{
			name:       "product type icon",
			regulation: ProductTypeIconRegulation,
			header:     &multipart.FileHeader{Filename: "and-period.png"},
			expect:     "product-types/icon/[a-zA-Z0-9]+.png",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Regexp(t, regexp.MustCompile(tt.expect), tt.regulation.GenerateFilePath(tt.header))
		})
	}
}

func testImageFile(t *testing.T) (io.Reader, *multipart.FileHeader) {
	const filename, format = "and-period.png", "image"

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, format, filename))
	header.Set("Content-Type", "multipart/form-data")
	part := &multipart.FileHeader{
		Filename: filepath,
		Header:   header,
		Size:     3 << 20, // 3MB
	}

	return file, part
}

func testVideoFile(t *testing.T) (io.Reader, *multipart.FileHeader) {
	const filename, format = "and-period.mp4", "video"

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, format, filename))
	header.Set("Content-Type", "multipart/form-data")
	part := &multipart.FileHeader{
		Filename: filepath,
		Header:   header,
		Size:     10 << 20, // 10MB
	}

	return file, part
}

func testFilePath(t *testing.T, filename string) string {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	strs := strings.Split(dir, "api/internal")
	if len(strs) == 0 {
		t.Fatal("test: invalid file path")
	}
	return filepath.Join(strs[0], "/api/tmp", filename)
}
