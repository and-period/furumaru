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

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegulation_Validate(t *testing.T) {
	t.Parallel()
	reg := &Regulation{
		MaxSize: 10 << 20, // 10MB
		Formats: set.New("image/png", "image/jpeg"),
		dir:     "image/png",
	}
	type args struct {
		fileType string
		fileSize int64
	}
	tests := []struct {
		name string
		args
		expect error
	}{
		{
			name: "success",
			args: args{
				fileType: "image/png",
				fileSize: 2 << 20, // 2MB
			},
			expect: nil,
		},
		{
			name: "invalid size",
			args: args{
				fileType: "image/png",
				fileSize: 20 << 20, // 20MB
			},
			expect: ErrTooLargeFileSize,
		},
		{
			name: "invalid format",
			args: args{
				fileType: "video/mp4",
				fileSize: 2 << 20, // 2MB
			},
			expect: ErrInvalidFileFormat,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := reg.Validate(tt.fileType, tt.fileSize)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}

func TestRegulation_GetObjectKey(t *testing.T) {
	t.Parallel()
	reg := &Regulation{
		MaxSize: 0,
		Formats: set.New[string]("image/jpeg", "image/png", "image/webp", "video/mp4"),
		dir:     "test",
	}
	tests := []struct {
		name        string
		regulation  *Regulation
		contentType string
		args        []interface{}
		expect      string
		expectErr   error
	}{
		{
			name:        "success image/jpeg",
			regulation:  reg,
			contentType: "image/jpeg",
			args:        []interface{}{},
			expect:      "test/[a-zA-Z0-9]+.jpg",
			expectErr:   nil,
		},
		{
			name:        "success image/png",
			regulation:  reg,
			contentType: "image/png",
			expect:      "test/[a-zA-Z0-9]+.png",
			expectErr:   nil,
		},
		{
			name:        "success with params",
			regulation:  BroadcastArchiveRegulation,
			contentType: "video/mp4",
			args:        []interface{}{"broadcast-id"},
			expect:      "schedules/archives/broadcast-id/mp4/[a-zA-Z0-9]+",
			expectErr:   nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.regulation.GetObjectKey(tt.contentType, tt.args...)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Regexp(t, regexp.MustCompile(tt.expect), actual, actual)
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
