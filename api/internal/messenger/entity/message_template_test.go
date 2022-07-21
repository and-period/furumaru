package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestMessageTemplate_Build(t *testing.T) {
	t.Parallel()
	var tmpl = `メッセージの内容です。`
	tests := []struct {
		name     string
		template *MessageTemplate
		fields   map[string]string
		expect   string
		hasErr   bool
	}{
		{
			name: "success",
			template: &MessageTemplate{
				TemplateID: MessageIDNotification,
				Template:   tmpl,
				CreatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
				UpdatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
			},
			fields: map[string]string{},
			expect: "メッセージの内容です。メッセージの内容です。t re",
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
