package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/database/tidb"
	"github.com/and-period/furumaru/api/internal/media/entity"
	transcribe "github.com/and-period/furumaru/api/pkg/aws/transcribe"
	translate "github.com/and-period/furumaru/api/pkg/aws/translate"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awstranscribe "github.com/aws/aws-sdk-go-v2/service/transcribe"
	transcribetype "github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	awstranslate "github.com/aws/aws-sdk-go-v2/service/translate"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	appName          string
	environment      string
	logLevel         string
	awsRegion        string
	assetsDomain     string
	dbSecretName     string
	dbDatabase       string
	dbTimeZone       string
	sentrySecretName string
	s3Bucket         string
	broadcastID      string

	tidbHost     string
	tidbPort     string
	tidbUsername string
	tidbPassword string
	sentryDSN    string

	dbLocation, _ = time.LoadLocation(dbTimeZone)
)

type app struct {
	logger     *zap.Logger
	db         *database.Database
	s3         storage.Bucket
	transcribe transcribe.Client
	translate  translate.Client
}

//nolint:forbidigo,gocritic
func main() {
	startedAt := jst.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	app, err := setup(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup: %v\n", err)
		os.Exit(1)
	}

	if err := app.run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %v\n", err)
		os.Exit(1)
	}

	endAt := jst.Now()

	const format = "2006-01-02 15:04:05"
	fmt.Printf("Done: %s -> %s (%s)\n",
		jst.Format(startedAt, format), jst.Format(endAt, format),
		endAt.Sub(startedAt).Truncate(time.Second).String(),
	)
}

