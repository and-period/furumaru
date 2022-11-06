package cmd

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	AppName               string `envconfig:"APP_NAME" default:"messenger-worker"`
	Environment           string `envconfig:"ENV" default:"none"`
	RunMethod             string `envconfig:"RUN_METHOD" default:"lambda"`
	LogPath               string `envconfig:"LOG_PATH" default:""`
	LogLevel              string `envconfig:"LOG_LEVEL" default:"info"`
	DBSocket              string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost                string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort                string `envconfig:"DB_PORT" default:"3306"`
	DBUsername            string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword            string `envconfig:"DB_PASSWORD" default:""`
	DBTimeZone            string `envconfig:"DB_TIMEZONE" default:""`
	DBEnabledTLS          bool   `envconfig:"DB_ENABLED_TLS" default:"false"`
	DBSecretName          string `envconfig:"DB_SECRET_NAME" default:""`
	AWSRegion             string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	SendGridAPIKey        string `envconfig:"SENDGRID_API_KEY" default:""`
	SendGridTemplatePath  string `envconfig:"SENDGRID_TEMPLATE_PATH" default:""`
	SendGridSecretName    string `envconfig:"SENDGRID_SECRET_NAME" default:""`
	MailFromName          string `envconfig:"MAIL_FROM_NAME" default:""`
	MailFromAddress       string `envconfig:"MAIL_FROM_ADDRESS" default:""`
	LINEChannelToken      string `envconfig:"LINE_CHANNEL_TOKEN" default:""`
	LINEChannelSecret     string `envconfig:"LINE_CHANNEL_SECRET" default:""`
	LINERoomID            string `envconfig:"LINE_ROOM_ID" default:""`
	LINESecretName        string `envconfig:"LINE_SECRET_NAME" default:""`
	GoogleCredentialsJSON string `envconfig:"GOOGLE_CREDENTIALS_JSON" default:""`
	GoogleSecretName      string `envconfig:"GOOGLE_SECRET_NAME" default:""`
}

func newConfig() (*config, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return conf, fmt.Errorf("config: failed to new config: %w", err)
	}
	return conf, nil
}
