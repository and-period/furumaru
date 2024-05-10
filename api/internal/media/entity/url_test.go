package entity

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminURLMaker(t *testing.T) {
	t.Parallel()
	webURL, err := url.Parse("http://example.com")
	require.NoError(t, err)
	maker := NewAdminURLMaker(webURL)
	res := maker.AuthYoutubeCallback("schedule-id")
	assert.Equal(t, "http://example.com/schedules/schedule-id/broadcasts/youtube", res)
}
