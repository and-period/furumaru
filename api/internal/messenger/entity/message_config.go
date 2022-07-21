package entity

import "time"

const (
	MessageIDNotification = "notification" // お知らせ送信
)

// MessageConfig - メッセージ作成設定
type MessageConfig struct {
	MessageID   string      `json:"messageId"`   // メッセージID
	MessageType MessageType `json:"messageType"` // メッセージ種別
	Title       string      `json:"title"`       // メッセージ件名
	Prepared    string      `json:"prepared"`    // メッセージ作成者
	Link        string      `json:"link"`        // 遷移先リンク
	ReceivedAt  time.Time   `json:"receivedAt"`  // 受信日時
}

func (c *MessageConfig) Fields() map[string]string {
	return map[string]string{
		"Title":    c.Title,
		"Prepared": c.Prepared,
	}
}
