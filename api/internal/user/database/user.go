package database

import (
	"context"
	"time"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/jst"
)

const userTable = "users"

var userFields = []string{
	"id", "cognito_id", "provider_type", "email", "phone_number",
	"created_at", "updated_at", "verified_at", "deleted_at",
}

type user struct {
	db  *database.Client
	now func() time.Time
}

func NewUser(db *database.Client) User {
	return &user{
		db:  db,
		now: jst.Now,
	}
}

func (u *user) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.User, error) {
	var user *entity.User
	if len(fields) == 0 {
		fields = userFields
	}

	stmt := u.db.DB.
		Table(userTable).Select(fields).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&user).Error; err != nil {
		return nil, dbError(err)
	}
	return user, nil
}
