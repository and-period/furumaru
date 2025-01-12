package common

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/zap"
)

type Client interface {
	Execute(ctx context.Context) error
}

type Params struct {
	Logger     *zap.Logger
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	SrcDir     string
}

func NewDBClient(p *Params, database string) (*mysql.Client, error) {
	params := &mysql.Params{
		Socket:   "tcp",
		Host:     p.DBHost,
		Port:     p.DBPort,
		Database: database,
		Username: p.DBUsername,
		Password: p.DBPassword,
	}
	opts := []mysql.Option{
		mysql.WithLocation(jst.Location()),
	}
	switch p.DBDriver {
	case "mysql":
		return mysql.NewClient(params, opts...)
	case "tidb":
		return mysql.NewTiDBClient(params, opts...)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", p.DBDriver)
	}
}
