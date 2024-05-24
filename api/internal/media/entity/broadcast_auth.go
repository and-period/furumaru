package entity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

var errInvalidBroadcastToken = errors.New("entity: invalid broadcast token")

// BroadcastAuthType - ライブ配信連携認証種別
type BroadcastAuthType string

const (
	BroadcastAuthTypeYoutube BroadcastAuthType = "youtube" // Youtube認証
)

// BroadcastAuth - ライブ配信連携認証情報
type BroadcastAuth struct {
	SessionID  string            `dynamodbav:"session_id"`          // セッションID
	Type       BroadcastAuthType `dynamodbav:"type"`                // 認証種別
	Account    string            `dynamodbav:"account"`             // アカウント
	ScheduleID string            `dynamodbav:"schedule_id"`         // スケジュールID
	Token      []byte            `dynamodbav:"token"`               // 認証トークン
	ExpiredAt  time.Time         `dynamodbav:"expired_at,unixtime"` // 有効期限
	CreatedAt  time.Time         `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt  time.Time         `dynamodbav:"updated_at"`          // 更新日時
}

type BroadcastAuthParams struct {
	SessionID  string
	Account    string
	ScheduleID string
	Now        time.Time
	TTL        time.Duration
}

func NewYoutubeBroadcastAuth(params *BroadcastAuthParams) *BroadcastAuth {
	return &BroadcastAuth{
		SessionID:  params.SessionID,
		Type:       BroadcastAuthTypeYoutube,
		Account:    params.Account,
		ScheduleID: params.ScheduleID,
		ExpiredAt:  params.Now.Add(params.TTL),
		CreatedAt:  params.Now,
		UpdatedAt:  params.Now,
	}
}

func (a *BroadcastAuth) TableName() string {
	return "broadcast-auth"
}

func (a *BroadcastAuth) PrimaryKey() map[string]interface{} {
	return map[string]interface{}{
		"session_id": a.SessionID,
	}
}

func (a *BroadcastAuth) SetToken(token *oauth2.Token) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(token); err != nil {
		return err
	}
	a.Token = buf.Bytes()
	return nil
}

func (a *BroadcastAuth) GetToken() (*oauth2.Token, error) {
	if a == nil || a.Token == nil {
		return nil, fmt.Errorf("entity: broadcast auth token is empty: %w", errInvalidBroadcastToken)
	}
	token := &oauth2.Token{}
	if err := json.Unmarshal(a.Token, token); err != nil {
		return nil, err
	}
	if !token.Valid() {
		return nil, fmt.Errorf("entity: broadcast auth token is invalid: %w", errInvalidBroadcastToken)
	}
	return token, nil
}

func (a *BroadcastAuth) ValidYoutubeAuth(channels []*youtube.Channel) bool {
	if a == nil || a.Type != BroadcastAuthTypeYoutube {
		return false
	}
	for _, channel := range channels {
		if channel.Snippet.CustomUrl == a.Account {
			return true
		}
	}
	return false
}
