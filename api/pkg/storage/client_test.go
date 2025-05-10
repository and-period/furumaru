package storage

import (
	"net/url"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBucket(t *testing.T) {
	t.Parallel()
	cfg, err := config.LoadDefaultConfig(t.Context())
	require.NoError(t, err)
	bucket := NewBucket(cfg, &Params{},
		WithMaxRetries(1),
		WithInterval(time.Millisecond),
		WithLogger(zap.NewNop()),
	)
	assert.NotNil(t, bucket)
}

func TestGenerateObjectURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		bucket string
		path   string
		expect string
		hasErr bool
	}{
		{
			name:   "success",
			bucket: "bucket",
			path:   "dir/path.png",
			expect: "https://bucket.s3.ap-northeast-1.amazonaws.com/dir/path.png",
			hasErr: false,
		},
		{
			name:   "failed to get host",
			bucket: " ",
			path:   "dir/path.png",
			expect: "",
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &bucket{
				s3:     &s3.Client{},
				name:   aws.String(tt.bucket),
				region: "ap-northeast-1",
				logger: &zap.Logger{},
			}
			actual, err := client.GenerateObjectURL(tt.path)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestGenerateS3URI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		bucket string
		path   string
		expect string
	}{
		{
			name:   "success",
			bucket: "bucket",
			path:   "dir/path.png",
			expect: "s3://bucket/dir/path.png",
		},
		{
			name:   "trim duplicate slash",
			bucket: "bucket",
			path:   "/dir/path.png",
			expect: "s3://bucket/dir/path.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &bucket{
				s3:     &s3.Client{},
				name:   aws.String(tt.bucket),
				region: "ap-northeast-1",
				logger: &zap.Logger{},
			}
			actual := client.GenerateS3URI(tt.path)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestReplaceURLToS3URI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		bucket string
		rawURL string
		expect string
		hasErr bool
	}{
		{
			name:   "success",
			bucket: "bucket",
			rawURL: "https://bucket.s3.ap-northeast-1.amazonaws.com/dir/path.png",
			expect: "s3://bucket/dir/path.png",
			hasErr: false,
		},
		{
			name:   "empty url",
			bucket: "bucket",
			rawURL: "",
			expect: "",
			hasErr: false,
		},
		{
			name:   "failed to parse url",
			bucket: "bucket",
			rawURL: "https:// ",
			expect: "",
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &bucket{
				s3:     &s3.Client{},
				name:   aws.String(tt.bucket),
				region: "ap-northeast-1",
				logger: &zap.Logger{},
			}
			actual, err := client.ReplaceURLToS3URI(tt.rawURL)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestGetHost(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		bucket string
		expect *url.URL
		hasErr bool
	}{
		{
			name:   "success",
			bucket: "bucket",
			expect: &url.URL{
				Scheme: "https",
				Host:   "bucket.s3.ap-northeast-1.amazonaws.com",
			},
			hasErr: false,
		},
		{
			name:   "failed to get host",
			bucket: " ",
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &bucket{
				s3:     &s3.Client{},
				name:   aws.String(tt.bucket),
				region: "ap-northeast-1",
				logger: &zap.Logger{},
			}
			actual, err := client.GetHost()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
