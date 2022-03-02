package database

import (
	"database/sql"
	"fmt"
	"strings"

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
	Socket     string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string
	TimeZone   string
	EnabledTLS bool
	Logger     *zap.Logger
}

// NewClient - DBクライアントの構造体
func NewClient(params *Params) (*Client, error) {
	con := getConfig(params)

	// プライマリレプリカの作成
	db, err := getDBClient(con, params)
	if err != nil {
		return nil, err
	}

	c := &Client{
		DB: db,
	}

	return c, nil
}

// Begin - トランザクションの開始処理
func (c *Client) Begin(opts ...*sql.TxOptions) (*gorm.DB, error) {
	tx := c.DB.Begin()
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
func (c *Client) Transaction(f func(tx *gorm.DB) (interface{}, error)) (data interface{}, err error) {
	tx, err := c.Begin()
	if err != nil {
		return nil, err
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
	data, err = f(tx)
	return
}

func getDBClient(config string, params *Params) (*gorm.DB, error) {
	if params.Logger == nil {
		params.Logger = zap.NewNop()
	}

	opt := &gorm.Config{
		Logger: zapgorm2.New(params.Logger),
	}

	return gorm.Open(mysql.Open(config), opt)
}

func getConfig(params *Params) string {
	switch params.Socket {
	case "tcp":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True%s%s",
			params.Username,
			params.Password,
			params.Host,
			params.Port,
			params.Database,
			withTLS(params.EnabledTLS),
			withTimeZone(params.TimeZone),
		)
	case "unix":
		return fmt.Sprintf(
			"%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=true%s",
			params.Username,
			params.Password,
			params.Host,
			params.Database,
			withTLS(params.EnabledTLS),
		)
	default:
		return ""
	}
}

func withTLS(enabled bool) string {
	if !enabled {
		return ""
	}
	return "&tls=true"
}

func withTimeZone(tz string) string {
	if tz == "" {
		tz = "Asia%2FTokyo"
	} else {
		tz = strings.Replace(tz, "/", "%2F", -1)
	}
	return fmt.Sprintf("&loc=%s", tz)
}
