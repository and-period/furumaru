package worker

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (w *worker) createMessages(ctx context.Context, payload *entity.WorkerPayload) error {
	template, err := w.db.MessageTemplate.Get(ctx, payload.Message.TemplateID)
	if err != nil {
		return err
	}
	title, body, err := template.Build(payload.Message.Fields())
	if err != nil {
		return err
	}
	params := &entity.NewMessagesParams{
		UserType:   payload.UserType,
		UserIDs:    payload.UserIDs,
		Type:       payload.Message.MessageType,
		Title:      title,
		Body:       body,
		Link:       payload.Message.Link,
		ReceivedAt: payload.Message.ReceivedAt,
	}
	messages := entity.NewMessages(params)
	return w.db.Message.MultiCreate(ctx, messages)
}
