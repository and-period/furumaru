package entity

import (
	"net/url"
	"strings"
)

/**
 * --------------------------
 * 管理者関連URL生成用
 * --------------------------
 */
type AdminURLMaker struct {
	url *url.URL
}

func NewAdminURLMaker(url *url.URL) *AdminURLMaker {
	return &AdminURLMaker{url: url}
}

// SignIn - サインイン
func (m *AdminURLMaker) SignIn() string {
	// e.g.) /signin
	paths := []string{"signin"}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")
	return webURL.String()
}

// Contact - お問い合わせ詳細
func (m *AdminURLMaker) Contact(contactID string) string {
	// e.g.) /contacts/:contact-id
	paths := []string{"contacts", contactID}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")
	return webURL.String()
}

// Notification - お知らせ詳細
func (m *AdminURLMaker) Notification(notificationID string) string {
	// e.g.) /notifications/:notification-id
	paths := []string{"notifications", notificationID}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")
	return webURL.String()
}

/**
 * --------------------------
 * 購入者関連URL生成用
 * --------------------------
 */
type UserURLMaker struct {
	url *url.URL
}

func NewUserURLMaker(url *url.URL) *UserURLMaker {
	return &UserURLMaker{url: url}
}

// SignIn - サインイン
func (m *UserURLMaker) SignIn() string {
	// e.g.) /signin
	paths := []string{"signin"}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")
	return webURL.String()
}

func (m *UserURLMaker) Live(scheduleID string) string {
	// e.g.) /live/:live-id
	paths := []string{"live", scheduleID}
	webURL := *m.url // copy
	webURL.Path = strings.Join(paths, "/")
	return webURL.String()
}
