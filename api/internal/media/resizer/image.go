package resizer

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/shopspring/decimal"
	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
)

func (r *resizer) uploadImages(
	ctx context.Context,
	originURL string,
	images map[common.ImageSize]io.Reader,
) (common.Images, error) {
	res := make(common.Images, 0, len(images))

	var mu sync.Mutex
	eg, ectx := errgroup.WithContext(ctx)
	for size, image := range images {
		size, image := size, image
		eg.Go(func() error {
			path, err := r.generateFilePath(originURL, size)
			if err != nil {
				return err
			}
			md := map[string]string{
				"Cache-Control": "max-age=" + r.ttl.String(),
			}
			url, err := r.storage.Upload(ectx, path, image, md)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			res = append(res, &common.Image{URL: url, Size: size})
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *resizer) generateFilePath(originURL string, size common.ImageSize) (string, error) {
	u, err := url.Parse(originURL)
	if err != nil {
		return "", err
	}
	path := strings.TrimPrefix(u.Path, "/")

	var suffix string
	switch size {
	case common.ImageSizeSmall:
		suffix = "240"
	case common.ImageSizeMedium:
		suffix = "675"
	case common.ImageSizeLarge:
		suffix = "900"
	default:
		return "", errUnsupportedImageSize
	}

	extension := filepath.Ext(originURL)
	keys := strings.Split(path, extension)
	filepath := fmt.Sprintf("%s_%s%s", keys[0], suffix, extension)
	return filepath, nil
}

func (r *resizer) resizeImages(originURL string, src io.Reader) (map[common.ImageSize]io.Reader, error) {
	targets := set.New(".png", ".jpg", ".jpeg")
	extension := filepath.Ext(originURL)
	if !targets.Contains(extension) || src == nil {
		return map[common.ImageSize]io.Reader{}, nil
	}

	images := make(map[common.ImageSize]io.Reader, len(common.ImageSizes))

	for _, size := range common.ImageSizes {
		buf := &bytes.Buffer{}
		teeReader := io.TeeReader(src, buf)
		src = buf

		img, err := r.resizeImage(teeReader, size)
		if err != nil {
			return nil, err
		}
		images[size] = img
	}

	return images, nil
}

func (r *resizer) resizeImage(src io.Reader, size common.ImageSize) (io.Reader, error) {
	const (
		imageSizeSmall  = 240
		imageSizeMedium = 675
		imageSizeLarge  = 900
	)

	img, data, err := image.Decode(src)
	if err != nil {
		return nil, err
	}
	rect := img.Bounds()

	var limit int
	switch size {
	case common.ImageSizeSmall:
		limit = imageSizeSmall
	case common.ImageSizeMedium:
		limit = imageSizeMedium
	case common.ImageSizeLarge:
		limit = imageSizeLarge
	default:
		return nil, errUnsupportedImageSize
	}

	width, height := func(rect image.Rectangle, limit int) (int, int) {
		width, height := rect.Dx(), rect.Dy()
		if width == height {
			return limit, limit
		}

		dlimit := decimal.New(int64(limit), 0)
		dwidth := decimal.NewFromInt(int64(width))
		dheight := decimal.NewFromInt(int64(height))
		if width > height {
			drate := dlimit.Div(dwidth)
			dwidth = dlimit
			dheight = dheight.Mul(drate)
		} else {
			drate := dlimit.Div(dheight)
			dwidth = dwidth.Mul(drate)
			dheight = dlimit
		}
		return int(dwidth.IntPart()), int(dheight.IntPart())
	}(rect, limit)

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, rect, draw.Over, nil)

	buf := &bytes.Buffer{}
	switch data {
	case "png":
		err = png.Encode(buf, dst)
	case "jpeg":
		err = jpeg.Encode(buf, dst, &jpeg.Options{Quality: 100})
	default:
		return nil, errUnsupportedImageFormat
	}
	return buf, err
}
