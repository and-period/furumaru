package mysql

import (
	"crypto/tls"
	"fmt"
	"time"

	dmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewTiDBClient(params *Params, opts ...Option) (*Client, error) {
	options := &options{
		now:              time.Now,
		location:         time.UTC,
		maxAllowedPacket: 4194304, // 4MiB
		maxRetries:       3,
		maxConnLifetime:  5 * time.Minute,
		maxConnIdleTime:  5 * time.Minute,
	}
	for _, opt := range opts {
		opt(options)
	}

	db, err := newTiDBClient(params, options)
	if err != nil {
		return nil, err
	}
	sql, err := db.DB()
	if err != nil {
		return nil, err
	}
	sql.SetConnMaxLifetime(options.maxConnLifetime)
	sql.SetConnMaxIdleTime(options.maxConnIdleTime)

	c := &Client{
		DB:         db,
		maxRetries: uint64(options.maxRetries),
	}
	return c, nil
}

func newTiDBClient(params *Params, opts *options) (*gorm.DB, error) {
	tls := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: params.Host,
	}
	if err := dmysql.RegisterTLSConfig("tidb", tls); err != nil {
		return nil, err
	}
	conf := &gorm.Config{
		NowFunc: opts.now,
	}
	dsn := newTiDBDSN(params, opts)
	return gorm.Open(mysql.Open(dsn), conf)
}

func newTiDBDSN(params *Params, opts *options) string {
	dsn := &dmysql.Config{
		User:                 params.Username,
		Passwd:               params.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", params.Host, params.Port),
		DBName:               params.Database,
		Loc:                  opts.location,
		MaxAllowedPacket:     opts.maxAllowedPacket,
		ParseTime:            true,
		CheckConnLiveness:    true,
		AllowNativePasswords: true,
		TLSConfig:            "tidb",
	}
	return dsn.FormatDSN()
}
