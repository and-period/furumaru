package mailer

import (
	"net/http"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

func TestSendGridError_Error(t *testing.T) {
	t.Parallel()
	err := &SendGridError{
		Code:    http.StatusBadRequest,
		Message: "test error",
		Field:   "field",
		Help:    "help message",
	}
	assert.Equal(t, "mailer: failed to SendGrid. code=400, field=test error", err.Error())
}

func TestMessage(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 2, 2, 18, 0, 0, 0)
	type args struct {
		templateID       string
		fromName         string
		fromAddress      string
		subject          string
		contents         []*Content
		personalizations []*Personalization
	}
	tests := []struct {
		name   string
		args   args
		expect *mail.SGMailV3
	}{
		{
			name: "success",
			args: args{
				templateID:  "template-id",
				fromName:    "test user",
				fromAddress: "test-user@and-period.jp",
				subject:     "テストタイトル",
				contents: []*Content{{
					ContentType: "type",
					Value:       "value",
				}},
				personalizations: []*Personalization{
					{
						Name:          "test to user",
						Address:       "test-to@and-period.jp",
						Type:          AddressTypeTo,
						Substitutions: map[string]interface{}{"key": "value"},
					},
					{
						Name:          "test cc user",
						Address:       "test-cc@and-period.jp",
						Type:          AddressTypeCC,
						Substitutions: map[string]interface{}{"key": "value"},
					},
					{
						Name:          "test bcc user",
						Address:       "test-bcc@and-period.jp",
						Type:          AddressTypeBCC,
						Substitutions: map[string]interface{}{"key": "value"},
					},
				},
			},
			expect: &mail.SGMailV3{
				From: &mail.Email{
					Name:    "test user",
					Address: "test-user@and-period.jp",
				},
				Subject: "テストタイトル",
				Personalizations: []*mail.Personalization{
					{
						To: []*mail.Email{{
							Name:    "test to user",
							Address: "test-to@and-period.jp",
						}},
						CC:            []*mail.Email{},
						BCC:           []*mail.Email{},
						Headers:       map[string]string{},
						Substitutions: map[string]string{},
						CustomArgs:    map[string]string{},
						DynamicTemplateData: map[string]interface{}{
							"key": "value",
						},
						Categories: []string{},
					},
					{
						To: []*mail.Email{},
						CC: []*mail.Email{{
							Name:    "test cc user",
							Address: "test-cc@and-period.jp",
						}},
						BCC:           []*mail.Email{},
						Headers:       map[string]string{},
						Substitutions: map[string]string{},
						CustomArgs:    map[string]string{},
						DynamicTemplateData: map[string]interface{}{
							"key": "value",
						},
						Categories: []string{},
					},
					{
						To: []*mail.Email{},
						CC: []*mail.Email{},
						BCC: []*mail.Email{{
							Name:    "test bcc user",
							Address: "test-bcc@and-period.jp",
						}},
						Headers:       map[string]string{},
						Substitutions: map[string]string{},
						CustomArgs:    map[string]string{},
						DynamicTemplateData: map[string]interface{}{
							"key": "value",
						},
						Categories: []string{},
					},
				},
				Content: []*mail.Content{{
					Type:  "type",
					Value: "value",
				}},
				Attachments: []*mail.Attachment{},
				Categories:  []string{"template-id"},
				TemplateID:  "d-templateid",
				SendAt:      int(now.Unix()),
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &client{
				now:         func() time.Time { return now },
				templateMap: map[string]string{"template-id": "d-templateid"},
			}
			assert.Equal(t, tt.expect, client.newMessage(
				tt.args.templateID,
				tt.args.fromName,
				tt.args.fromAddress,
				tt.args.subject,
				tt.args.contents,
				tt.args.personalizations,
			))
		})
	}
}
