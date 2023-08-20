package medialive

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/medialive"
)

func (c *client) StartChannel(ctx context.Context, channelID string) error {
	in := &medialive.StartChannelInput{
		ChannelId: aws.String(channelID),
	}
	_, err := c.media.StartChannel(ctx, in)
	return err
}

func (c *client) StopChannel(ctx context.Context, channelID string) error {
	in := &medialive.StopChannelInput{
		ChannelId: aws.String(channelID),
	}
	_, err := c.media.StopChannel(ctx, in)
	return err
}
