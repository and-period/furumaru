package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/stretchr/testify/assert"
)

func TestImageSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		size   common.ImageSize
		expect ImageSize
	}{
		{
			name:   "success to small",
			size:   common.ImageSizeSmall,
			expect: ImageSizeSmall,
		},
		{
			name:   "success to medium",
			size:   common.ImageSizeMedium,
			expect: ImageSizeMedium,
		},
		{
			name:   "success to large",
			size:   common.ImageSizeLarge,
			expect: ImageSizeLarge,
		},
		{
			name:   "success to unknown",
			size:   common.ImageSizeUnknown,
			expect: ImageSizeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewImageSize(tt.size))
		})
	}
}

func TestImageSize_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		size   ImageSize
		expect int32
	}{
		{
			name:   "success",
			size:   ImageSizeSmall,
			expect: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.size.Response())
		})
	}
}

func TestImage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		image  *common.Image
		expect *Image
	}{
		{
			name: "success",
			image: &common.Image{
				URL:  "http://example.com/media/image.png",
				Size: common.ImageSizeSmall,
			},
			expect: &Image{
				Image: response.Image{
					URL:  "http://example.com/media/image.png",
					Size: int32(ImageSizeSmall),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewImage(tt.image))
		})
	}
}

func TestImage_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		image  *Image
		expect *response.Image
	}{
		{
			name: "success",
			image: &Image{
				Image: response.Image{
					URL:  "http://example.com/media/image.png",
					Size: int32(ImageSizeSmall),
				},
			},
			expect: &response.Image{
				URL:  "http://example.com/media/image.png",
				Size: int32(ImageSizeSmall),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.image.Response())
		})
	}
}

func TestImages(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		images common.Images
		expect Images
	}{
		{
			name: "success",
			images: common.Images{
				{
					URL:  "http://example.com/media/image.png",
					Size: common.ImageSizeSmall,
				},
			},
			expect: Images{
				{
					Image: response.Image{
						URL:  "http://example.com/media/image.png",
						Size: int32(ImageSizeSmall),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewImages(tt.images))
		})
	}
}

func TestImages_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		images Images
		expect []*response.Image
	}{
		{
			name: "success",
			images: Images{
				{
					Image: response.Image{
						URL:  "http://example.com/media/image.png",
						Size: int32(ImageSizeSmall),
					},
				},
			},
			expect: []*response.Image{
				{
					URL:  "http://example.com/media/image.png",
					Size: int32(ImageSizeSmall),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.images.Response())
		})
	}
}