func setup(ctx context.Context) (*app, error) {
	flag.StringVar(&appName, "app-name", "media-update-archive", "application name")
	flag.StringVar(&environment, "environment", "", "environment")
	flag.StringVar(&logLevel, "log-level", "debug", "log level")
	flag.StringVar(&awsRegion, "aws-region", "ap-northeast-1", "AWS region")
	flag.StringVar(&dbSecretName, "db-secret-name", "", "AWS Secrets Manager secret name for TiDB")
	flag.StringVar(&dbDatabase, "db-database", "media", "TiDB database name")
	flag.StringVar(&dbTimeZone, "db-timezone", "Asia/Tokyo", "TiDB timezone")
	flag.StringVar(&sentrySecretName, "sentry-secret-name", "", "AWS Secrets Manager secret name for Sentry")
	flag.StringVar(&assetsDomain, "assets-domain", "", "assets domain")
	flag.StringVar(&s3Bucket, "s3-bucket", "", "target S3 bucket name")
	flag.StringVar(&broadcastID, "broadcast-id", "", "target broadcast id")
	flag.Parse()

	if dbSecretName == "" {
		return nil, fmt.Errorf("db-secret-name is required")
	}
	if sentrySecretName == "" {
		return nil, fmt.Errorf("sentry-secret-name is required")
	}
	if assetsDomain == "" {
		return nil, fmt.Errorf("assets-domain is required")
	}
	if s3Bucket == "" {
		return nil, fmt.Errorf("s3-bucket is required")
	}
	if broadcastID == "" {
		return nil, fmt.Errorf("broadcast-id is required")
	}

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(awsRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	// AWS Secrets Managerの設定
	secret := secret.NewClient(awscfg)

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		secrets, err := secret.Get(ectx, dbSecretName)
		if err != nil {
			return fmt.Errorf("failed to get db secret: %w", err)
		}
		tidbHost = secrets["host"]
		tidbPort = secrets["port"]
		tidbUsername = secrets["username"]
		tidbPassword = secrets["password"]
		return nil
	})
	eg.Go(func() error {
		secrets, err := secret.Get(ectx, sentrySecretName)
		if err != nil {
			return fmt.Errorf("failed to get sentry secret: %w", err)
		}
		sentryDSN = secrets["dsn"]
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get secrets: %w", err)
	}

	// Loggerの設定
	logger, err := log.NewSentryLogger(sentryDSN,
		log.WithLogLevel(logLevel),
		log.WithSentryServerName(appName),
		log.WithSentryEnvironment(environment),
		log.WithSentryLevel("error"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sentry logger: %w", err)
	}

	// TiDBの設定
	db, err := mysql.NewTiDBClient(&mysql.Params{
		Host:     tidbHost,
		Port:     tidbPort,
		Database: dbDatabase,
		Username: tidbUsername,
		Password: tidbPassword,
	}, mysql.WithNow(jst.Now), mysql.WithLocation(dbLocation))
	if err != nil {
		return nil, fmt.Errorf("failed to create tidb client: %w", err)
	}

	// Storageの設定
	storage := storage.NewBucket(awscfg, &storage.Params{
		Bucket: s3Bucket,
	}, storage.WithLogger(logger))

	app := &app{
		logger:     logger,
		db:         tidb.NewDatabase(db),
		s3:         storage,
		transcribe: transcribe.NewClient(awscfg, transcribe.WithLogger(logger)),
		translate:  translate.NewClient(awscfg, translate.WithLogger(logger)),
	}
	return app, nil
}

func (a *app) run(ctx context.Context) error {
	a.logger.Info("start", zap.String("broadcastId", broadcastID))

	broadcast, err := a.db.Broadcast.Get(ctx, broadcastID)
	if errors.Is(err, database.ErrNotFound) {
		a.logger.Warn("broadcast not found", zap.String("broadcastId", broadcastID))
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get broadcast: %w", err)
	}

	a.logger.Info("archive url", zap.String("broadcastId", broadcastID), zap.String("archiveUrl", broadcast.ArchiveURL))

	u, err := url.Parse(broadcast.ArchiveURL)
	if err != nil {
		return fmt.Errorf("failed to parse archive url: %w", err)
	}

	metadata, err := a.s3.GetMetadata(ctx, u.Path)
	if err != nil {
		return fmt.Errorf("failed to get metadata for archive: %w", err)
	}
	if metadata.ContentType != "video/mp4" {
		return fmt.Errorf("invalid content type: %s", metadata.ContentType)
	}

	a.logger.Info("start execute transcibe", zap.String("broadcastId", broadcastID))

	japaneseTextKey, err := a.executeTranscribe(ctx, broadcast)
	if err != nil {
		return fmt.Errorf("failed to execute transcribe: %w", err)
	}

	a.logger.Info("finished execute transcribe", zap.String("broadcastId", broadcastID), zap.String("japaneseTextKey", japaneseTextKey))
	a.logger.Info("start execute translate", zap.String("broadcastId", broadcastID))

	englishTextKey, err := a.executeTranslate(ctx, broadcast)
	if err != nil {
		return fmt.Errorf("failed to execute translate: %w", err)
	}

	a.logger.Info("finished execute translate", zap.String("broadcastId", broadcastID), zap.String("englishTextKey", englishTextKey))
	a.logger.Info("start upload fixed archive", zap.String("broadcastId", broadcastID))

	archiveKey, err := a.uploadFixedArchive(ctx, broadcast, japaneseTextKey, englishTextKey)
	if err != nil {
		return fmt.Errorf("failed to convert archive: %w", err)
	}

	a.logger.Info("finished upload fixed archive", zap.String("broadcastId", broadcastID))
	a.logger.Info("start update archive", zap.String("broadcastId", broadcastID))

	params := &database.UpdateBroadcastParams{
		UploadBroadcastArchiveParams: &database.UploadBroadcastArchiveParams{
			ArchiveURL:   fmt.Sprintf("https://%s/%s", assetsDomain, archiveKey),
			ArchiveFixed: true,
		},
	}
	if err := a.db.Broadcast.Update(ctx, broadcastID, params); err != nil {
		return fmt.Errorf("failed to update broadcast: %w", err)
	}

	a.logger.Info("updated archive", zap.String("broadcastId", broadcastID))
	return nil
}

// 動画から日本語テキストを抽出
func (a *app) executeTranscribe(ctx context.Context, broadcast *entity.Broadcast) (string, error) {
	u, err := url.Parse(broadcast.ArchiveURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse archive url: %w", err)
	}

	inputKey := strings.TrimPrefix(u.Path, "/")
	inputURI := a.s3.GenerateS3URI(inputKey)

	filename := strings.Split(filepath.Base(inputKey), ".")[0]
	outputDir := fmt.Sprintf(entity.BroadcastArchiveTextPath, broadcast.ScheduleID)
	outputKey := fmt.Sprintf("%s/%s-ja", outputDir, filename)
	outputFilename := fmt.Sprintf("%s/%s-ja.srt", outputDir, filename)

	jobName := fmt.Sprintf("%s-%s-%s", environment, broadcast.ScheduleID, filename)

	currentIn := &awstranscribe.GetTranscriptionJobInput{
		TranscriptionJobName: aws.String(jobName),
	}
	current, err := a.transcribe.GetTranscriptionJob(ctx, currentIn)
	if err == nil && current.TranscriptionJob.TranscriptionJobStatus == transcribetype.TranscriptionJobStatusCompleted {
		a.logger.Info("transcription job already completed", zap.String("broadcastId", broadcast.ID))
		return outputFilename, nil
	}

	in := &awstranscribe.StartTranscriptionJobInput{
		Media: &transcribetype.Media{
			MediaFileUri: aws.String(inputURI),
		},
		TranscriptionJobName: aws.String(jobName),
		LanguageCode:         transcribetype.LanguageCodeJaJp, // 日本語
		MediaFormat:          transcribetype.MediaFormatMp4,
		OutputBucketName:     aws.String(s3Bucket),
		OutputKey:            aws.String(outputKey),
		Subtitles: &transcribetype.Subtitles{
			Formats: []transcribetype.SubtitleFormat{
				transcribetype.SubtitleFormatSrt,
			},
			OutputStartIndex: aws.Int32(0),
		},
	}
	transcribe, err := a.transcribe.StartTranscriptionJob(ctx, in)
	if err != nil {
		return "", fmt.Errorf("failed to start transcription job: %w", err)
	}
	a.logger.Info("started transcription job", zap.String("broadcastId", broadcast.ID), zap.String("jobName", jobName))

	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("transcription job timeout")
		case <-ticker.C:
			in := &awstranscribe.GetTranscriptionJobInput{
				TranscriptionJobName: transcribe.TranscriptionJob.TranscriptionJobName,
			}
			out, err := a.transcribe.GetTranscriptionJob(ctx, in)
			if err != nil {
				return "", fmt.Errorf("failed to get transcription job: %w", err)
			}

			switch out.TranscriptionJob.TranscriptionJobStatus {
			case transcribetype.TranscriptionJobStatusCompleted:
				return outputFilename, nil
			case transcribetype.TranscriptionJobStatusFailed:
				return "", fmt.Errorf("transcription job failed: reason=%s", aws.ToString(out.TranscriptionJob.FailureReason))
			}

			a.logger.Info("translation job in progress", zap.String("broadcastId", broadcast.ID), zap.String("jobName", jobName))
		}
	}
}

