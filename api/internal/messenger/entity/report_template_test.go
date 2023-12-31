package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
)

func TestReportTemplate_Build(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		template *ReportTemplate
		fields   map[string]string
		expect   linebot.FlexContainer
		hasErr   bool
	}{
		{
			name: "success",
			template: &ReportTemplate{
				TemplateID: ReportTemplateIDReceivedContact,
				Template:   tmpl,
				CreatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
				UpdatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
			},
			fields: map[string]string{
				"Overview": "レポートの概要です。",
				"Detail":   "レポートの詳細です。",
				"Link":     "https://and-period.jp",
			},
			expect: &linebot.BubbleContainer{
				Type: linebot.FlexContainerTypeBubble,
				Body: &linebot.BoxComponent{
					Type: linebot.FlexComponentTypeBox,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type: linebot.FlexComponentTypeText,
							Text: "レポートの概要です。",
						},
						&linebot.TextComponent{
							Type: linebot.FlexComponentTypeText,
							Text: "レポートの詳細です。",
						},
						&linebot.TextComponent{
							Type: linebot.FlexComponentTypeText,
							Text: "https://and-period.jp",
						},
					},
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.template.Build(tt.fields)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

var tmpl = `{
  "type": "bubble",
  "body": {
    "type": "box",
    "contents": [
      {"type": "text", "text": "{{.Overview}}"},
      {"type": "text", "text": "{{.Detail}}"},
      {"type": "text", "text": "{{.Link}}"}
    ]
  }
}`
