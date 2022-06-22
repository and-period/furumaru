package cmd

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	RunMethod            string `envconfig:"RUN_METHOD" default:"lambda"`
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
	SendGridAPIKey       string `envconfig:"SENDGRID_API_KEY" default:""`
	SendGridTemplatePath string `envconfig:"SENDGRID_TEMPLATE_PATH" default:""`
	MailFromName         string `envconfig:"MAIL_FROM_NAME" default:""`
	MailFromAddress      string `envconfig:"MAIL_FROM_ADDRESS" default:""`
}

func newConfig() (*config, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return conf, fmt.Errorf("config: failed to new config: %w", err)
	}
	return conf, nil
}
