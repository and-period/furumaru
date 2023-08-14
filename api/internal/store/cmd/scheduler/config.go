package cmd

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	AppName        string `envconfig:"APP_NAME" default:"broadcast-scheduler"`
	Environment    string `envconfig:"ENV" default:"none"`
	RunMethod      string `envconfig:"RUN_METHOD" default:"lambda"`
	RunType        string `envconfig:"RUN_TYPE" default:""`
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
	TargetDatetime string `envconfig:"TARGET_DATETIME" default:""`
}

func newConfig() (*config, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return conf, fmt.Errorf("config: failed to new config: %w", err)
	}
	return conf, nil
}
