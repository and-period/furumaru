package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContact(t *testing.T) {
	t.Parallel()
	type input struct {
		title       string
		content     string
		username    string
		email       string
		phoneNumber string
	}
	tests := []struct {
		name   string
		input  input
		expect *Contact
	}{
		{
			name: "success",
			input: input{
				title:       "お問い合わせ件名",
				content:     "お問い合わせ内容",
				username:    "お問い合わせ氏名",
				email:       "test-user@and-period.jp",
				phoneNumber: "+819012345678",
			},
			expect: &Contact{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				Priority:    ContactPriorityUnknown,
				Note:        "",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewContact(
				tt.input.title,
				tt.input.content,
				tt.input.username,
				tt.input.email,
				tt.input.phoneNumber,
			)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
