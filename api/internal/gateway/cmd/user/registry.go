package user

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/auth"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type params struct {
	serviceName              string
	waitGroup                *sync.WaitGroup
	aws                      aws.Config
	secret                   secret.Client
	storage                  storage.Bucket
	tmpStorage               storage.Bucket
	userAuth                 cognito.Client
	cache                    dynamodb.Client
	producer                 sqs.Producer
	slack                    slack.Client
	newRelic                 *newrelic.Application
	sentry                   sentry.Client
	komoju                   *komoju.Komoju
	adminWebURL              *url.URL
	userWebURL               *url.URL
	postalCode               postalcode.Client
	geolocation              geolocation.Client
	liffVerifier             auth.OIDCVerifier[auth.LIFFClaims]
	jwtVerifier              auth.JWTVerifier
	jwtGenerator             auth.JWTGenerator
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
	googleMapsPlatformAPIKey string
	jwtSecret                string
}

func (a *app) inject(ctx context.Context) error {
	params := &params{
		serviceName: fmt.Sprintf("%s-%s", a.AppName, a.Environment),
		now:         jst.Now,
		waitGroup:   &sync.WaitGroup{},
		debugMode:   a.LogLevel == "debug",
	}

	// AWS関連の初期化
	if err := a.injectAWS(ctx, params); err != nil {
		return err
	}

	// 認証クライアントの初期化
	a.injectAuth(params)

	// 外部サービスの初期化
	if err := a.injectExternal(ctx, params); err != nil {
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
