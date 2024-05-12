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
	res := maker.AuthYoutubeCallback()
	assert.Equal(t, "http://example.com/auth/youtube/callback", res)
}
