package entity

import "time"

// BroadcastAuthType - ライブ配信連携認証種別
type BroadcastAuthType string

const (
	BroadcastAuthTypeYouTube BroadcastAuthType = "youtube" // YouTube認証
)

// BroadcastAuth - ライブ配信連携認証情報
type BroadcastAuth struct {
	SessionID  string            `dynamodbav:"session_id"`          // セッションID
	Type       BroadcastAuthType `dynamodbav:"type"`                // 認証種別
	Account    string            `dynamodbav:"account"`             // アカウント
	ScheduleID string            `dynamodbav:"schedule_id"`         // スケジュールID
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

func NewYouTubeBroadcastAuth(params *BroadcastAuthParams) *BroadcastAuth {
	return &BroadcastAuth{
		SessionID:  params.SessionID,
		Type:       BroadcastAuthTypeYouTube,
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

func (a *BroadcastAuth) ValidYouTubeAuth(email string) bool {
	if a == nil {
		return false
	}
	return a.Type == BroadcastAuthTypeYouTube && a.Account == email
}
