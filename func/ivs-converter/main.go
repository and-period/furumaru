package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type client struct {
	logger  *zap.Logger
	s3      *s3.Client
	convert *mediaconvert.Client
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
	app.s3 = s3.NewFromConfig(cfg)
	opts := mediaconvert.WithEndpointResolver(mediaconvert.EndpointResolverFromURL(os.Getenv("MEDIA_CONVERT_ENDPOINT")))
	app.convert = mediaconvert.NewFromConfig(cfg, opts)
}

func Handler(ctx context.Context, e events.CloudWatchEvent) error {
	if e.Source != "aws.ivs" && e.DetailType != "IVS Recording State Change" {
		app.logger.Debug("Event type is not ivs recording state change", zap.Any("event", e))
		return nil
	}
	img, err := unmarshal(e.Detail)
	if err != nil {
		app.logger.Error("Failed to unmarshal detail", zap.Any("event", e), zap.Error(err))
		return nil
	}
	if img.RecordingStatus != "Recording End" {
		app.logger.Debug("Recording status is not ended", zap.Any("event", e), zap.Any("image", e))
		return nil
	}
	conf, err := getIVSConfiguration(ctx, img)
	if err != nil {
		app.logger.Error("Failed to get ivs configuration", zap.Any("event", e), zap.Any("image", img), zap.Error(err))
		return err
	}
	if conf.RecordingStatus != "RECORDING_ENDED" {
		app.logger.Debug("Recording status is not ended", zap.Any("image", img))
		return nil
	}
	if err := createMediaConvertJob(ctx, img, conf); err != nil {
		app.logger.Error("Failed to create media convert job", zap.Any("event", e), zap.Any("image", img), zap.Error(err))
		return err
	}
	app.logger.Info("Success to create media convert job", zap.Any("event", e), zap.Any("image", img))
	return nil
}

// image - EventBridgeからのイベント詳細
type image struct {
	ChannelName           string `json:"channel_name"`
	StreamID              string `json:"stream_id"`
	RecordingStatus       string `json:"recording_status"`
	RecordingStatusReason string `json:"recording_status_reason"`
	RecordingS3BucketName string `json:"recording_s3_bucket_name"`
	RecordingS3KeyPrefix  string `json:"recording_s3_key_prefix"`
	RecordingDurationMs   int64  `json:"recording_duration_ms"`
	RecordingSessionID    string `json:"recording_session_id"`
}

// unmarshal - EventBridgeからのイベント処理
func unmarshal(detail json.RawMessage) (*image, error) {
	buf, err := detail.MarshalJSON()
	if err != nil {
		app.logger.Error("Failed to marshal detail", zap.Error(err))
		return nil, err
	}
	img := &image{}
	if err := json.Unmarshal(buf, img); err != nil {
		app.logger.Error("Failed to unmarshal detail", zap.Error(err))
		return nil, err
	}
	return img, nil
}

type configuration struct {
	Version            string    `json:"version"`
	RecordingStartedAt time.Time `json:"recording_started_at"`
	RecordingEndedAt   time.Time `json:"recording_ended_at"`
	ChannelARN         string    `json:"channel_arn"`
	RecordingStatus    string    `json:"recording_status"`
	Media              Media     `json:"media"`
}

type Media struct {
	HLS        HLS        `json:"hls"`
	Thumbnails Thumbnails `json:"thumbnails"`
}

type HLS struct {
	DurationMS int64       `json:"duration_ms"`
	Path       string      `json:"path"`
	Playlist   string      `json:"playlist"`
	Renditions []Rendition `json:"renditions"`
}

type Rendition struct {
	Path             string `json:"path"`
	Playlist         string `json:"playlist"`
	ResolutionWidth  int64  `json:"resolution_width"`
	ResolutionHeight int64  `json:"resolution_height"`
}

type Thumbnails struct {
	Path string `json:"path"`
}

// getIVSConfiguration - Amaozon IVSの情報をS3から取得
func getIVSConfiguration(ctx context.Context, img *image) (*configuration, error) {
	const filename = "events/recording-ended.json"
	in := &s3.GetObjectInput{
		Bucket: aws.String(img.RecordingS3BucketName),
		Key:    aws.String(strings.Join([]string{img.RecordingS3KeyPrefix, filename}, "/")),
	}
	object, err := app.s3.GetObject(ctx, in)
	if err != nil {
		app.logger.Error("Failed to get s3 object", zap.Any("params", in), zap.Error(err))
		return nil, err
	}
	defer object.Body.Close()
	cbuf, err := io.ReadAll(object.Body)
	if err != nil {
		app.logger.Error("Failed to read template file", zap.Error(err))
		return nil, err
	}
	conf := &configuration{}
	if err := json.Unmarshal(cbuf, conf); err != nil {
		app.logger.Error("Failed to unmarshal configuration", zap.Error(err))
		return nil, err
	}
	return conf, nil
}

// createMediaConvertJob - Media ConvertのJobを作成
func createMediaConvertJob(ctx context.Context, img *image, conf *configuration) error {
	in := &mediaconvert.CreateJobInput{
		Role:        aws.String(os.Getenv("MEDIA_CONVERT_ROLE")),
		Settings:    newMediaconvertJobSettings(img, conf),
		JobTemplate: aws.String(os.Getenv("MEDIA_CONVERT_JOB_TEMPLATE")),
	}
	out, err := app.convert.CreateJob(ctx, in)
	if err != nil {
		app.logger.Error("Failed to create job", zap.Any("params", in), zap.Error(err))
		return err
	}
	app.logger.Debug("Success create job", zap.Any("output", out))
	return nil
}

func newMediaconvertJobSettings(img *image, conf *configuration) *types.JobSettings {
	base := filepath.Join(img.RecordingS3BucketName, img.RecordingS3KeyPrefix)
	src := fmt.Sprintf("s3://%s", filepath.Join(base, conf.Media.HLS.Path, conf.Media.HLS.Playlist))
	dst := fmt.Sprintf("s3://%s", filepath.Join(base, "media"))

	return &types.JobSettings{
		Inputs: []types.Input{{
			FileInput:      aws.String(src),
			TimecodeSource: types.InputTimecodeSourceZerobased,
		}},
		OutputGroups: []types.OutputGroup{{
			OutputGroupSettings: &types.OutputGroupSettings{
				Type: types.OutputGroupTypeFileGroupSettings,
				FileGroupSettings: &types.FileGroupSettings{
					Destination: aws.String(dst),
				},
			},
		}},
	}
}
