package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
)

const (
	EmailIDAdminRegister       = "admin-register"        // 管理者登録
	EmailIDUserReceivedContact = "user-received-contact" // お問い合わせ受領
)

// MailConfig - メール送信設定
type MailConfig struct {
	EmailID       string            `json:"emailId"`       // メールテンプレートID
	Substitutions map[string]string `json:"substitutions"` // メール動的内容
}

type TemplateDataBuilder struct {
	data map[string]string
}

func NewTemplateDataBuilder() *TemplateDataBuilder {
	return &TemplateDataBuilder{
		data: map[string]string{},
	}
}

func (b *TemplateDataBuilder) Build() map[string]string {
	return b.data
}

func (b *TemplateDataBuilder) Data(data map[string]string) *TemplateDataBuilder {
	if data != nil {
		b.data = data
	}
	return b
}

func (b *TemplateDataBuilder) YearMonth(yearMonth time.Time) *TemplateDataBuilder {
	b.data["年月"] = jst.Format(yearMonth, "2006年01月")
	return b
}

func (b *TemplateDataBuilder) Name(name string) *TemplateDataBuilder {
	b.data["氏名"] = name
	return b
}

func (b *TemplateDataBuilder) Email(email string) *TemplateDataBuilder {
	b.data["メールアドレス"] = email
	return b
}

func (b *TemplateDataBuilder) Password(password string) *TemplateDataBuilder {
	b.data["パスワード"] = password
	return b
}

func (b *TemplateDataBuilder) WebURL(url string) *TemplateDataBuilder {
	b.data["サイトURL"] = url
	return b
}

func (b *TemplateDataBuilder) Contact(title, body string) *TemplateDataBuilder {
	b.data["件名"] = title
	b.data["本文"] = body
	return b
}
