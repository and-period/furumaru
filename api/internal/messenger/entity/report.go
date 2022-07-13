package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
)

const (
	ReportIDReceivedContact = "received-contact" // お問い合わせ受領
)

// Report - システムレポート情報
type Report struct {
	ReportID   string    `gorm:"" json:"reportId"`   // レポートID
	Overview   string    `gorm:"" json:"overview"`   // レポート概要
	Link       string    `gorm:"" json:"link"`       // 詳細リンク
	ReceivedAt time.Time `gorm:"" json:"receivedAt"` // 受領日時
}

func (r *Report) Fields() map[string]string {
	return map[string]string{
		"Overview":   r.Overview,
		"Link":       r.Link,
		"ReceivedAt": jst.Format(r.ReceivedAt, "2006-01-02 15:04:05"),
	}
}
