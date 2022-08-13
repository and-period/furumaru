package messaging

import (
	"context"

	"firebase.google.com/go/v4/messaging"
)

type Notification struct {
	Title    string
	Body     string
	ImageURL string
}

func (c *client) Send(ctx context.Context, n *Notification, token string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    n.Title,
			Body:     n.Body,
			ImageURL: n.ImageURL,
		},
	}
	sendFn := func() error {
		_, err := c.messageing.Send(ctx, message)
		return err
	}
	err := c.do(ctx, sendFn)
	return c.sendError(err)
}

func (c *client) MultiSend(ctx context.Context, n *Notification, tokens ...string) (int64, int64, error) {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title:    n.Title,
			Body:     n.Body,
			ImageURL: n.ImageURL,
		},
	}
	var success, failure int64
	sendFn := func() error {
		res, err := c.messageing.SendMulticast(ctx, message)
		if err != nil {
			return err
		}
		success, failure = int64(res.SuccessCount), int64(res.FailureCount)
		return nil
	}
	if err := c.do(ctx, sendFn); err != nil {
		return 0, 0, c.sendError(err)
	}
	return success, failure, nil
}
