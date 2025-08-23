package tidb

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

func NewUser(db *mysql.Client) database.User {
	return &user{
		db:  db,
		now: jst.Now,
	}
}

type listUsersParams database.ListUsersParams

func (p listUsersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.WithDeleted {
		stmt = stmt.Unscoped()
	}
	if p.OnlyRegistered {
		stmt = stmt.Where("registered = ?", true)
	}
	if p.OnlyVerified {
		stmt = stmt.Joins("INNER JOIN members ON members.user_id = users.id").
			Where("members.verified_at IS NOT NULL").
			Order("users.updated_at DESC")
	} else {
		stmt = stmt.Order("updated_at DESC")
	}
	return stmt
}

func (p listUsersParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (u *user) List(ctx context.Context, params *database.ListUsersParams, fields ...string) (entity.Users, error) {
	var users entity.Users

	p := listUsersParams(*params)

	stmt := u.db.Statement(ctx, u.db.DB, userTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&users).Error; err != nil {
		return nil, dbError(err)
	}
	if err := u.fill(ctx, u.db.DB, users...); err != nil {
		return nil, dbError(err)
	}
	return users, nil
}

func (u *user) Count(ctx context.Context, params *database.ListUsersParams) (int64, error) {
	p := listUsersParams(*params)

	total, err := u.db.Count(ctx, u.db.DB, &entity.User{}, p.stmt)
	return total, dbError(err)
}

func (u *user) MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.Users, error) {
	var users entity.Users

	stmt := u.db.Statement(ctx, u.db.DB, userTable, fields...).Unscoped().Where("id IN (?)", userIDs)

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
		Unscoped().
		Where("id = ?", userID).
		First(&user).Error
	return user, err
}

func (u *user) fill(ctx context.Context, tx *gorm.DB, users ...*entity.User) error {
	ids := entity.Users(users).IDs()
	if len(ids) == 0 {
		return nil
	}

	usersMap := entity.Users(users).GroupByUserType()

	var (
		members       entity.Members
		guests        entity.Guests
		facilityUsers entity.FacilityUsers
	)
	for userType, us := range usersMap {
		var err error
		switch userType {
		case entity.UserTypeMember:
			members, err = u.fetchMembers(ctx, tx, us.IDs())
		case entity.UserTypeGuest:
			guests, err = u.fetchGuests(ctx, tx, us.IDs())
		case entity.UserTypeFacilityUser:
			facilityUsers, err = u.fetchFacilityUsers(ctx, tx, us.IDs())
		}
		if err != nil {
			return err
		}
	}
	entity.Users(users).Fill(members.Map(), guests.Map(), facilityUsers.Map())
	return nil
}

func (u *user) fetchMembers(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.Members, error) {
	var members entity.Members

	stmt := u.db.Statement(ctx, tx, memberTable).Where("user_id IN (?)", userIDs)

	err := stmt.Find(&members).Error
	return members, err
}

func (u *user) fetchGuests(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.Guests, error) {
	var guests entity.Guests

	err := u.db.Statement(ctx, tx, guestTable).Where("user_id IN (?)", userIDs).Find(&guests).Error
	return guests, err
}

func (u *user) fetchFacilityUsers(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.FacilityUsers, error) {
	var facilityUsers entity.FacilityUsers

	err := u.db.Statement(ctx, tx, facilityUserTable).Where("user_id IN (?)", userIDs).Find(&facilityUsers).Error
	return facilityUsers, err
}
