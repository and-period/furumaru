package database

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// const coordinatorTable = "coordinators"

// var coordinatorFields = []string{
// 	"id", "email", "phone_number",
// 	"lastname", "firstname", "lastname_kana", "firstname_kana",
// 	"company_name", "store_name", "thumbnail_url", "header_url",
// 	"twitter_account", "instagram_account", "facebook_account",
// 	"postal_code", "prefecture", "city", "address_line1", "address_line2",
// 	"created_at", "updated_at", "deleted_at",
// }

type coordinator struct {
	db  *database.Client
	now func() time.Time
}

func NewCoordinator(db *database.Client) Coordinator {
	return &coordinator{
		db:  db,
		now: jst.Now,
	}
}
