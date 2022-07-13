package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

var tmpl = `レポート概要: {{.Overview}}
レポートリンク: {{.Link}}`

func TestReportTemplate_Build(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		template *ReportTemplate
		fields   map[string]string
		expect   string
		hasErr   bool
	}{
		{
			name: "success",
			template: &ReportTemplate{
				TemplateID: ReportIDReceivedContact,
				Template:   tmpl,
				CreatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
				UpdatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
			},
			fields: map[string]string{
				"Overview": "レポートの概要です。",
				"Link":     "https://and-period.jp",
			},
			expect: "レポート概要: レポートの概要です。\nレポートリンク: https://and-period.jp",
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
