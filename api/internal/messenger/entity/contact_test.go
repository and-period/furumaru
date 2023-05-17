package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContact(t *testing.T) {
	t.Parallel()
	type params struct {
		title       string
		content     string
		username    string
		email       string
		phoneNumber string
	}
	tests := []struct {
		name   string
		params *NewContactParams
		expect *Contact
	}{
		{
			name: "success",
			params: &NewContactParams{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
			},
			expect: &Contact{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				Note:        "",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewContact(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestContact_Fill(t *testing.T) {
	t.Parallel()
	contact := &Contact{
		Title:       "お問い合わせ件名",
		Content:     "お問い合わせ内容",
		Username:    "お問い合わせ氏名",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
		Status:      ContactStatusUnknown,
		Note:        "",
	}

	tests := []struct {
		name        string
		contact     *Contact
		categoryId  string
		userId      string
		responderId string
		expect      *Contact
	}{
		{
			name:        "success",
			contact:     contact,
			categoryId:  "category-id",
			userId:      "user-id",
			responderId: "responder-id",
			expect: &Contact{
				Title:       "お問い合わせ件名",
				CategoryID:  "category-id",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				UserID:      "user-id",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				ResponderID: "responder-id",
				Note:        "",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.contact.Fill(tt.categoryId, tt.userId, tt.responderId)
			assert.Equal(t, tt.expect, tt.contact)
		})
	}

}
