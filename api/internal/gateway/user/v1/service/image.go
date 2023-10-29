package service

import (
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
)

// ImageSize - 画像サイズ
type ImageSize int32

const (
	ImageSizeUnknown ImageSize = 0
	ImageSizeSmall   ImageSize = 1 // 画像サイズ: 小
	ImageSizeMedium  ImageSize = 2 // 画像サイズ: 中
	ImageSizeLarge   ImageSize = 3 // 画像サイズ: 大
)

type Image struct {
	response.Image
}

type Images []*Image

func NewImageSize(size common.ImageSize) ImageSize {
	switch size {
	case common.ImageSizeSmall:
		return ImageSizeSmall
	case common.ImageSizeMedium:
		return ImageSizeMedium
	case common.ImageSizeLarge:
		return ImageSizeLarge
	default:
		return ImageSizeUnknown
	}
}

func (s ImageSize) Response() int32 {
	return int32(s)
}

func NewImage(image *common.Image) *Image {
	return &Image{
		Image: response.Image{
			URL:  image.URL,
			Size: NewImageSize(image.Size).Response(),
		},
	}
}

func (i *Image) Response() *response.Image {
	return &i.Image
}

func NewImages(images common.Images) Images {
	res := make(Images, len(images))
	for i := range images {
		res[i] = NewImage(images[i])
	}
	return res
}

func (is Images) Response() []*response.Image {
	res := make([]*response.Image, len(is))
	for i := range is {
		res[i] = is[i].Response()
	}
	return res
}
