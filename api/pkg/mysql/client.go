package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	dmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Client - DB操作用のクライアント構造体
type Client struct {
	DB         *gorm.DB
	maxRetries uint64
}

type Params struct {
	Socket   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type options struct {
	now                  func() time.Time
	location             *time.Location
	charset              string
	collation            string
	enabledTLS           bool
	allowNativePasswords bool
	maxAllowedPacket     int
	maxRetries           int
	maxOpenConns         int
	maxIdleConns         int
	maxConnLifetime      time.Duration
	maxConnIdleTime      time.Duration
	dialTimeout          time.Duration
	readTimeout          time.Duration
	writeTimeout         time.Duration
}

type Option func(opts *options)

func WithNow(now func() time.Time) Option {
	return func(opts *options) {
		opts.now = now
	}
}

func WithLocation(location *time.Location) Option {
	return func(opts *options) {
		opts.location = location
	}
}

func WithCharset(charset string) Option {
	return func(opts *options) {
		opts.charset = charset
	}
}

func WithCollation(collation string) Option {
	return func(opts *options) {
		opts.collation = collation
	}
}

func WithTLS(enable bool) Option {
	return func(opts *options) {
		opts.enabledTLS = enable
	}
}

func WithNativePasswords(enable bool) Option {
	return func(opts *options) {
		opts.allowNativePasswords = enable
	}
}

func WithMaxAllowedPacket(size int) Option {
	return func(opts *options) {
		opts.maxAllowedPacket = size
	}
}

func WithMaxRetries(retries int) Option {
	return func(opts *options) {
		opts.maxRetries = retries
	}
}

func WithMaxConnLifetime(d time.Duration) Option {
	return func(opts *options) {
		opts.maxConnLifetime = d
	}
}

func WithMaxConnIdleTime(d time.Duration) Option {
	return func(opts *options) {
		opts.maxConnIdleTime = d
	}
}

func WithMaxOpenConns(n int) Option {
	return func(opts *options) {
		opts.maxOpenConns = n
	}
}

func WithMaxIdleConns(n int) Option {
	return func(opts *options) {
		opts.maxIdleConns = n
	}
}

func WithDialTimeout(d time.Duration) Option {
	return func(opts *options) {
		opts.dialTimeout = d
	}
}

func WithReadTimeout(d time.Duration) Option {
	return func(opts *options) {
		opts.readTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) Option {
	return func(opts *options) {
		opts.writeTimeout = d
	}
}

// NewClient - DBクライアントの構造体
func NewClient(params *Params, opts ...Option) (*Client, error) {
	dopts := &options{
		now:                  time.Now,
		location:             time.UTC,
		charset:              "utf8mb4",
		collation:            "utf8mb4_general_ci",
		enabledTLS:           false,
		allowNativePasswords: true,
		maxAllowedPacket:     4194304, // 4MiB
		maxRetries:           3,
		maxConnLifetime:      5 * time.Minute,
		maxConnIdleTime:      5 * time.Minute,
	}
	for i := range opts {
		opts[i](dopts)
	}

	// プライマリレプリカの作成
	db, err := newDBClient(params, dopts)
	if err != nil {
		return nil, err
	}
	sql, err := db.DB()
	if err != nil {
		return nil, err
	}
	sql.SetConnMaxLifetime(dopts.maxConnLifetime)
	sql.SetConnMaxIdleTime(dopts.maxConnIdleTime)

	c := &Client{
		DB:         db,
		maxRetries: uint64(dopts.maxRetries),
	}
	return c, nil
}

// Begin - トランザクションの開始処理
func (c *Client) Begin(ctx context.Context, opts ...*sql.TxOptions) (*gorm.DB, error) {
	tx := c.DB.WithContext(ctx).Begin(opts...)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return tx, nil
}

// Close - トランザクションの終了処理
func (c *Client) Close(tx *gorm.DB) func() {
	return func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}
}

// Transaction - トランザクション処理
func (c *Client) Transaction(ctx context.Context, f func(tx *gorm.DB) error) (err error) {
	txFn := func() error {
		tx, err := c.Begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				return
			}
			if err == nil {
				err = tx.Commit().Error
				return
			}
			if !IsRetryable(err) {
				err = backoff.Permanent(err)
			}
			tx.Rollback()
		}()
		return f(tx)
	}
	retry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), c.maxRetries)
	return backoff.Retry(txFn, backoff.WithContext(retry, ctx))
}

// Statement - セレクトクエリの生成
func (c *Client) Statement(ctx context.Context, tx *gorm.DB, table string, fields ...string) *gorm.DB {
	stmt := tx.WithContext(ctx).Table(table)
	if len(fields) == 0 {
		stmt = stmt.Select("*")
	} else {
		stmt = stmt.Select(fields)
	}
	return stmt
}

// Statement - カウントクエリの生成
func (c *Client) Count(ctx context.Context, tx *gorm.DB, model interface{}, fn func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64

	stmt := tx.WithContext(ctx).Model(model).Select("COUNT(*)")
	if fn != nil {
		stmt = fn(stmt)
	}
	return total, stmt.Find(&total).Error
}

// IsRetryable - リトライ対象のエラーかの判定
func IsRetryable(err error) bool {
	if errors.Is(err, driver.ErrBadConn) {
		return true
	}
	var e *dmysql.MySQLError
	if errors.As(err, &e) {
		switch e.Number {
		case 1213: // デッドロック
			return true
		case 1105: // TiProxy 一時的エラー
			return strings.Contains(e.Message, "No available TiDB instances") ||
				strings.Contains(e.Message, "TiProxy fails to connect to TiDB")
		}
	}
	return false
}

// Retryable - リトライ可能かの判定
func Retryable(err error) bool {
	var e *dmysql.MySQLError
	if errors.As(err, &e) {
		return e.Number == 1213 // デッドロック
	}
	return false
}

func newDBClient(params *Params, opts *options) (*gorm.DB, error) {
	conf := &gorm.Config{
		NowFunc: opts.now,
	}
	dsn := newDSN(params, opts)
	return gorm.Open(mysql.Open(dsn), conf)
}

func newDSN(params *Params, opts *options) string {
	dsn := &dmysql.Config{
		User:                 params.Username,
		Passwd:               params.Password,
		Net:                  params.Socket,
		Addr:                 fmt.Sprintf("%s:%s", params.Host, params.Port),
		DBName:               params.Database,
		Collation:            opts.collation,
		Loc:                  opts.location,
		MaxAllowedPacket:     opts.maxAllowedPacket,
		ParseTime:            true,
		AllowNativePasswords: opts.allowNativePasswords,
		CheckConnLiveness:    true,
		Params:               map[string]string{"charset": opts.charset},
	}
	if opts.enabledTLS {
		dsn.TLSConfig = "true"
	}
	return dsn.FormatDSN()
}
