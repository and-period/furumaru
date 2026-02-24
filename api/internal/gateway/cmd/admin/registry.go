package admin

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/youtube"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/payment"
	"github.com/and-period/furumaru/api/pkg/batch"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type params struct {
	waitGroup                *sync.WaitGroup
	aws                      aws.Config
	secret                   secret.Client
	storage                  storage.Bucket
	tmpStorage               storage.Bucket
	adminAuth                cognito.Client
	userAuth                 cognito.Client
	cache                    dynamodb.Client
	messengerQueue           sqs.Producer
	mediaQueue               sqs.Producer
	batch                    batch.Client
	medialive                medialive.MediaLive
	youtube                  youtube.Youtube
	slack                    slack.Client
	newRelic                 *newrelic.Application
	sentry                   sentry.Client
	providers                map[entity.PaymentProviderType]payment.Provider
	komojuWebhookSecret      string
	stripeSecretKey          string
	stripeWebhookSecret      string
	adminWebURL              *url.URL
	userWebURL               *url.URL
	assetsURL                *url.URL
	postalCode               postalcode.Client
	geolocation              geolocation.Client
	now                      func() time.Time
	debugMode                bool
	tidbHost                 string
	tidbPort                 string
	tidbUsername             string
	tidbPassword             string
	slackToken               string
	slackChannelID           string
	newRelicLicense          string
	sentryDsn                string
	komojuClientID           string
	komojuClientPassword     string
	googleClientID           string
	googleClientSecret       string
	googleMapsPlatformAPIKey string
}

//nolint:funlen,maintidx
func (a *app) inject(ctx context.Context) error {
	params := &params{
		now:       jst.Now,
		waitGroup: &sync.WaitGroup{},
		debugMode: a.LogLevel == "debug",
	}

	// AWS関連の初期化
	if err := a.injectAWS(ctx, params); err != nil {
		return err
	}

	// 認証クライアントの初期化
	a.injectAuth(params)

	// 外部サービスの初期化
	if err := a.injectExternal(params); err != nil {
		return err
	}

	// サービス層・ハンドラーの初期化
	if err := a.injectServices(params); err != nil {
		return err
	}

	a.debugMode = params.debugMode
	a.waitGroup = params.waitGroup
	a.slack = params.slack
	a.newRelic = params.newRelic
	return nil
}
