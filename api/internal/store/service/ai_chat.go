package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) CreateAiChatSession(
	ctx context.Context, in *store.CreateAiChatSessionInput,
) (*entity.AiChatSession, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewAiChatSessionParams{
		AdminID:   in.AdminID,
		ProductID: in.ProductID,
		Title:     in.Title,
	}
	session := entity.NewAiChatSession(params)
	if err := s.db.AiChat.CreateSession(ctx, session); err != nil {
		return nil, internalError(err)
	}
	return session, nil
}

func (s *service) GetAiChatSession(
	ctx context.Context, in *store.GetAiChatSessionInput,
) (*entity.AiChatSession, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	session, err := s.db.AiChat.GetSession(ctx, in.SessionID)
	if err != nil {
		return nil, internalError(err)
	}
	return session, nil
}

func (s *service) ListAiChatSessions(
	ctx context.Context, in *store.ListAiChatSessionsInput,
) (entity.AiChatSessions, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	sessions, err := s.db.AiChat.ListSessionsByAdminID(ctx, in.AdminID, int(in.Limit), int(in.Offset))
	if err != nil {
		return nil, internalError(err)
	}
	return sessions, nil
}

func (s *service) CreateAiChatMessage(
	ctx context.Context, in *store.CreateAiChatMessageInput,
) (*entity.AiChatMessage, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// セッションの存在確認
	if _, err := s.db.AiChat.GetSession(ctx, in.SessionID); err != nil {
		return nil, fmt.Errorf("service: session not found: %w", exception.ErrNotFound)
	}
	params := &entity.NewAiChatMessageParams{
		SessionID: in.SessionID,
		Role:      in.Role,
		Content:   in.Content,
	}
	message := entity.NewAiChatMessage(params)
	if err := s.db.AiChat.CreateMessage(ctx, message); err != nil {
		return nil, internalError(err)
	}
	return message, nil
}

func (s *service) ListAiChatMessages(
	ctx context.Context, in *store.ListAiChatMessagesInput,
) (entity.AiChatMessages, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	messages, err := s.db.AiChat.ListMessagesBySessionID(ctx, in.SessionID)
	if err != nil {
		return nil, internalError(err)
	}
	return messages, nil
}
