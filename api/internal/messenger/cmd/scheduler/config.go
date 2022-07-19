package cmd

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	LogPath        string `envconfig:"LOG_PATH" default:""`
	LogLevel       string `envconfig:"LOG_LEVEL" default:"info"`
	DBSocket       string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost         string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort         string `envconfig:"DB_PORT" default:"3306"`
	DBUsername     string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword     string `envconfig:"DB_PASSWORD" default:""`
	DBTimeZone     string `envconfig:"DB_TIMEZONE" default:""`
	DBEnabledTLS   bool   `envconfig:"DB_ENABLED_TLS" default:"false"`
	DBSecretName   string `envconfig:"DB_SECRET_NAME" default:""`
	AWSRegion      string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	SQSQueueURL    string `envconfig:"SQS_QUEUE_URL" default:""`
	SQSMockEnabled bool   `envconfig:"SQS_MOCK_ENABLED" default:"false"`
	AminWebURL     string `envconfig:"ADMIN_WEB_URL" default:""`
	UserWebURL     string `envconfig:"USER_WEB_URL" default:""`
	TargetDatetime string `envconfig:"TARGET_DATETIME" default:""`
}

func newConfig() (*config, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return conf, fmt.Errorf("config: failed to new config: %w", err)
	}
	return conf, nil
}
