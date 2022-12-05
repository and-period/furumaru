package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImagesFromBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		buf    []byte
		expect Images
		hasErr bool
	}{
		{
			name: "success",
			buf:  []byte(`[{"url":"http://example.com/media.png","size":1}]`),
			expect: Images{
				{
					Size: ImageSizeSmall,
					URL:  "http://example.com/media.png",
				},
			},
			hasErr: false,
		},
		{
			name:   "success empty",
			buf:    nil,
			expect: Images{},
			hasErr: false,
		},
		{
			name:   "failed to parse error",
			buf:    []byte(`[`),
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewImagesFromBytes(tt.buf)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestImages_Marshal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		images Images
		expect []byte
		hasErr bool
	}{
		{
			name: "success",
			images: Images{
				{
					Size: ImageSizeSmall,
					URL:  "http://example.com/media.png",
				},
			},
			expect: []byte(`[{"url":"http://example.com/media.png","size":1}]`),
			hasErr: false,
		},
		{
			name:   "success empty",
			images: Images{},
			expect: []byte{},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.images.Marshal()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
