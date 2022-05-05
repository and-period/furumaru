package cmd

import (
	"context"
	"sync"

	v1 "github.com/and-period/marche/api/internal/gateway/user/v1/handler"
	storedb "github.com/and-period/marche/api/internal/store/database"
	store "github.com/and-period/marche/api/internal/store/service"
	userdb "github.com/and-period/marche/api/internal/user/database"
	user "github.com/and-period/marche/api/internal/user/service"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/zap"
)

type registry struct {
	v1 v1.APIV1Handler
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

	userService, err := newUserService(ctx, conf, dopts)
	if err != nil {
		return nil, err
	}
	storeService, err := newStoreService(ctx, conf, dopts)
	if err != nil {
		return nil, err
	}

	// Handlerの設定
	v1Params := &v1.Params{
		WaitGroup:    &sync.WaitGroup{},
		UserService:  userService,
		StoreService: storeService,
	}

	return &registry{
		v1: v1.NewAPIV1Handler(v1Params, v1.WithLogger(dopts.logger)),
	}, nil
}

func newDatabase(db string, conf *config, opts *options) (*database.Client, error) {
	params := &database.Params{
		Socket:   conf.DBSocket,
		Host:     conf.DBHost,
		Port:     conf.DBPort,
		Database: db,
		Username: conf.DBUsername,
		Password: conf.DBPassword,
	}
	return database.NewClient(
		params,
		database.WithLogger(opts.logger),
		database.WithTLS(conf.DBEnabledTLS),
		database.WithTimeZone(conf.DBTimeZone),
	)
}

func newUserService(ctx context.Context, conf *config, opts *options) (user.UserService, error) {
	// MySQLの設定
	mysql, err := newDatabase("users", conf, opts)
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
		AppClientID:     conf.CognitoUserClientID,
		AppClientSecret: conf.CognitoUserClientSecret,
	}

	// Databaseの設定
	dbParams := &userdb.Params{
		Database: mysql,
	}

	// User Serviceの設定
	params := &user.Params{
		Database:  userdb.NewDatabase(dbParams),
		AdminAuth: cognito.NewClient(awscfg, &cognito.Params{}),
		ShopAuth:  cognito.NewClient(awscfg, &cognito.Params{}),
		UserAuth:  cognito.NewClient(awscfg, userAuthParams),
	}
	return user.NewUserService(
		params,
		user.WithLogger(opts.logger),
	), nil
}

func newStoreService(ctx context.Context, conf *config, opts *options) (store.StoreService, error) {
	// MySQLの設定
	mysql, err := newDatabase("stores", conf, opts)
	if err != nil {
		return nil, err
	}

	// Databaseの設定
	dbParams := storedb.Params{
		Database: mysql,
	}

	// Store Serviceの設定
	params := &store.Params{
		Database: storedb.NewDatabase(&dbParams),
	}
	return store.NewStoreService(
		params,
		store.WithLogger(opts.logger),
	), nil
}
