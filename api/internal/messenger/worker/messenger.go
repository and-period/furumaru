package worker

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (w *worker) messenger(ctx context.Context, payload *entity.WorkerPayload) error {
	template, err := w.db.MessageTemplate.Get(ctx, payload.Message.MessageID)
	if err != nil {
		return err
	}
	body, err := template.Build(payload.Message.Fields())
	if err != nil {
		return err
	}
	params := &entity.NewMessagesParams{
		UserType:   payload.UserType,
		UserIDs:    payload.UserIDs,
		Type:       payload.Message.MessageType,
		Title:      payload.Message.Title,
		Body:       body,
		Link:       payload.Message.Link,
		ReceivedAt: payload.Message.ReceivedAt,
	}
	messages := entity.NewMessages(params)
	return w.db.Message.MultiCreate(ctx, messages)
}
