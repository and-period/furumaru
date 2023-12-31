package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPushTemplate_Build(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		template    *PushTemplate
		fields      map[string]string
		expectTitle string
		expectBody  string
		hasErr      bool
	}{
		{
			name: "success",
			template: &PushTemplate{
				TemplateID:    PushTemplateIDContact,
				TitleTemplate: "件名: {{.Title}}",
				BodyTemplate:  "内容: {{.Body}}",
				ImageURL:      "https://and-period.jp/image.png",
				CreatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
				UpdatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
			},
			fields: map[string]string{
				"Title": "テスト",
				"Body":  "テスト",
			},
			expectTitle: "件名: テスト",
			expectBody:  "内容: テスト",
			hasErr:      false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			title, body, err := tt.template.Build(tt.fields)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expectTitle, title)
			assert.Equal(t, tt.expectBody, body)
		})
	}
}
