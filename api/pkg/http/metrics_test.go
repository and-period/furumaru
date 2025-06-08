package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsServer(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	tests := []struct {
		name   string
		port   int64
		expect int
	}{
		{
			name:   "success",
			port:   20081,
			expect: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewMetricsServer(tt.port)
			go server.Serve()
			defer func() {
				server.Stop(ctx)
				time.Sleep(2 * time.Second)
			}()

			for range 10 {
				time.Sleep(time.Microsecond * 100)
				url := fmt.Sprintf("http://localhost:%d/metrics", tt.port)
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				require.NoError(t, err)
				res, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				defer res.Body.Close()
				if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
					continue
				}
				assert.Equal(t, tt.expect, res.StatusCode)
				break
			}
		})
	}
}
