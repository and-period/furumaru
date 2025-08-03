// データベースにスキーマを適用します
//
//	usage: go run ./database-migrate-tidb/main.go \
//	 -db-host='127.0.0.1' -db-port='3316' \
//	 -db-username='root' -db-password='12345678'
//
//nolint:lll
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mysql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	migrateDB      = "migrations"
	schemaTable    = "schemas"
	createDBDDL    = "CREATE SCHEMA IF NOT EXISTS `%s`;"
	schemaTableDDL = "CREATE TABLE IF NOT EXISTS `migrations`.`schemas` (`database` varchar(256) NOT NULL, `version` varchar(10) NOT NULL, `filename` varchar(256) NOT NULL, `created_at` int NOT NULL, `updated_at` int NOT NULL, PRIMARY KEY (`database`, `version`));"
)

var (
	driver   string
	host     string
	port     string
	username string
	password string
	srcDir   string

	databases = []string{
		"users",
		"stores",
		"messengers",
		"media",
	}
)

type app struct {
	host     string
	port     string
	username string
	password string
	srcDir   string
}

func main() {
	startedAt := jst.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	app, err := setup(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup: %v\n", err)
		os.Exit(1)
	}

	if err := app.precheck(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to precheck: %v\n", err)
		os.Exit(1)
	}

	if err := app.run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %v\n", err)
		os.Exit(1)
	}

	endAt := jst.Now()

	const format = "2006-01-02 15:04:05"
	fmt.Printf("Done: %s -> %s (%s)\n",
		jst.Format(startedAt, format), jst.Format(endAt, format),
		endAt.Sub(startedAt).Truncate(time.Second).String(),
	)
}

func setup(_ context.Context) (*app, error) {
	flag.StringVar(&driver, "db-driver", "tidb", "target mysql driver")
	flag.StringVar(&host, "db-host", "127.0.0.1", "target mysql host")
	flag.StringVar(&port, "db-port", "3306", "target mysql port")
	flag.StringVar(&username, "db-username", "root", "target mysql username")
	flag.StringVar(&password, "db-password", "12345678", "target mysql password")
	flag.StringVar(&srcDir, "src", "./../infra/tidb/schema", "ddl source directory")
	flag.Parse()

	return &app{
		host:     host,
		port:     port,
		username: username,
		password: password,
		srcDir:   srcDir,
	}, nil
}

type schema struct {
	version  string
	database string
	filename string
	path     string
}

func (a *app) precheck(ctx context.Context) error {
	client, err := a.setup("")
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	// DDLの管理用DBの作成
	if err := client.DB.WithContext(ctx).Exec(fmt.Sprintf(createDBDDL, migrateDB)).Error; err != nil {
		return fmt.Errorf("failed to create migrate database: %w", err)
	}
	// DDLの管理用テーブルの作成
	if err := client.DB.WithContext(ctx).Exec(schemaTableDDL).Error; err != nil {
		return fmt.Errorf("failed to create schemas table: %w", err)
	}
	// 各種DBの作成
	for _, database := range databases {
		if err := client.DB.WithContext(ctx).Exec(fmt.Sprintf(createDBDDL, database)).Error; err != nil {
			return fmt.Errorf("failed to create database. database=%s: %w", database, err)
		}
	}
	return nil
}

func (a *app) run(ctx context.Context) error {
	// DDLの管理用DBの接続
	migrate, err := a.setup(migrateDB)
	if err != nil {
		return fmt.Errorf("failed to setup. database=%s: %w", migrateDB, err)
	}

	slog.Info("start to apply schema")

	// データベースごとにDDLを適用
	for _, database := range databases {
		slog.Info("start to apply schema", slog.String("database", database))

		if err := a.execute(ctx, migrate, database); err != nil {
			return fmt.Errorf("failed to execute. database=%s: %w", database, err)
		}

		slog.Info("finish to apply schema", slog.String("database", database))
	}

	slog.Info("finish to apply schema")
	return nil
}

/**
 * instance
 */
