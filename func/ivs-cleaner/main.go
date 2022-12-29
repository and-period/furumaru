package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	itypes "github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"go.uber.org/zap"
)

type client struct {
	logger   *zap.Logger
	dynamodb *dynamodb.Client
	ivs      *ivs.Client
}

var app = &client{}

func main() {
	lambda.Start(Handler)
}

func init() {
	ctx := context.Background()
	// Loggerの設定
	var err error
	app.logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to new logger: err=%s", err.Error())
		return
	}
	// AWS SDKの設定
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		app.logger.Error("Failed to load aws config", zap.Error(err))
		return
	}
	app.dynamodb = dynamodb.NewFromConfig(cfg)
	app.ivs = ivs.NewFromConfig(cfg)
}

func Handler(ctx context.Context, e events.DynamoDBEvent) error {
	for _, record := range e.Records {
		if events.DynamoDBOperationType(record.EventName) != events.DynamoDBOperationTypeRemove {
			app.logger.Debug("Operation Type is not remove", zap.Any("record", record))
			return nil
		}
		r, err := unmarshal(record.Change.OldImage)
		if err != nil {
			app.logger.Error("Failed to parse record", zap.Any("record", record), zap.Error(err))
			return nil
		}
		in := &ivs.DeleteChannelInput{
			Arn: &r.ChannelArn,
		}
		var e *itypes.ResourceNotFoundException
		_, err = app.ivs.DeleteChannel(ctx, in)
		if errors.As(err, &e) {
			app.logger.Info("Amazon IVS Channel has already deleted", zap.Any("record", record))
			return nil
		}
		if err != nil {
			app.logger.Error("Failed to delete Amazon IVS Channel", zap.Any("record", record), zap.Error(err))
			return err
		}
		app.logger.Info("Success to delete Amazon IVS Channel", zap.Any("record", record))
	}
	return nil
}

// rehearsal - リハーサルスケジュール
type rehearsal struct {
	LiveID         string    `dynamodbav:"live_id"`                       // ライブ配信ID
	ScheduleID     string    `dynamodbav:"schedule_id,omitempty"`         // 開催スケジュールID
	ChannelName    string    `dyanmodbav:"channel_name,omitempty"`        // IVS チャンネル名
	ChannelArn     string    `dynamodbav:"channel_arn,omitempty"`         // IVS チャンネルARN
	IngestEndpoint string    `dynamodbav:"ingest_endpoint,omitempty"`     // IVS 配信取り込みエンドポイント
	StreamKey      string    `dynamodbav:"stream_key,omitempty"`          // IVS 配信用ストリームキー
	PlaybackURL    string    `dynamodbav:"playback_url,omitempty"`        // IVS 再生用URL
	StartAt        time.Time `dynamodbav:"start_at,omitempty"`            // 配信開始日時
	ExpiresAt      time.Time `dynamodbav:"expires_at,unixtime,omitempty"` // リハーサルモード実施有効期限
	CreatedAt      time.Time `dynamodbav:"created_at,omitempty"`          // 登録日時
	UpdatedAt      time.Time `dynamodbav:"updated_at,omitempty"`          // 更新日時
}

func unmarshal(image map[string]events.DynamoDBAttributeValue) (*rehearsal, error) {
	attrs := make(map[string]dtypes.AttributeValue, len(image))
	for k, v := range image {
		attr, err := unmarshalValue(v)
		if err != nil {
			app.logger.Error("Failed to marshal", zap.String("key", k), zap.Any("value", v), zap.Error(err))
			return nil, fmt.Errorf("failed to marshal: %w", err)
		}
		attrs[k] = attr
	}
	res := &rehearsal{}
	if err := attributevalue.UnmarshalMap(attrs, res); err != nil {
		app.logger.Error("Failed to unmarshal map", zap.Any("attributes", attrs), zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal map: %w", err)
	}
	return res, nil
}

func unmarshalValue(av events.DynamoDBAttributeValue) (dtypes.AttributeValue, error) {
	var v interface{}
	switch av.DataType() {
	case events.DataTypeString:
		v = av.String()
	case events.DataTypeNumber:
		var err error
		v, err = av.Integer()
		if err != nil {
			return nil, err
		}
	}
	return attributevalue.Marshal(v)
}
