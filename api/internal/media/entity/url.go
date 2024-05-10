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
	// e.g.) /schedules/:scheduleId/broadcasts/youtube
	paths := []string{"schedules", scheduleID, "broadcasts", "youtube"}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")
	return webURL.String()
}
