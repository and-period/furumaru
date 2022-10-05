package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const userTable = "users"

var userFields = []string{
	"id", "registered", "created_at", "updated_at",
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

func (u *user) MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.Users, error) {
	var users entity.Users
	if len(fields) == 0 {
		fields = userFields
	}

	stmt := u.db.DB.WithContext(ctx).
		Table(userTable).Select(fields).
		Where("id IN (?)", userIDs)

	if err := stmt.Find(&users).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := u.fill(ctx, u.db.DB, users...); err != nil {
		return nil, exception.InternalError(err)
	}
	return users, nil
}

func (u *user) Get(ctx context.Context, userID string, fields ...string) (*entity.User, error) {
	user, err := u.get(ctx, u.db.DB, userID, fields...)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if err := u.fill(ctx, u.db.DB, user); err != nil {
		return nil, exception.InternalError(err)
	}
	return user, nil
}

func (u *user) Create(ctx context.Context, user *entity.User) error {
	_, err := u.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := u.now()
		user.CreatedAt, user.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(userTable).Create(&user).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (u *user) get(ctx context.Context, tx *gorm.DB, userID string, fields ...string) (*entity.User, error) {
	var user *entity.User
	if len(fields) == 0 {
		fields = userFields
	}

	err := tx.WithContext(ctx).
		Table(userTable).Select(fields).
		Where("id = ?", userID).
		First(&user).Error
	return user, err
}

func (u *user) fill(ctx context.Context, tx *gorm.DB, users ...*entity.User) error {
	ids := entity.Users(users).IDs()
	if len(ids) == 0 {
		return nil
	}

	customers, err := u.fetchCustomers(ctx, tx, ids)
	if err != nil {
		return err
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

	customerMap := customers.Map()
	memberMap := members.Map()
	guestMap := guests.Map()

	for _, u := range users {
		customer, ok := customerMap[u.ID]
		if !ok {
			customer = &entity.Customer{}
		}
		member, ok := memberMap[u.ID]
		if !ok {
			member = &entity.Member{}
		}
		guest, ok := guestMap[u.ID]
		if !ok {
			guest = &entity.Guest{}
		}

		u.Fill(customer, member, guest)
	}
	return nil
}

func (u *user) fetchCustomers(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.Customers, error) {
	var customers entity.Customers
	err := tx.WithContext(ctx).
		Table(customerTable).Select(customerFields).
		Where("user_id IN (?)", userIDs).
		Find(&customers).Error
	return customers, err
}

func (u *user) fetchMembers(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.Members, error) {
	var members entity.Members
	err := tx.WithContext(ctx).
		Table(memberTable).Select(memberFields).
		Where("user_id IN (?)", userIDs).
		Find(&members).Error
	return members, err
}

func (u *user) fetchGuests(ctx context.Context, tx *gorm.DB, userIDs []string) (entity.Guests, error) {
	var guests entity.Guests
	err := tx.WithContext(ctx).
		Table(guestTable).Select(guestFields).
		Where("user_id IN (?)", userIDs).
		Find(&guests).Error
	return guests, err
}
