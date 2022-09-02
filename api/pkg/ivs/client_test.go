package ivs

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
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
