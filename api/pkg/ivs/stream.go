package ivs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivs"
)

type GetStreamParams struct {
	ChannelArn string
}

func (c *client) GetStream(ctx context.Context, params *GetStreamParams) (*ivs.GetStreamOutput, error) {
	in := &ivs.GetStreamInput{
		ChannelArn: &params.ChannelArn,
	}
	out, err := c.ivs.GetStream(ctx, in)
	if err != nil {
		return nil, c.streamError(err)
	}
	return out, nil
}
