package media

import (
	"io"
	"mime/multipart"
)

type GenerateFileInput struct {
	File   io.Reader             `validate:"required"`
	Header *multipart.FileHeader `validate:"required"`
}

type UploadFileInput struct {
	URL string `validate:"required,url"`
}

type ResizeFileInput struct {
	TargetID string   `validate:"required"`
	URLs     []string `validate:"min=1,dive,required,url"`
}

type GetBroadcastByScheduleIDInput struct {
	ScheduleID string `validate:"required"`
}

type CreateBroadcastInput struct {
	ScheduleID string `validate:"required"`
}
