package entity

import (
	"regexp"
	"testing"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
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

func TestRegulation_ShouldConvert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		regulation  *Regulation
		contentType string
		expect      bool
	}{
		{
			name: "none",
			regulation: &Regulation{
				Convert: ConvertTypeNone,
			},
			contentType: "image/jpeg",
			expect:      false,
		},
		{
			name: "jpeg to png should convert",
			regulation: &Regulation{
				Convert: ConvertTypeJPEGToPNG,
			},
			contentType: "image/jpeg",
			expect:      true,
		},
		{
			name: "jpeg to png should not convert",
			regulation: &Regulation{
				Convert: ConvertTypeJPEGToPNG,
			},
			contentType: "image/png",
			expect:      false,
		},
		{
			name: "unknown",
			regulation: &Regulation{
				Convert: -1,
			},
			contentType: "image/jpeg",
			expect:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.regulation.ShouldConvert(tt.contentType)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
