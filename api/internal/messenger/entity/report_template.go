package entity

import (
	"bytes"
	"io"
	"text/template"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// ReportTemplate - システムレポートテンプレート
type ReportTemplate struct {
	TemplateID ReportTemplateID `gorm:"primaryKey;column:id;<-:create"` // テンプレートID
	Template   string           `gorm:""`                               // テンプレート
	CreatedAt  time.Time        `gorm:"<-:create"`                      // 登録日時
	UpdatedAt  time.Time        `gorm:""`                               // 更新日時
}

func (t *ReportTemplate) Build(fields map[string]string) (linebot.FlexContainer, error) {
	text := template.Must(template.New("report").Parse(t.Template))
	var buf bytes.Buffer
	if err := text.Execute(io.Writer(&buf), fields); err != nil {
		return nil, err
	}
	return linebot.UnmarshalFlexMessageJSON(buf.Bytes())
}
