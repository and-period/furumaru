package entity

import (
	"bytes"
	"io"
	"text/template"
	"time"
)

// MessageTemplate - メッセージテンプレート
type MessageTemplate struct {
	TemplateID    MessageTemplateID `gorm:"primaryKey;column:id;<-:create"` // テンプレートID
	TitleTemplate string            `gorm:""`                               // テンプレート(件名)
	BodyTemplate  string            `gorm:""`                               // テンプレート(内容)
	CreatedAt     time.Time         `gorm:"<-:create"`                      // 登録日時
	UpdatedAt     time.Time         `gorm:""`                               // 更新日時
}

func (t *MessageTemplate) Build(fields map[string]string) (string, string, error) {
	title, err := t.build(t.TitleTemplate, fields)
	if err != nil {
		return "", "", err
	}
	body, err := t.build(t.BodyTemplate, fields)
	if err != nil {
		return "", "", err
	}
	return title, body, nil
}

func (t *MessageTemplate) build(tmpl string, fields map[string]string) (string, error) {
	text := template.Must(template.New("message").Parse(tmpl))
	var buf bytes.Buffer
	if err := text.Execute(io.Writer(&buf), fields); err != nil {
		return "", err
	}
	return buf.String(), nil
}
