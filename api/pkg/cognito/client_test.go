package cognito

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Parallel()
	cfg, err := config.LoadDefaultConfig(t.Context())
	require.NoError(t, err)
	auth := NewClient(cfg, &Params{},
		WithMaxRetries(1),
		WithInterval(time.Millisecond),
	)
	assert.NotNil(t, auth)
}

func TestAuthError(t *testing.T) {
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
			err:    &types.CodeMismatchException{Message: aws.String("some error")},
			expect: ErrInvalidArgument,
		},
		{
			name:   "unauthenticated",
			err:    &types.NotAuthorizedException{Message: aws.String("some error")},
			expect: ErrUnauthenticated,
		},
		{
			name:   "not found",
			err:    &types.ResourceNotFoundException{Message: aws.String("some error")},
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    &types.UsernameExistsException{Message: aws.String("some error")},
			expect: ErrAlreadyExists,
		},
		{
			name:   "resource exhausted",
			err:    &types.LimitExceededException{Message: aws.String("some error")},
			expect: ErrResourceExhausted,
		},
		{
			name:   "internal",
			err:    &types.InternalErrorException{Message: aws.String("some error")},
			expect: ErrInternal,
		},
		{
			name:   "code expired",
			err:    &types.ExpiredCodeException{Message: aws.String("some error")},
			expect: ErrCodeExpired,
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
			cli := &client{}
			err := cli.authError(tt.err)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
