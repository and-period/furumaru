package cmd

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	AppName              string `envconfig:"APP_NAME" default:"admin-gateway"`
	Environment          string `envconfig:"ENV" default:"none"`
	Port                 int64  `envconfig:"PORT" default:"8080"`
	MetricsPort          int64  `envconfig:"METRICS_PORT" default:"9090"`
	ShutdownDelaySec     int64  `envconfig:"SHUTDOWN_DELAY_SEC" default:"20"`
	LogPath              string `envconfig:"LOG_PATH" default:""`
	LogLevel             string `envconfig:"LOG_LEVEL" default:"info"`
	DBSocket             string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost               string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort               string `envconfig:"DB_PORT" default:"3306"`
	DBUsername           string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword           string `envconfig:"DB_PASSWORD" default:""`
	DBTimeZone           string `envconfig:"DB_TIMEZONE" default:""`
	DBEnabledTLS         bool   `envconfig:"DB_ENABLED_TLS" default:"false"`
	DBSecretName         string `envconfig:"DB_SECRET_NAME" default:""`
	NewRelicLicense      string `envconfig:"NEW_RELIC_LICENSE" default:""`
	NewRelicSecretName   string `envconfig:"NEW_RELIC_SECRET_NAME" default:""`
	StripeSecretKey      string `envconfig:"STRIPE_SECRET_KEY" default:""`
	StripeWebhookKey     string `envconfig:"STRIPE_WEBHOOK_KEY" default:""`
	StripeSecretName     string `envconfig:"STRIPE_SECRET_NAME" default:""`
	AWSRegion            string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	S3Bucket             string `envconfig:"S3_BUCKET" default:""`
	CognitoAdminPoolID   string `envconfig:"COGNITO_Admin_POOL_ID" default:""`
	CognitoAdminClientID string `envconfig:"COGNITO_Admin_CLIENT_ID" default:""`
	CognitoUserPoolID    string `envconfig:"COGNITO_USER_POOL_ID" default:""`
	CognitoUserClientID  string `envconfig:"COGNITO_USER_CLIENT_ID" default:""`
	SQSQueueURL          string `envconfig:"SQS_QUEUE_URL" default:""`
	SQSMockEnabled       bool   `envconfig:"SQS_MOCK_ENABLED" default:"false"`
	AminWebURL           string `envconfig:"ADMIN_WEB_URL" default:""`
	UserWebURL           string `envconfig:"USER_WEB_URL" default:""`
	SlackAPIToken        string `envconfig:"SLACK_API_TOKEN" default:""`
	SlackChannelID       string `envconfig:"SLACK_CHANNEL_ID" default:""`
	SlackSecretName      string `envconfig:"SLACK_SECRET_NAME" default:""`
	RBACPolicyPath       string `envconfig:"RBAC_POLICY_PATH" default:""`
	RBACModelPath        string `envconfig:"RBAC_MODEL_PATH" default:""`
}

func newConfig() (*config, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return conf, fmt.Errorf("config: failed to new config: %w", err)
	}
	return conf, nil
}
