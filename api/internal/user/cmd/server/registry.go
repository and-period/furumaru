package cmd

import (
	"context"

	"github.com/and-period/marche/api/internal/user/api"
	db "github.com/and-period/marche/api/internal/user/database"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/proto/user"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/zap"
)

type registry struct {
	userServer user.UserServiceServer
}

type options struct {
	logger *zap.Logger
}

type option func(opts *options)

func withLogger(logger *zap.Logger) option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func newRegistry(ctx context.Context, conf *config, opts ...option) (*registry, error) {
	// オプション設定の取得
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}

	// MySQLの設定
	mysqlParams := &database.Params{
		Socket:     conf.DBSocket,
		Host:       conf.DBHost,
		Port:       conf.DBPort,
		Database:   conf.DBDatabase,
		Username:   conf.DBUsername,
		Password:   conf.DBPassword,
		TimeZone:   conf.DBTimeZone,
		EnabledTLS: conf.DBEnabledTLS,
		Logger:     dopts.logger,
	}
	mysql, err := database.NewClient(mysqlParams)
	if err != nil {
		return nil, err
	}

	// Amazon Cognitoの設定
	awscreds := aws.NewCredentialsCache(
		awscredentials.NewStaticCredentialsProvider(conf.AWSAccessKey, conf.AWSSecretKey, ""),
	)
	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(conf.AWSRegion),
		awsconfig.WithCredentialsProvider(awscreds),
	)
	if err != nil {
		return nil, err
	}
	userAuthParams := &cognito.Params{
		UserPoolID:      conf.CognitoUserPoolID,
		AppClientID:     conf.CognitoClientID,
		AppClientSecret: conf.CognitoClientSecret,
	}
	userAuth := cognito.NewClient(awscfg, userAuthParams)

	// Databaseの設定
	dbParams := &db.Params{
		Database: mysql,
	}

	// User Serviceの設定
	apiParams := &api.Params{
		Logger:   dopts.logger,
		Database: db.NewDatabase(dbParams),
		UserAuth: userAuth,
	}

	return &registry{
		userServer: api.NewUserService(apiParams),
	}, nil
}
