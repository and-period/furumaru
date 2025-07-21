package mediaconvert

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
)

func (c *client) CreateJob(
	ctx context.Context,
	template string,
	settings *types.JobSettings,
) error {
	in := &mediaconvert.CreateJobInput{
		Role:        c.role,
		Settings:    settings,
		JobTemplate: aws.String(template),
	}
	_, err := c.convert.CreateJob(ctx, in)
	return err
}
