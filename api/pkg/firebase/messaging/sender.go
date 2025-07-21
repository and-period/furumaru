package messaging

import (
	"context"

	"firebase.google.com/go/v4/messaging"
)

type Message struct {
	Title    string
	Body     string
	ImageURL string
	Data     map[string]string
}

func (c *client) Send(ctx context.Context, msg *Message, token string) error {
	message := &messaging.Message{
		Token: token,
		Data:  msg.Data,
		Notification: &messaging.Notification{
			Title:    msg.Title,
			Body:     msg.Body,
			ImageURL: msg.ImageURL,
		},
	}
	_, err := c.messageing.Send(ctx, message)
	return c.sendError(err)
}

func (c *client) MultiSend(
	ctx context.Context,
	msg *Message,
	tokens ...string,
) (int64, int64, error) {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Data:   msg.Data,
		Notification: &messaging.Notification{
			Title:    msg.Title,
			Body:     msg.Body,
			ImageURL: msg.ImageURL,
		},
	}
	res, err := c.messageing.SendEachForMulticast(ctx, message)
	if err != nil {
		return 0, 0, c.sendError(err)
	}
	return int64(res.SuccessCount), int64(res.FailureCount), nil
}
