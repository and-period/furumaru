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
		maxAllowedPacket: 4194304,        // 4MiB
		maxRetries:       3,
		maxOpenConns:     50,
		maxIdleConns:     50,
		maxConnLifetime:  5 * time.Minute,
		maxConnIdleTime:  5 * time.Minute,
		dialTimeout:      10 * time.Second,
		readTimeout:      30 * time.Second,
		writeTimeout:     30 * time.Second,
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
	sql.SetMaxOpenConns(options.maxOpenConns)
	sql.SetMaxIdleConns(options.maxIdleConns)
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
		InterpolateParams:    true,
		TLSConfig:            "tidb",
		Timeout:              opts.dialTimeout,
		ReadTimeout:          opts.readTimeout,
		WriteTimeout:         opts.writeTimeout,
	}
	return dsn.FormatDSN()
}
