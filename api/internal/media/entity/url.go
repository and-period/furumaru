package entity

import (
	"net/url"
	"strings"
)

type AdminURLMaker struct {
	url *url.URL
}

func NewAdminURLMaker(url *url.URL) *AdminURLMaker {
	return &AdminURLMaker{url: url}
}

func (m *AdminURLMaker) AuthYoutubeCallback() string {
	// e.g.) /auth/youtube/callback
	paths := []string{"auth", "youtube", "callback"}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")
	return webURL.String()
}
