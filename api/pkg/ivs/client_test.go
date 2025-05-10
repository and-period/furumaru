package ivs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestClient(t *testing.T) {
	t.Parallel()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoError(t, err)
	ivs := NewClient(cfg, &Params{
		RecordingConfigurationArn: "arn:aws:iam::123456789012:user/Development/product_1234/*",
	},
		WithMaxRetries(1),
		WithInterval(time.Millisecond),
		WithLogger(zap.NewNop()),
	)
	assert.NotNil(t, ivs)
}

func TestStreamError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "not error",
			err:    nil,
			expect: nil,
		},
		{
			name:   "invalid argument",
			err:    &types.ValidationException{Message: aws.String("some error")},
			expect: ErrInvalidArgument,
		},
		{
			name:   "acccess denied",
			err:    &types.AccessDeniedException{Message: aws.String("some error")},
			expect: ErrForbidden,
		},
		{
			name:   "not found",
			err:    &types.ResourceNotFoundException{Message: aws.String("some error")},
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    &types.ConflictException{Message: aws.String("some error")},
			expect: ErrAlreadyExists,
		},
		{
			name:   "resource exhausted",
			err:    &types.ServiceQuotaExceededException{Message: aws.String("some error")},
			expect: ErrResourceExhausted,
		},
		{
			name:   "resource exhausted",
			err:    &types.ThrottlingException{Message: aws.String("some error")},
			expect: ErrResourceExhausted,
		},
		{
			name:   "internal",
			err:    &types.InternalServerException{Message: aws.String("some error")},
			expect: ErrInternal,
		},
		{
			name:   "canceled",
			err:    context.Canceled,
			expect: ErrCanceled,
		},
		{
			name:   "timeout",
			err:    context.DeadlineExceeded,
			expect: ErrTimeout,
		},
		{
			name:   "unknown",
			err:    errors.New("some error"),
			expect: ErrUnknown,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli := &client{logger: zap.NewNop()}
			err := cli.streamError(tt.err)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
