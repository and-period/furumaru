package entity

import (
	"bytes"
	"io"
	"text/template"
	"time"
)

// MessageTemplate - メッセージテンプレート
type MessageTemplate struct {
	TemplateID string    `gorm:"primaryKey;column:id;<-:create"` // メッセージID
	Template   string    `gorm:""`                               // テンプレート
	CreatedAt  time.Time `gorm:"<-:create"`                      // 登録日時
	UpdatedAt  time.Time `gorm:""`                               // 更新日時
}

func (t *MessageTemplate) Build(fields map[string]string) (string, error) {
	text := template.Must(template.New("report").Parse(t.Template))
	var buf bytes.Buffer
	if err := text.Execute(io.Writer(&buf), fields); err != nil {
		return "", err
	}
	return buf.String(), nil
}
