// データベースにスキーマを適用します
//
//	usage: go run ./database-migrate-tidb/main.go \
//	 -db-host='127.0.0.1' -db-port='3316' \
//	 -db-username='root' -db-password='12345678'
//
// 事前に以下DDLを実行すること

/**
 * CREATE SCHEMA IF NOT EXISTS `migrations` DEFAULT CHARACTER SET utf8mb4;
 *
 * CREATE TABLE IF NOT EXISTS `migrations`.`schemas` (
 *   `database` varchar(256) NOT NULL,
 *   `version` varchar(10) NOT NULL,
 *   `filename` varchar(256) NOT NULL,
 *   `created_at` int NOT NULL,
 *   `updated_at` int NOT NULL,
 *   PRIMARY KEY (`database`, `version`)
 * );
 */
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

const (
	migrateDB   = "migrations"
	schemaTable = "schemas"
)

var (
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
	logger   *zap.Logger
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
	flag.StringVar(&host, "db-host", "mysql", "target mysql host")
	flag.StringVar(&port, "db-port", "3306", "target mysql port")
	flag.StringVar(&username, "db-username", "root", "target mysql username")
	flag.StringVar(&password, "db-password", "12345678", "target mysql password")
	flag.StringVar(&srcDir, "src", "./../infra/tidb/schema", "ddl source directory")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &app{
		logger:   logger,
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

func (a *app) run(ctx context.Context) error {
	// DDLの管理用DBの接続
	migrate, err := a.setup(migrateDB)
	if err != nil {
		return fmt.Errorf("failed to setup. database=%s: %w", migrateDB, err)
	}

	a.logger.Info("start to apply schema")

	// データベースごとにDDLを適用
	for _, database := range databases {
		a.logger.Info("start to apply schema", zap.String("database", database))

		if err := a.execute(ctx, migrate, database); err != nil {
			return fmt.Errorf("failed to execute. database=%s: %w", database, err)
		}

		a.logger.Info("finish to apply schema", zap.String("database", database))
	}

	a.logger.Info("finish to apply schema")
	return nil
}

/**
 * instance
 */
func (a *app) setup(database string) (*mysql.Client, error) {
	params := &mysql.Params{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
	return mysql.NewTiDBClient(params, mysql.WithLogger(a.logger))
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

//nolint:unparam
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
		isApplied, err := a.getSchema(migrateTx, schemas[i])
		if err != nil {
			a.logger.Error("failed to get schema", zap.Error(err))
			return a.rollback(tx, err)
		}
		if isApplied {
			a.logger.Info("already applied schema", zap.String("filename", schemas[i].filename))
			continue
		}

		a.logger.Info("applying schema...", zap.String("filename", schemas[i].filename))
		if err := a.applySchema(migrateTx, tx, schemas[i]); err != nil {
			a.logger.Error("failed to apply schema", zap.Error(err))
			return a.rollback(tx, err)
		}
		a.logger.Info("applied schema", zap.String("filename", schemas[i].filename))
	}

	if err := tx.Commit(); err != nil {
		a.logger.Error("failed to commit transaction", zap.String("database", database), zap.Error(err))
		return a.rollback(tx, err)
	}

	if err := migrateTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. database=%s: %w", migrateDB, err)
	}
	return nil
}

func (a *app) getSchema(tx *sql.Tx, schema *schema) (bool, error) {
	const format = "SELECT `database`, `version` FROM `%s` WHERE `database` = '%s' AND `version` = '%s' LIMIT 1"
	stmt := fmt.Sprintf(format, schemaTable, schema.database, schema.version)
	rs, err := tx.Query(stmt)
	a.logger.Debug("get schema", zap.String("stmt", stmt), zap.Error(err))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	defer rs.Close()
	return rs.Next(), nil
}

func (a *app) applySchema(migrate, tx *sql.Tx, schema *schema) error {
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
		if _, err := tx.Exec(sql); err != nil {
			return err
		}
	}

	now := jst.Now().Unix()
	const format = "INSERT INTO `%s` (`database`, `version`, `filename`, `created_at`, `updated_at`) VALUES ('%s', '%s', '%s', '%d', '%d')"
	stmt := fmt.Sprintf(format, schemaTable, schema.database, schema.version, schema.filename, now, now)
	if _, err := migrate.Exec(stmt); err != nil {
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
