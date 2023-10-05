package resizer

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
)

func (r *resizer) scheduleThumbnail(ctx context.Context, payload *entity.ResizerPayload) error {
	if len(payload.URLs) == 0 || payload.URLs[0] == "" {
		return errRequiredMediaURL
	}
	url := payload.URLs[0]
	file, err := r.storage.Download(ctx, payload.URLs[0])
	if err != nil {
		return err
	}
	resizedImages, err := r.resizeImages(url, file)
	if err != nil {
		return err
	}
	images, err := r.uploadImages(ctx, url, resizedImages)
	if err != nil {
		return err
	}
	in := &store.UpdateScheduleThumbnailsInput{
		ScheduleID: payload.TargetID,
		Thumbnails: images,
	}
	updateFn := func() error {
		return r.store.UpdateScheduleThumbnails(ctx, in)
	}
	return r.notify(ctx, payload, updateFn)
}
