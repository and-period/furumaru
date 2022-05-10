// 管理者を登録します
// usage: go run ./main.go \
//          -db-host=127.0.0.1 -db-port=3316 -db-password=1234567 \
//          -aws-access-key=xxx -aws-secret-key=xxx \
//          -cognito-client-id=xxx -cognito-pool-id=xxx \
//          -email=test-admin@and-period.jp
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/and-period/marche/api/internal/user/database"
	"github.com/and-period/marche/api/internal/user/entity"
	user "github.com/and-period/marche/api/internal/user/service"
	"github.com/and-period/marche/api/pkg/cognito"
	db "github.com/and-period/marche/api/pkg/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
)

const (
	dbName    = "users"
	awsRegion = "ap-northeast-1"
)

type app struct {
	db   *db.Client
	auth cognito.Client
	user user.UserService
}

func main() {
	start := time.Now()
	fmt.Println("Start..")
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Printf("Done: %s\n", time.Since(start))
}

func run() error {
	var (
		dbHost, dbPort             string
		dbUsername, dbPassword     string
		dbEnabledTLS               bool
		awsAccessKey, awsSecretKey string
		authClientID, authPoolID   string
		email                      string
		err                        error
	)

	app := app{}
	flag.StringVar(&dbHost, "db-host", "mysql", "target mysql host")
	flag.StringVar(&dbPort, "db-port", "3306", "target mysql port")
	flag.StringVar(&dbUsername, "db-username", "root", "target mysql username")
	flag.StringVar(&dbPassword, "db-password", "12345678", "target mysql password")
	flag.BoolVar(&dbEnabledTLS, "db-enabled-tls", false, "target mysql enabled tls")
	flag.StringVar(&awsAccessKey, "aws-access-key", "", "aws access key for cognito")
	flag.StringVar(&awsSecretKey, "aws-secret-key", "", "aws secret key for cognito")
	flag.StringVar(&authClientID, "cognito-client-id", "", "target cognito client id")
	flag.StringVar(&authPoolID, "cognito-pool-id", "", "target cognito user pool id")
	flag.StringVar(&email, "email", "", "target email for created admin")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.db, err = app.setupDB(dbHost, dbPort, dbUsername, dbPassword, dbEnabledTLS)
	if err != nil {
		return err
	}
	app.auth, err = app.setupAuth(ctx, awsAccessKey, awsSecretKey, authClientID, authPoolID)
	if err != nil {
		return err
	}
	app.user = app.newUserService()

	in := &user.CreateAdminInput{
		Email: email,
		Role:  entity.AdminRoleAdministrator,
	}
	_, err = app.user.CreateAdmin(ctx, in)
	return err
}

func (a *app) newUserService() user.UserService {
	params := &user.Params{
		Database:  database.NewDatabase(&database.Params{Database: a.db}),
		AdminAuth: a.auth,
	}
	return user.NewUserService(params)
}

func (a *app) setupDB(host, port, username, password string, tls bool) (*db.Client, error) {
	params := &db.Params{
		Socket:   "tcp",
		Host:     host,
		Port:     port,
		Database: dbName,
		Username: username,
		Password: password,
	}
	return db.NewClient(params, db.WithTLS(tls))
}

func (a *app) setupAuth(ctx context.Context, accessKey, secretKey, clientID, poolID string) (cognito.Client, error) {
	awscreds := aws.NewCredentialsCache(
		awscredentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	)
	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(awsRegion),
		awsconfig.WithCredentialsProvider(awscreds),
	)
	if err != nil {
		return nil, err
	}
	params := &cognito.Params{
		UserPoolID:  poolID,
		AppClientID: clientID,
	}
	return cognito.NewClient(awscfg, params), nil
}
