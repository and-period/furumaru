package admin

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/pkg/batch"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func (a *app) injectAWS(ctx context.Context, p *params) error {
	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(a.AWSRegion))
	if err != nil {
		return fmt.Errorf("cmd: failed to load aws config: %w", err)
	}
	p.aws = awscfg

	// AWS Secrets Managerの設定
	p.secret = secret.NewClient(awscfg)
	if err := a.getSecret(ctx, p); err != nil {
		return fmt.Errorf("cmd: failed to get secret: %w", err)
	}

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: a.S3Bucket,
	}
	p.storage = storage.NewBucket(awscfg, storageParams)
	tmpStorageParams := &storage.Params{
		Bucket: a.S3TmpBucket,
	}
	p.tmpStorage = storage.NewBucket(awscfg, tmpStorageParams)

	// Amazon SQSの設定
	messengerSQSParams := &sqs.Params{
		QueueURL: a.SQSMessengerQueueURL,
	}
	p.messengerQueue = sqs.NewProducer(awscfg, messengerSQSParams, sqs.WithDryRun(a.SQSMockEnabled))
	mediaSQSParams := &sqs.Params{
		QueueURL: a.SQSMediaQueueURL,
	}
	p.mediaQueue = sqs.NewProducer(awscfg, mediaSQSParams, sqs.WithDryRun(a.SQSMockEnabled))

	// Amazon DynamoDBの設定
	dbParams := &dynamodb.Params{
		TablePrefix: "furumaru",
		TableSuffix: a.Environment,
	}
	p.cache = dynamodb.NewClient(awscfg, dbParams)

	// AWS Batchの設定
	p.batch = batch.NewClient(awscfg)

	// AWS MediaLiveの設定
	p.medialive = medialive.NewMediaLive(awscfg)

	return nil
}
