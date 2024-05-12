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

func (m *AdminURLMaker) AuthYoutubeCallback(scheduleID string) string {
	// e.g.) /auth/youtube/callback
	paths := []string{"auth", "youtube", "callback"}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")

	query := webURL.Query()
	query.Set("schedule-id", scheduleID)

	webURL.RawQuery = query.Encode()
	return webURL.String()
}
