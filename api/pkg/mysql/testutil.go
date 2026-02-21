package mysql

import (
	"context"
	"fmt"
	"os"
	"time"

	mysqlmodule "github.com/testcontainers/testcontainers-go/modules/mysql"
)

const (
	defaultMySQLImage    = "mysql:8.0"
	defaultMySQLDatabase = "test"
	defaultMySQLUsername = "root"
	defaultMySQLPassword = "password"
)

// ContainerDBOption はコンテナDB起動時のオプションを設定する関数型
type ContainerDBOption func(*containerDBConfig)

type containerDBConfig struct {
	image    string
	database string
	username string
	password string
	now      func() time.Time
}

// WithContainerImage はコンテナイメージを指定する
func WithContainerImage(image string) ContainerDBOption {
	return func(c *containerDBConfig) {
		c.image = image
	}
}

// WithContainerDatabase はデータベース名を指定する
func WithContainerDatabase(database string) ContainerDBOption {
	return func(c *containerDBConfig) {
		c.database = database
	}
}

// WithContainerUsername はユーザー名を指定する
func WithContainerUsername(username string) ContainerDBOption {
	return func(c *containerDBConfig) {
		c.username = username
	}
}

// WithContainerPassword はパスワードを指定する
func WithContainerPassword(password string) ContainerDBOption {
	return func(c *containerDBConfig) {
		c.password = password
	}
}

// WithContainerNow はテスト用のNow関数を指定する
func WithContainerNow(now func() time.Time) ContainerDBOption {
	return func(c *containerDBConfig) {
		c.now = now
	}
}

// NewContainerDB は testcontainers-go を使ってMySQLコンテナを起動し、
// 接続済みの *Client とクリーンアップ関数を返す。
//
// 環境変数 DISABLE_CONTAINER_DB=true が設定されている場合は、
// 従来の環境変数ベースの外部DB接続にフォールバックする。
func NewContainerDB(ctx context.Context, opts ...ContainerDBOption) (*Client, func(), error) {
	if os.Getenv("DISABLE_CONTAINER_DB") == "true" {
		return newExternalDB()
	}

	conf := &containerDBConfig{
		image:    defaultMySQLImage,
		database: defaultMySQLDatabase,
		username: defaultMySQLUsername,
		password: defaultMySQLPassword,
	}
	for _, opt := range opts {
		opt(conf)
	}

	container, err := mysqlmodule.Run(ctx, conf.image,
		mysqlmodule.WithDatabase(conf.database),
		mysqlmodule.WithUsername(conf.username),
		mysqlmodule.WithPassword(conf.password),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("mysql testcontainer: failed to start: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("mysql testcontainer: failed to get host: %w", err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("mysql testcontainer: failed to get port: %w", err)
	}

	params := &Params{
		Socket:   "tcp",
		Host:     host,
		Port:     port.Port(),
		Database: conf.database,
		Username: conf.username,
		Password: conf.password,
	}

	var clientOpts []Option
	if conf.now != nil {
		clientOpts = append(clientOpts, WithNow(conf.now))
	}
	clientOpts = append(clientOpts,
		WithNativePasswords(true),
		WithMaxRetries(3),
	)

	client, err := NewClient(params, clientOpts...)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("mysql testcontainer: failed to create client: %w", err)
	}

	cleanup := func() {
		_ = container.Terminate(context.Background())
	}

	return client, cleanup, nil
}

// newExternalDB は環境変数ベースで外部DBに接続する（従来方式のフォールバック）
func newExternalDB() (*Client, func(), error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	database := os.Getenv("DB_DATABASE")
	if database == "" {
		database = "test"
	}
	username := os.Getenv("DB_USERNAME")
	if username == "" {
		username = "root"
	}
	password := os.Getenv("DB_PASSWORD")

	params := &Params{
		Socket:   "tcp",
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}

	driver := os.Getenv("DB_DRIVER")
	var client *Client
	var err error
	switch driver {
	case "tidb":
		client, err = NewTiDBClient(params)
	default:
		client, err = NewClient(params)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("external db: failed to create client: %w", err)
	}

	// 外部DBの場合はクリーンアップ不要
	noop := func() {}
	return client, noop, nil
}
