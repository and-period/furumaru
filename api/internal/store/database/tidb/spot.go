package tidb

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const spotTable = "spots"

type spot struct {
	database.Spot
	db  *mysql.Client
	now func() time.Time
}

func newSpot(db *mysql.Client, mysql database.Spot) database.Spot {
	return &spot{
		Spot: mysql,
		db:   db,
		now:  jst.Now,
	}
}
