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
