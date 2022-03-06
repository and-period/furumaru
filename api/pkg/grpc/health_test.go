package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestHealthServer(t *testing.T) {
	tests := []struct {
		name   string
		server *grpc.Server
	}{
		{
			name:   "success",
			server: grpc.NewServer(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterHealthServer(tt.server)
			assert.NotNil(t, tt.server)
		})
	}
}
