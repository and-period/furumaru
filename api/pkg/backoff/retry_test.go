package backoff

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	backoff := NewFixedIntervalBackoff(time.Second, 2)
	testFn := func() error { return errors.New("some error") }
	retryable := WithRetryablel(func(err error) bool { return true })
	err := Retry(ctx, backoff, testFn, retryable)
	assert.Error(t, err)
}
