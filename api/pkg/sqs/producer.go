package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.uber.org/zap"
)

type producer struct {
	sqs          *sqs.Client
	logger       *zap.Logger
	dryRun       bool
	queueURL     *string
	delaySeconds int32
}

func NewProducer(cfg aws.Config, params *Params, opts ...Option) Producer {
	dopts := &options{
		dryRun:       false,
		maxRetries:   retry.DefaultMaxAttempts,
		interval:     retry.DefaultMaxBackoff,
		logger:       zap.NewNop(),
		delaySeconds: 0,
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &producer{
		sqs:          cli,
		logger:       dopts.logger,
		dryRun:       dopts.dryRun,
		queueURL:     aws.String(params.QueueURL),
		delaySeconds: dopts.delaySeconds,
	}
}

func (p *producer) SendMessage(ctx context.Context, b []byte) (string, error) {
	if p.dryRun {
		p.logger.Info("Receive SendMessage", zap.Binary("body", b))
		return "", nil
	}
	in := &sqs.SendMessageInput{
		QueueUrl:     p.queueURL,
		DelaySeconds: p.delaySeconds,
		MessageBody:  aws.String(string(b)),
	}
	out, err := p.sqs.SendMessage(ctx, in)
	if err != nil {
		return "", err
	}
	return aws.ToString(out.MessageId), nil
}
