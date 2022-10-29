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
