package mysql

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	image     string
	database  string
	username  string
	password  string
	now       func() time.Time
	schemaDir string
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

// WithSchemaDir はスキーマSQLファイルのディレクトリパスを指定する。
// 指定されたディレクトリ内の .sql ファイルをファイル名順（バージョン順）に読み込み、
// コンテナDB起動後に自動的に実行してスキーマを初期化する。
func WithSchemaDir(dir string) ContainerDBOption {
	return func(c *containerDBConfig) {
		c.schemaDir = dir
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

	// スキーマディレクトリが指定されている場合、SQLファイルを実行してスキーマを初期化する
	if conf.schemaDir != "" {
		if err := execSchemaFiles(ctx, client, conf.schemaDir); err != nil {
			_ = container.Terminate(ctx)
			return nil, nil, fmt.Errorf("mysql testcontainer: failed to initialize schema: %w", err)
		}
	}

	cleanup := func() {
		_ = container.Terminate(context.Background())
	}

	return client, cleanup, nil
}

// execSchemaFiles は指定ディレクトリ内の .sql ファイルをファイル名順に読み込み実行する。
// マルチステートメントSQLに対応するため、underlying *sql.DB を使って実行する。
func execSchemaFiles(ctx context.Context, client *Client, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read schema dir %q: %w", dir, err)
	}

	// .sql ファイルのみ抽出してファイル名順にソート
	var sqlFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		return nil
	}

	// underlying *sql.DB を取得（multiStatements 対応の接続を作成）
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	for _, file := range sqlFiles {
		path := filepath.Join(dir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read schema file %q: %w", path, err)
		}

		sql := string(content)
		if strings.TrimSpace(sql) == "" {
			continue
		}

		// multiStatements=true が DSN に含まれていない場合でも動作するよう、
		// セミコロンで分割して個別に実行する
		statements := splitSQLStatements(sql)
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := sqlDB.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("failed to execute schema file %q: %w", file, err)
			}
		}
	}

	return nil
}

// splitSQLStatements はSQLテキストをセミコロンで分割する。
// 文字列リテラル内のセミコロンは無視する簡易パーサー。
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]

		// エスケープ文字の処理
		if (inSingleQuote || inDoubleQuote) && ch == '\\' && i+1 < len(sql) {
			current.WriteByte(ch)
			i++
			current.WriteByte(sql[i])
			continue
		}

		switch ch {
		case '\'':
			if !inDoubleQuote && !inBacktick {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote && !inBacktick {
				inDoubleQuote = !inDoubleQuote
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote && !inBacktick {
				stmt := strings.TrimSpace(current.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				current.Reset()
				continue
			}
		}

		current.WriteByte(ch)
	}

	// 最後のステートメント（セミコロンなしの場合）
	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

// FindProjectRoot はカレントディレクトリから親ディレクトリを辿り、
// go.mod が存在するディレクトリをプロジェクトルートとして返す。
// テストファイルからスキーマディレクトリの相対パスを解決するために使用する。
func FindProjectRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %q", start)
		}
		dir = parent
	}
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
