package common

import (
	"context"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

type Client interface {
	Execute(ctx context.Context) error
}

type Params struct {
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
	return mysql.NewTiDBClient(params, opts...)
}
