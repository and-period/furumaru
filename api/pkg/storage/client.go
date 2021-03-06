//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

var (
	ErrInvalidURL = errors.New("s3: invalid s3 object url")
	ErrNotFound   = errors.New("s3: not found object")
)

const domain = "%s.s3.amazonaws.com"

type Bucket interface {
	// オブジェクトURLの生成
	GenerateObjectURL(path string) string
	// S3 Bucketからオブジェクトを取得
	Download(ctx context.Context, url string) (io.Reader, error)
	// S3 Bucketからオブジェクトを取得とByte型へ変換
	DownloadAndReadAll(ctx context.Context, url string) ([]byte, error)
	// S3 Bucketへオブジェクトをアップロード
	Upload(ctx context.Context, path string, body io.Reader) (string, error)
}

type Params struct {
	Bucket string
}

type bucket struct {
	s3     *s3.Client
	name   *string
	logger *zap.Logger
}

type options struct {
	maxRetries int
	interval   time.Duration
	logger     *zap.Logger
}

type Option func(*options)

func WithMaxRetries(maxRetries int) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.interval = interval
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewBucket(cfg aws.Config, params *Params, opts ...Option) Bucket {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
		logger:     zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &bucket{
		s3:     cli,
		name:   aws.String(params.Bucket),
		logger: dopts.logger,
	}
}

func (b *bucket) GenerateObjectURL(path string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf(domain, aws.ToString(b.name)),
		Path:   path,
	}
	return u.String()
}

func (b *bucket) Download(ctx context.Context, url string) (io.Reader, error) {
	path, err := b.generateKeyFromObjectURL(url)
	if err != nil {
		return nil, err
	}
	in := &s3.GetObjectInput{
		Bucket: b.name,
		Key:    aws.String(path),
	}
	out, err := b.s3.GetObject(ctx, in)
	var bne *types.NotFound
	if errors.As(err, &bne) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (b *bucket) DownloadAndReadAll(ctx context.Context, url string) ([]byte, error) {
	reader, err := b.Download(ctx, url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(reader)
}

func (b *bucket) Upload(ctx context.Context, path string, body io.Reader) (string, error) {
	in := &s3.PutObjectInput{
		Bucket: b.name,
		Key:    aws.String(path),
		Body:   body,
	}
	_, err := b.s3.PutObject(ctx, in)
	if err != nil {
		return "", err
	}
	return b.GenerateObjectURL(path), nil
}

func (b *bucket) generateKeyFromObjectURL(objectURL string) (string, error) {
	u, err := url.Parse(objectURL)
	if err != nil {
		return "", ErrInvalidURL
	}
	return u.Path, nil
}
