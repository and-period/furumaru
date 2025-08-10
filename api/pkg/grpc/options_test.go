package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrpcServerOptions(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewGRPCOptions())
}
