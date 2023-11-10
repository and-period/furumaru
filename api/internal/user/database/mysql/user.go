package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const userTable = "users"

type user struct {
	db  *mysql.Client
	now func() time.Time
}

func newUser(db *mysql.Client) database.User {
	return &user{
		db:  db,
		now: jst.Now,
	}
}

func (u *user) List(ctx context.Context, params *database.ListUsersParams, fields ...string) (entity.Users, error) {
	var users entity.Users

	stmt := u.db.Statement(ctx, u.db.DB, userTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&users).Error; err != nil {
		return nil, dbError(err)
	}
	if err := u.fill(ctx, u.db.DB, users...); err != nil {
		return nil, dbError(err)
	}
	return users, nil
}

func (u *user) Count(ctx context.Context, _ *database.ListUsersParams) (int64, error) {
	total, err := u.db.Count(ctx, u.db.DB, &entity.User{}, nil)
	return total, dbError(err)
}

func (u *user) MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.Users, error) {
	var users entity.Users

	stmt := u.db.Statement(ctx, u.db.DB, userTable, fields...).
		Table(userTable).Select(fields).
		Where("id IN (?)", userIDs)

	if err := stmt.Find(&users).Error; err != nil {
		return nil, dbError(err)
	}
	if err := u.fill(ctx, u.db.DB, users...); err != nil {
		return nil, dbError(err)
	}
	return users, nil
}

func (u *user) Get(ctx context.Context, userID string, fields ...string) (*entity.User, error) {
	user, err := u.get(ctx, u.db.DB, userID, fields...)
	if err != nil {
		return nil, dbError(err)
	}
	if err := u.fill(ctx, u.db.DB, user); err != nil {
		return nil, dbError(err)
	}
	return user, nil
}

func (u *user) Create(ctx context.Context, user *entity.User) error {
	now := u.now()
	user.CreatedAt, user.UpdatedAt = now, now

	err := u.db.DB.WithContext(ctx).Table(userTable).Create(&user).Error
	return dbError(err)
}

func (u *user) get(ctx context.Context, tx *gorm.DB, userID string, fields ...string) (*entity.User, error) {
	var user *entity.User

	err := u.db.Statement(ctx, tx, userTable, fields...).
		Where("id = ?", userID).
		First(&user).Error
	return user, err
}

func (u *user) fill(ctx context.Context, tx *gorm.DB, users ...*entity.User) error {
	ids := entity.Users(users).IDs()
	if len(ids) == 0 {
		return nil
	}

	usersMap := entity.Users(users).GroupByRegistered()

	var (
		members entity.Members
		guests  entity.Guests
	)
	for registered, us := range usersMap {
		var err error
		if registered {
			members, err = u.fetchMembers(ctx, tx, us.IDs())
		} else {
			guests, err = u.fetchGuests(ctx, tx, us.IDs())
		}
		if err != nil {
			return err
		}
	}

	memberMap := members.Map()
	guestMap := guests.Map()

	for _, u := range users {
		member, ok := memberMap[u.ID]
		if !ok {
			member = &entity.Member{}
		}
		guest, ok := guestMap[u.ID]
		if !ok {
			guest = &entity.Guest{}
		}

		u.Fill(member, guest)
	}
	return nil
}

func (u *user) fetchMembers(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.Members, error) {
	var members entity.Members

	err := u.db.Statement(ctx, tx, memberTable).
		Where("user_id IN (?)", userIDs).
		Find(&members).Error
	return members, err
}

func (u *user) fetchGuests(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.Guests, error) {
	var guests entity.Guests

	err := u.db.Statement(ctx, tx, guestTable).
		Where("user_id IN (?)", userIDs).
		Find(&guests).Error
	return guests, err
}
