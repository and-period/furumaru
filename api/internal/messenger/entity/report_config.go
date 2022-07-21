package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
)

const (
	ReportIDReceivedContact = "received-contact" // お問い合わせ受領
)

// ReportConfig - システムレポート送信設定
type ReportConfig struct {
	ReportID   string    `json:"reportId"`   // レポートID
	Overview   string    `json:"overview"`   // レポート概要
	Detail     string    `json:"detail"`     // レポート詳細
	Link       string    `json:"link"`       // 詳細リンク
	ReceivedAt time.Time `json:"receivedAt"` // 受信日時
}

func (r *ReportConfig) Fields() map[string]string {
	return map[string]string{
		"Overview":   r.Overview,
		"Detail":     r.Detail,
		"Link":       r.Link,
		"ReceivedAt": jst.Format(r.ReceivedAt, "2006-01-02 15:04:05"),
	}
}
