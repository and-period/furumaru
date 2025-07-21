package ivs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
)

type GetStreamParams struct {
	ChannelArn string
}

type GetStreamKeyParams struct {
	StreamKeyArn string
}

func (c *client) GetStream(ctx context.Context, params *GetStreamParams) (*types.Stream, error) {
	in := &ivs.GetStreamInput{
		ChannelArn: aws.String(params.ChannelArn),
	}
	out, err := c.ivs.GetStream(ctx, in)
	if err != nil {
		return nil, c.streamError(err)
	}
	return out.Stream, nil
}

func (c *client) GetStreamKey(
	ctx context.Context,
	params *GetStreamKeyParams,
) (*types.StreamKey, error) {
	in := &ivs.GetStreamKeyInput{
		Arn: aws.String(params.StreamKeyArn),
	}
	out, err := c.ivs.GetStreamKey(ctx, in)
	if err != nil {
		return nil, c.streamError(err)
	}
	return out.StreamKey, nil
}
