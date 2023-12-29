package entity

import "time"

// MessageTemplateID - メッセージテンプレートID
type MessageTemplateID string

const (
	MessageTemplateIDNotificationSystem    MessageTemplateID = "notification-system"    // お知らせ送信（システム関連）
	MessageTemplateIDNotificationLive      MessageTemplateID = "notification-live"      // お知らせ送信（ライブ関連）
	MessageTemplateIDNotificationPromotion MessageTemplateID = "notification-promotion" // お知らせ送信（セール関連）
	MessageTemplateIDNotificationOther     MessageTemplateID = "notification-other"     // お知らせ送信（その他）
)

// MessageConfig - メッセージ作成設定
type MessageConfig struct {
	TemplateID  MessageTemplateID `json:"templateId"`  // メッセージテンプレートID
	MessageType MessageType       `json:"messageType"` // メッセージ種別
	Title       string            `json:"title"`       // メッセージ件名
	Detail      string            `json:"detail"`      // メッセージ詳細
	Author      string            `json:"author"`      // メッセージ作成者
	Link        string            `json:"link"`        // 遷移先リンク
	ReceivedAt  time.Time         `json:"receivedAt"`  // 受信日時
}

func (c *MessageConfig) Fields() map[string]string {
	return map[string]string{
		"Title":  c.Title,
		"Detail": c.Detail,
		"Author": c.Author,
	}
}
