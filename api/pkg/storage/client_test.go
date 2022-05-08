package storage

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBucket(t *testing.T) {
	t.Parallel()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoError(t, err)
	bucket := NewBucket(cfg, &Params{},
		WithMaxRetries(1),
		WithInterval(time.Millisecond),
		WithLogger(zap.NewNop()),
	)
	assert.NotNil(t, bucket)
}
