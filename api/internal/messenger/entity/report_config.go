package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
)

// ReportTemplateID - レポートテンプレートID
type ReportTemplateID string

const (
	ReportTemplateIDReceivedContact           ReportTemplateID = "received-contact"            // お問い合わせ受領
	ReportTemplateIDNotification              ReportTemplateID = "notification"                // お知らせ投稿
	ReportTemplateIDOrderProductAuthorized    ReportTemplateID = "order-product-authorized"    // 支払い完了（商品）
	ReportTemplateIDOrderExperienceAuthorized ReportTemplateID = "order-experience-authorized" // 支払い完了（体験）
)

// ReportConfig - システムレポート送信設定
type ReportConfig struct {
	TemplateID  ReportTemplateID `json:"reportId"`    // レポートテンプレートID
	Overview    string           `json:"overview"`    // レポート概要
	Detail      string           `json:"detail"`      // レポート詳細
	Author      string           `json:"author"`      // 管理者名
	Link        string           `json:"link"`        // 詳細リンク
	PublishedAt time.Time        `json:"publishedAt"` // 公開日時
	ReceivedAt  time.Time        `json:"receivedAt"`  // 受信日時
}

func (c *ReportConfig) Fields() map[string]string {
	return map[string]string{
		"Overview":    c.Overview,
		"Detail":      c.Detail,
		"Author":      c.Author,
		"Link":        c.Link,
		"PublishedAt": jst.Format(c.PublishedAt, "2006-01-02 15:04:05"),
		"ReceivedAt":  jst.Format(c.ReceivedAt, "2006-01-02 15:04:05"),
	}
}
