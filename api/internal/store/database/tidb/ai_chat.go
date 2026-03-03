package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const (
	aiChatSessionTable = "ai_chat_sessions"
	aiChatMessageTable = "ai_chat_messages"
)

type aiChat struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAiChat(db *mysql.Client) database.AiChat {
	return &aiChat{
		db:  db,
		now: jst.Now,
	}
}

func (a *aiChat) CreateSession(ctx context.Context, session *entity.AiChatSession) error {
	now := a.now()
	session.CreatedAt, session.UpdatedAt = now, now

	err := a.db.DB.WithContext(ctx).Table(aiChatSessionTable).Create(&session).Error
	return dbError(err)
}

func (a *aiChat) GetSession(ctx context.Context, sessionID string) (*entity.AiChatSession, error) {
	var session *entity.AiChatSession

	err := a.db.Statement(ctx, a.db.DB, aiChatSessionTable).
		Where("id = ?", sessionID).
		First(&session).Error
	return session, dbError(err)
}

func (a *aiChat) ListSessionsByAdminID(
	ctx context.Context, adminID string, limit, offset int,
) (entity.AiChatSessions, error) {
	var sessions entity.AiChatSessions

	stmt := a.db.Statement(ctx, a.db.DB, aiChatSessionTable).
		Where("admin_id = ?", adminID).
		Order("created_at DESC")

	if limit > 0 {
		stmt = stmt.Limit(limit)
	}
	if offset > 0 {
		stmt = stmt.Offset(offset)
	}

	err := stmt.Find(&sessions).Error
	return sessions, dbError(err)
}

func (a *aiChat) CreateMessage(ctx context.Context, message *entity.AiChatMessage) error {
	now := a.now()
	message.CreatedAt = now

	err := a.db.DB.WithContext(ctx).Table(aiChatMessageTable).Create(&message).Error
	return dbError(err)
}

func (a *aiChat) ListMessagesBySessionID(
	ctx context.Context, sessionID string,
) (entity.AiChatMessages, error) {
	var messages entity.AiChatMessages

	err := a.db.Statement(ctx, a.db.DB, aiChatMessageTable).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, dbError(err)
}
