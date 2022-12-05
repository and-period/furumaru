package resizer

import (
	"context"
	"sync"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"golang.org/x/sync/errgroup"
)

func (r *resizer) productMedia(ctx context.Context, payload *entity.ResizerPayload) error {
	if len(payload.URLs) == 0 {
		return errRequiredMediaURL
	}
	var mu sync.Mutex
	media := make([]*store.UpdateProductMediaImage, 0, len(payload.URLs))
	eg, ectx := errgroup.WithContext(ctx)
	for _, url := range payload.URLs {
		url := url
		eg.Go(func() error {
			file, err := r.storage.Download(ectx, url)
			if err != nil {
				return err
			}
			resizedImages, err := r.resizeImages(url, file)
			if err != nil || len(resizedImages) == 0 {
				return err
			}
			images, err := r.uploadImages(ectx, url, resizedImages)
			if err != nil {
				return err
			}
			mu.Lock()
			media = append(media, &store.UpdateProductMediaImage{OriginURL: url, Images: images})
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return exception.InternalError(err)
	}
	in := &store.UpdateProductMediaInput{
		ProductID: payload.TargetID,
		Images:    media,
	}
	updateFn := func() error {
		err := r.store.UpdateProductMedia(ctx, in)
		return exception.InternalError(err)
	}
	return r.notify(ctx, payload, updateFn)
}
