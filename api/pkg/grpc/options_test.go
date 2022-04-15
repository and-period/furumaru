package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGrpcServerOptions(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewGRPCOptions(WithLogger(zap.NewNop())))
}
