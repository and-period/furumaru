package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/stretchr/testify/assert"
)

func TestReportTemplate_Build(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		template *ReportTemplate
		fields   map[string]string
		expect   messaging_api.FlexContainerInterface
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
			expect: messaging_api.FlexBubble{
				FlexContainer: messaging_api.FlexContainer{Type: "bubble"},
				Body: &messaging_api.FlexBox{
					FlexComponent: messaging_api.FlexComponent{Type: "box"},
					Contents: []messaging_api.FlexComponentInterface{
						messaging_api.FlexText{
							FlexComponent: messaging_api.FlexComponent{Type: "text"},
							Text:          "レポートの概要です。",
						},
						messaging_api.FlexText{
							FlexComponent: messaging_api.FlexComponent{Type: "text"},
							Text:          "レポートの詳細です。",
						},
						messaging_api.FlexText{
							FlexComponent: messaging_api.FlexComponent{Type: "text"},
							Text:          "https://and-period.jp",
						},
					},
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
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
