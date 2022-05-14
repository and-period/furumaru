package cmd

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port                    int64  `envconfig:"PORT" default:"8080"`
	MetricsPort             int64  `envconfig:"METRICS_PORT" default:"9090"`
	ShutdownDelaySec        int64  `envconfig:"SHUTDOWN_DELAY_SEC" default:"20"`
	LogPath                 string `envconfig:"LOG_PATH" default:""`
	LogLevel                string `envconfig:"LOG_LEVEL" default:"info"`
	DBStoreSocket           string `envconfig:"DB_STORE_SOCKET" default:"tcp"`
	DBStoreHost             string `envconfig:"DB_STORE_HOST" default:"127.0.0.1"`
	DBStorePort             string `envconfig:"DB_STORE_PORT" default:"3306"`
	DBStoreUsername         string `envconfig:"DB_STORE_USERNAME" default:"root"`
	DBStorePassword         string `envconfig:"DB_STORE_PASSWORD" default:""`
	DBStoreTimeZone         string `envconfig:"DB_STORE_TIMEZONE" default:""`
	DBStoreEnabledTLS       bool   `envconfig:"DB_STORE_ENABLED_TLS" default:"false"`
	DBUserSocket            string `envconfig:"DB_USER_SOCKET" default:"tcp"`
	DBUserHost              string `envconfig:"DB_USER_HOST" default:"127.0.0.1"`
	DBUserPort              string `envconfig:"DB_USER_PORT" default:"3306"`
	DBUserUsername          string `envconfig:"DB_USER_USERNAME" default:"root"`
	DBUserPassword          string `envconfig:"DB_USER_PASSWORD" default:""`
	DBUserTimeZone          string `envconfig:"DB_USER_TIMEZONE" default:""`
	DBUserEnabledTLS        bool   `envconfig:"DB_USER_ENABLED_TLS" default:"false"`
	AWSRegion               string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	AWSAccessKey            string `envconfig:"AWS_ACCESS_KEY" default:""`
	AWSSecretKey            string `envconfig:"AWS_SECRET_KEY" default:""`
	S3BucketName            string `envconfig:"S3_BUCKET_NAME" default:""`
	SendGridAPIKey          string `envconfig:"SENDGRID_API_KEY" default:""`
	SendGridTemplatePath    string `envconfig:"SENDGRID_TEMPLATE_PATH" default:""`
	MailFromName            string `envconfig:"MAIL_FROM_NAME" default:""`
	MailFromAddress         string `envconfig:"MAIL_FROM_ADDRESS" default:""`
	UserWebURL              string `envconfig:"USER_WEB_URL" default:""`
	CognitoUserPoolID       string `envconfig:"COGNITO_USER_POOL_ID" default:""`
	CognitoUserClientID     string `envconfig:"COGNITO_USER_CLIENT_ID" default:""`
	CognitoUserClientSecret string `envconfig:"COGNITO_USER_CLIENT_SECRET" default:""`
}

func newConfig() (*config, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return conf, fmt.Errorf("config: failed to new config: %w", err)
	}
	return conf, nil
}
