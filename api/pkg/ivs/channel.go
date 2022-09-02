package ivs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	ivs "github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
)

type CreateChannelParams struct {
	authorized                bool
	LatencyMode               types.ChannelLatencyMode
	Name                      string
	RecordingConfigurationArn string
	ChannelType               types.ChannelType
}

func (c *client) CreateChannel(ctx context.Context, params *CreateChannelParams) (*ivs.CreateChannelOutput, error) {
	in := &ivs.CreateChannelInput{
		Authorized:                params.authorized,
		LatencyMode:               params.LatencyMode,
		Name:                      aws.String(params.Name),
		RecordingConfigurationArn: aws.String(params.RecordingConfigurationArn),
		Type:                      params.ChannelType,
	}
	out, err := c.ivs.CreateChannel(ctx, in)
	if err != nil {
		return nil, c.streamError(err)
	}
	return out, nil
}