func (a *app) setup(database string) (*mysql.Client, error) {
	params := &mysql.Params{
		Socket:   "tcp",
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
	switch driver {
	case "mysql":
		return mysql.NewClient(params)
	case "tidb":
		return mysql.NewTiDBClient(params)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}

func (a *app) begin(ctx context.Context, db *mysql.Client) (*sql.Tx, error) {
	sql, err := db.DB.DB()
	if err != nil {
		return nil, err
	}
	tx, err := sql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (a *app) close(tx *sql.Tx) func() {
	return func() {
		if r := recover(); r != nil {
			_ = a.rollback(tx, fmt.Errorf("panic: %v", r))
		}
	}
}

func (a *app) rollback(tx *sql.Tx, err error) error {
	return fmt.Errorf("%w: %s", err, tx.Rollback().Error())
}

func (a *app) execute(ctx context.Context, migrate *mysql.Client, database string) error {
	migrateTx, err := a.begin(ctx, migrate)
	if err != nil {
		return fmt.Errorf("failed to begin transaction. database=%s: %w", migrateDB, err)
	}
	defer a.close(migrateTx)

	// DBへ接続
	db, err := a.setup(database)
	if err != nil {
		return fmt.Errorf("failed to setup: %w", err)
	}

	// DDLの取得
	schemas, err := a.newSchemas(a.srcDir, database)
	if err != nil {
		return fmt.Errorf("failed to get schemas: %w", err)
	}

	tx, err := a.begin(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer a.close(tx)

	// DDLの適用
	for i := range schemas {
		isApplied, err := a.getSchema(ctx, migrateTx, schemas[i])
		if err != nil {
			slog.Error("failed to get schema", log.Error(err))
			return a.rollback(tx, err)
		}
		if isApplied {
			slog.Info("already applied schema", slog.String("filename", schemas[i].filename))
			continue
		}

		slog.Info("applying schema...", slog.String("filename", schemas[i].filename))
		if err := a.applySchema(ctx, migrateTx, tx, schemas[i]); err != nil {
			slog.Error("failed to apply schema", log.Error(err))
			return a.rollback(tx, err)
		}
		slog.Info("applied schema", slog.String("filename", schemas[i].filename))
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", slog.String("database", database), log.Error(err))
		return a.rollback(tx, err)
	}

	if err := migrateTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. database=%s: %w", migrateDB, err)
	}
	return nil
}

func (a *app) getSchema(ctx context.Context, tx *sql.Tx, schema *schema) (bool, error) {
	const format = "SELECT `database`, `version` FROM `%s` WHERE `database` = '%s' AND `version` = '%s' LIMIT 1"
	stmt := fmt.Sprintf(format, schemaTable, schema.database, schema.version)
	rs, err := tx.QueryContext(ctx, stmt)
	slog.Debug("get schema", slog.String("stmt", stmt), log.Error(err))
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer rs.Close()
	return rs.Next(), nil
}

func (a *app) applySchema(ctx context.Context, migrate, tx *sql.Tx, schema *schema) error {
	bytes, err := os.ReadFile(schema.path)
	if err != nil {
		return err
	}

	// 1つのファイルに複数定義が書いてある場合はエラーになるため、分割して適用
	sqls := strings.Split(string(bytes), ";")
	for _, sql := range sqls {
		if sql == "" || sql == "\n" {
			continue // split時、配列の最後に空文字が入るため
		}
		if _, err := tx.ExecContext(ctx, sql); err != nil {
			return err
		}
	}

	now := jst.Now().Unix()
	const format = "INSERT INTO `%s` (`database`, `version`, `filename`, `created_at`, `updated_at`) VALUES ('%s', '%s', '%s', '%d', '%d')"
	stmt := fmt.Sprintf(format, schemaTable, schema.database, schema.version, schema.filename, now, now)
	if _, err := migrate.ExecContext(ctx, stmt); err != nil {
		return err
	}
	return nil
}

func (a *app) newSchemas(dir, database string) ([]*schema, error) {
	path := filepath.Join(dir, database)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	schemas := make([]*schema, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// filename: {version}-{message}.sql
		// e.g.) 2021120902-create-table-teachers.sql
		filename := file.Name()
		strs := strings.Split(filename, "-")

		schema := &schema{
			database: database,
			version:  strs[0],
			filename: filename,
			path:     filepath.Join(path, filename),
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}