// 日本語テキストを英語に翻訳
func (a *app) executeTranslate(ctx context.Context, broadcast *entity.Broadcast) (string, error) {
	u, err := url.Parse(broadcast.ArchiveURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse archive url: %w", err)
	}

	archiveKey := strings.TrimPrefix(u.Path, "/")

	dir := fmt.Sprintf(entity.BroadcastArchiveTextPath, broadcast.ScheduleID)
	filename := strings.Split(filepath.Base(archiveKey), ".")[0]

	japaneseTextKey := fmt.Sprintf("%s/%s-ja.srt", dir, filename)
	englishTextKey := fmt.Sprintf("%s/%s-en.srt", dir, filename)

	current, err := a.s3.GetMetadata(ctx, englishTextKey)
	if err == nil && current.ContentType == "srt" {
		a.logger.Info("text translation already completed", zap.String("broadcastId", broadcast.ID))
		return englishTextKey, nil
	}

	japanese, err := a.s3.DownloadAndReadAll(ctx, japaneseTextKey)
	if err != nil {
		return "", fmt.Errorf("failed to get object: %w", err)
	}

	buf := &bytes.Buffer{}
	for chunk := range slices.Chunk(japanese, 8000) {
		in := &awstranslate.TranslateTextInput{
			SourceLanguageCode: aws.String("ja"),
			TargetLanguageCode: aws.String("en"),
			Text:               aws.String(string(chunk)),
		}
		out, err := a.translate.TranslateText(ctx, in)
		if err != nil {
			return "", fmt.Errorf("failed to translate text: %w", err)
		}
		if _, err := buf.WriteString(aws.ToString(out.TranslatedText)); err != nil {
			return "", fmt.Errorf("failed to write object: %w", err)
		}
	}

	metadata := map[string]string{
		"Content-Type": "text/plain",
	}
	if _, err := a.s3.Upload(ctx, englishTextKey, buf, metadata); err != nil {
		return "", fmt.Errorf("failed to put translated object: %w", err)
	}

	return englishTextKey, nil
}

func (a *app) uploadFixedArchive(ctx context.Context, broadcast *entity.Broadcast, japaneseTextKey, englishTextKey string) (string, error) {
	u, err := url.Parse(broadcast.ArchiveURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse archive url: %w", err)
	}
	archiveURL, err := a.s3.GenerateObjectURL(u.Path)
	if err != nil {
		return "", fmt.Errorf("failed to generate object url from archive: %w", err)
	}
	japaneseTextURL, err := a.s3.GenerateObjectURL(japaneseTextKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate object url from japanese text key: %w", err)
	}
	englishTextURL, err := a.s3.GenerateObjectURL(englishTextKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate object url from english text key: %w", err)
	}

	args := []string{
		"-i", archiveURL,
		"-i", japaneseTextURL,
		"-i", englishTextURL,
		"-map", "0:v",
		"-map", "0:a",
		"-map", "1",
		"-map", "2",
		"-metadata:s:s:0", "language=jpn",
		"-metadata:s:s:1", "language=eng",
		"-c:v", "copy",
		"-c:a", "copy",
		"-c:s", "mov_text",
		"/tmp/output.mp4",
	}

	buf := &bytes.Buffer{}
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stderr = buf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute command. err=%s: %w", buf.String(), err)
	}
	a.logger.Info("finished execute ffmpeg", zap.String("broadcastId", broadcast.ID))

	f, err := os.Open("/tmp/output.mp4")
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	dir := fmt.Sprintf(entity.BroadcastArchiveMP4Path, broadcast.ScheduleID)
	key := fmt.Sprintf("%s/%s.mp4", dir, uuid.Base58Encode(uuid.New()))

	metadata := map[string]string{
		"Content-Type":  "video/mp4",
		"Cache-Control": "s-maxage=" + entity.BroadcastArchiveRegulation.CacheTTL.String(),
	}
	if _, err := a.s3.Upload(ctx, key, f, metadata); err != nil {
		return "", fmt.Errorf("failed to put fixed archive video: %w", err)
	}
	return key, nil
}
