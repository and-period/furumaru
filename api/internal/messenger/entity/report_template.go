package entity

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// ReportTemplate - システムレポートテンプレート
type ReportTemplate struct {
	TemplateID ReportTemplateID `gorm:"primaryKey;column:id;<-:create"` // テンプレートID
	Template   string           `gorm:""`                               // テンプレート
	CreatedAt  time.Time        `gorm:"<-:create"`                      // 登録日時
	UpdatedAt  time.Time        `gorm:""`                               // 更新日時
}

func (t *ReportTemplate) Build(fields map[string]string) (messaging_api.FlexContainerInterface, error) {
	text := template.Must(template.New("report").Parse(t.Template))
	var buf bytes.Buffer
	if err := text.Execute(io.Writer(&buf), fields); err != nil {
		return nil, fmt.Errorf("entity: failed to execute report template: %w", err)
	}
	return messaging_api.UnmarshalFlexContainer(buf.Bytes())
}
