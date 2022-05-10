package cmd

import (
	"context"
	"sync"

	v1 "github.com/and-period/marche/api/internal/gateway/admin/v1/handler"
	storedb "github.com/and-period/marche/api/internal/store/database"
	store "github.com/and-period/marche/api/internal/store/service"
	userdb "github.com/and-period/marche/api/internal/user/database"
	user "github.com/and-period/marche/api/internal/user/service"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/rbac"
	"github.com/and-period/marche/api/pkg/storage"
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

	// Casbinの設定
	enforcer, err := rbac.NewEnforcer(conf.RBACModelPath, conf.RBACPolicyPath)
	if err != nil {
		return nil, err
	}

	// AWS SDKの設定
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

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: conf.S3BucketName,
	}

	// Serviceの設定
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
		Enforcer:     enforcer,
		Storage:      storage.NewBucket(awscfg, storageParams),
		UserService:  userService,
		StoreService: storeService,
	}

	return &registry{
		v1: v1.NewAPIV1Handler(v1Params, v1.WithLogger(dopts.logger)),
	}, nil
}

func newDatabase(params *database.Params, tls bool, timezone string, opts *options) (*database.Client, error) {
	return database.NewClient(
		params,
		database.WithLogger(opts.logger),
		database.WithTLS(tls),
		database.WithTimeZone(timezone),
	)
}

func newUserService(ctx context.Context, conf *config, opts *options) (user.UserService, error) {
	// MySQLの設定
	mysqlParams := &database.Params{
		Socket:   conf.DBUserSocket,
		Host:     conf.DBUserHost,
		Port:     conf.DBUserPort,
		Database: "users",
		Username: conf.DBUserUsername,
		Password: conf.DBUserPassword,
	}
	mysql, err := newDatabase(mysqlParams, conf.DBUserEnabledTLS, conf.DBUserTimeZone, opts)
	if err != nil {
		return nil, err
	}

	// Databaseの設定
	dbParams := &userdb.Params{
		Database: mysql,
	}

	// AWS SDKの設定
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

	// Amazon Cognitoの設定
	adminAuthParams := &cognito.Params{
		UserPoolID:      conf.CognitoAdminPoolID,
		AppClientID:     conf.CognitoAdminClientID,
		AppClientSecret: conf.CognitoAdminClientSecret,
	}
	userAuthParams := &cognito.Params{
		UserPoolID:      conf.CognitoUserPoolID,
		AppClientID:     conf.CognitoUserClientID,
		AppClientSecret: conf.CognitoUserClientSecret,
	}

	// User Serviceの設定
	params := &user.Params{
		Database:  userdb.NewDatabase(dbParams),
		AdminAuth: cognito.NewClient(awscfg, adminAuthParams),
		UserAuth:  cognito.NewClient(awscfg, userAuthParams),
	}
	return user.NewUserService(
		params,
		user.WithLogger(opts.logger),
	), nil
}

func newStoreService(ctx context.Context, conf *config, opts *options) (store.StoreService, error) {
	// MySQLの設定
	mysqlParams := &database.Params{
		Socket:   conf.DBStoreSocket,
		Host:     conf.DBStoreHost,
		Port:     conf.DBStorePort,
		Database: "stores",
		Username: conf.DBStoreUsername,
		Password: conf.DBStorePassword,
	}
	mysql, err := newDatabase(mysqlParams, conf.DBStoreEnabledTLS, conf.DBStoreTimeZone, opts)
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
