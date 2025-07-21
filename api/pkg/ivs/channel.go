package ivs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	ivs "github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
)

type CreateChannelParams struct {
	LatencyMode types.ChannelLatencyMode
	Name        string
	ChannelType types.ChannelType
}

type GetChannelParams struct {
	Arn string
}

type DeleteChannelParams struct {
	Arn string
}

func (c *client) CreateChannel(
	ctx context.Context,
	params *CreateChannelParams,
) (*ivs.CreateChannelOutput, error) {
	in := &ivs.CreateChannelInput{
		LatencyMode:               params.LatencyMode,
		Name:                      aws.String(params.Name),
		RecordingConfigurationArn: c.recordingConfigurationArn,
		Type:                      params.ChannelType,
	}
	out, err := c.ivs.CreateChannel(ctx, in)
	if err != nil {
		return nil, c.streamError(err)
	}
	return out, nil
}

func (c *client) GetChannel(ctx context.Context, params *GetChannelParams) (*types.Channel, error) {
	in := &ivs.GetChannelInput{
		Arn: aws.String(params.Arn),
	}

	out, err := c.ivs.GetChannel(ctx, in)
	if err != nil {
		return nil, c.streamError(err)
	}
	return out.Channel, nil
}

func (c *client) DeleteChannel(ctx context.Context, params *DeleteChannelParams) error {
	in := &ivs.DeleteChannelInput{
		Arn: aws.String(params.Arn),
	}
	_, err := c.ivs.DeleteChannel(ctx, in)
	return c.streamError(err)
}
