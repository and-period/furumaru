package medialive

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/medialive"
	"github.com/aws/aws-sdk-go-v2/service/medialive/types"
)

func (c *client) CreateSchedule(ctx context.Context, channelID string, actions ...types.ScheduleAction) error {
	in := &medialive.BatchUpdateScheduleInput{
		ChannelId: aws.String(channelID),
		Creates:   &types.BatchScheduleActionCreateRequest{ScheduleActions: actions},
	}
	_, err := c.media.BatchUpdateSchedule(ctx, in)
	return err
}
