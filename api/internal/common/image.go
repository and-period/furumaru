package common

import "encoding/json"

// ImageSize - 画像サイズ種別
type ImageSize int32

const (
	ImageSizeUnknown ImageSize = 0
	ImageSizeSmall   ImageSize = 1 // 画像サイズ: 小
	ImageSizeMedium  ImageSize = 2 // 画像サイズ: 中
	ImageSizeLarge   ImageSize = 3 // 画像サイズ: 大
)

var ImageSizes = []ImageSize{
	ImageSizeSmall,
	ImageSizeMedium,
	ImageSizeLarge,
}

// Image - メディア画像情報
type Image struct {
	URL  string    `json:"url"`  // 参照先URL
	Size ImageSize `json:"size"` // 画像サイズ
}

type Images []*Image

func NewImagesFromBytes(b []byte) (Images, error) {
	if b == nil {
		return Images{}, nil
	}
	var images Images
	return images, json.Unmarshal(b, &images)
}

func (is Images) Marshal() ([]byte, error) {
	if len(is) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(is)
}
