package tidb

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const messageTable = "messages"

type message struct {
	db  *mysql.Client
	now func() time.Time
}

func NewMessage(db *mysql.Client) database.Message {
	return &message{
		db:  db,
		now: jst.Now,
	}
}

type listMessagesParams database.ListMessagesParams

func (p *listMessagesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.UserType != entity.UserTypeNone {
		stmt = stmt.Where("user_type = ?", p.UserType)
	}
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("`%s` ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("`%s` DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

func (p *listMessagesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (m *message) List(
	ctx context.Context,
	params *database.ListMessagesParams,
	fields ...string,
) (entity.Messages, error) {
	var messages entity.Messages

	p := listMessagesParams(*params)

	stmt := m.db.Statement(ctx, m.db.DB, messageTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	err := stmt.Find(&messages).Error
	return messages, dbError(err)
}

func (m *message) Count(ctx context.Context, params *database.ListMessagesParams) (int64, error) {
	p := listMessagesParams(*params)

	total, err := m.db.Count(ctx, m.db.DB, &entity.Message{}, p.stmt)
	return total, dbError(err)
}

func (m *message) Get(
	ctx context.Context,
	messageID string,
	fields ...string,
) (*entity.Message, error) {
	message, err := m.get(ctx, m.db.DB, messageID, fields...)
	return message, dbError(err)
}

func (m *message) MultiCreate(ctx context.Context, messages entity.Messages) error {
	now := m.now()
	for i := range messages {
		messages[i].CreatedAt = now
		messages[i].UpdatedAt = now
	}

	err := m.db.DB.WithContext(ctx).Table(messageTable).Create(&messages).Error
	return dbError(err)
}

func (m *message) UpdateRead(ctx context.Context, messageID string) error {
	params := map[string]interface{}{
		"read":       true,
		"updated_at": m.now(),
	}
	stmt := m.db.DB.WithContext(ctx).
		Table(messageTable).
		Where("id = ?", messageID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (m *message) get(
	ctx context.Context,
	tx *gorm.DB,
	messageID string,
	fields ...string,
) (*entity.Message, error) {
	var message *entity.Message

	stmt := m.db.Statement(ctx, tx, messageTable, fields...).
		Where("id = ?", messageID)

	if err := stmt.First(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}
