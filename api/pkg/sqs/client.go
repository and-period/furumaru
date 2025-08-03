//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package sqs

import (
	"context"
	"time"
)

type Producer interface {
	// メッセージの送信
	SendMessage(ctx context.Context, b []byte) (string, error)
}

type Params struct {
	QueueURL string
}

type options struct {
	dryRun       bool
	maxRetries   int
	interval     time.Duration
	delaySeconds int32 // for producer
	timeout      int32 // for consumer
}

type Option func(*options)

func WithDryRun(dryRun bool) Option {
	return func(opts *options) {
		opts.dryRun = dryRun
	}
}

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

func WithDelaySeconds(delaySeconds int32) Option {
	return func(opts *options) {
		opts.delaySeconds = delaySeconds
	}
}

func WithTimeout(timeout int32) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}
