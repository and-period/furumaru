package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dmysql "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

// Client - DB操作用のクライアント構造体
type Client struct {
	DB *gorm.DB
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
	logger               *zap.Logger
	now                  func() time.Time
	location             *time.Location
	charset              string
	collation            string
	enabledTLS           bool
	allowNativePasswords bool
	maxAllowedPacket     int
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

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

// NewClient - DBクライアントの構造体
func NewClient(params *Params, opts ...Option) (*Client, error) {
	dopts := &options{
		logger:               zap.NewNop(),
		now:                  time.Now,
		location:             time.UTC,
		charset:              "utf8mb4",
		collation:            "utf8mb4_general_ci",
		enabledTLS:           false,
		allowNativePasswords: true,
		maxAllowedPacket:     4194304, // 4MiB
	}
	for i := range opts {
		opts[i](dopts)
	}

	// プライマリレプリカの作成
	db, err := newDBClient(params, dopts)
	if err != nil {
		return nil, err
	}

	c := &Client{
		DB: db,
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
	tx, err := c.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit().Error
	}()
	err = f(tx)
	return
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

func newDBClient(params *Params, opts *options) (*gorm.DB, error) {
	conf := &gorm.Config{
		Logger:  zapgorm2.New(opts.logger),
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
