package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
)

const (
	ReportIDReceivedContact = "received-contact" // お問い合わせ受領
	ReportIDNotification    = "notification"     // お知らせ投稿
)

// ReportConfig - システムレポート送信設定
type ReportConfig struct {
	ReportID    string    `json:"reportId"`    // レポートID
	Overview    string    `json:"overview"`    // レポート概要
	Detail      string    `json:"detail"`      // レポート詳細
	Author      string    `json:"author"`      // 作成者
	Link        string    `json:"link"`        // 詳細リンク
	PublishedAt time.Time `json:"publishedAt"` // 公開日時
	ReceivedAt  time.Time `json:"receivedAt"`  // 受信日時
}

func (r *ReportConfig) Fields() map[string]string {
	return map[string]string{
		"Overview":    r.Overview,
		"Detail":      r.Detail,
		"Author":      r.Author,
		"Link":        r.Link,
		"PublishedAt": jst.Format(r.PublishedAt, "2006-01-02 15:04:05"),
		"ReceivedAt":  jst.Format(r.ReceivedAt, "2006-01-02 15:04:05"),
	}
}
