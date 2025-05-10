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

	tests := []struct {
		name        string
		contact     *Contact
		categoryID  string
		userID      string
		responderID string
		expect      *Contact
	}{
		{
			name: "success",
			contact: &Contact{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				Note:        "",
			},
			categoryID:  "category-id",
			userID:      "user-id",
			responderID: "responder-id",
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
		{
			name: "success without userID",
			contact: &Contact{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				Note:        "",
			},
			categoryID:  "category-id",
			userID:      "",
			responderID: "responder-id",
			expect: &Contact{
				Title:       "お問い合わせ件名",
				CategoryID:  "category-id",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				UserID:      "",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				ResponderID: "responder-id",
				Note:        "",
			},
		},
		{
			name: "success without responderID",
			contact: &Contact{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				Note:        "",
			},
			categoryID:  "category-id",
			userID:      "user-id",
			responderID: "",
			expect: &Contact{
				Title:       "お問い合わせ件名",
				CategoryID:  "category-id",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				UserID:      "user-id",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				ResponderID: "",
				Note:        "",
			},
		},
		{
			name: "success only categoryID",
			contact: &Contact{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				Note:        "",
			},
			categoryID:  "category-id",
			userID:      "",
			responderID: "",
			expect: &Contact{
				Title:       "お問い合わせ件名",
				CategoryID:  "category-id",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				UserID:      "",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      ContactStatusUnknown,
				ResponderID: "",
				Note:        "",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.contact.Fill(tt.categoryID, tt.userID, tt.responderID)
			assert.Equal(t, tt.expect, tt.contact)
		})
	}
}

func TestContacts_ContactCategoryIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []string
	}{
		{
			name: "success",
			contacts: Contacts{
				{
					ID:          "contact-id",
					Title:       "お問い合わせ件名",
					Content:     "お問い合わせ内容",
					Username:    "お問い合わせ氏名",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      ContactStatusUnknown,
					Note:        "",
					CategoryID:  "category-id",
					UserID:      "",
					ResponderID: "responder-id",
				},
			},
			expect: []string{"category-id"},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.expect, tt.contacts.CategoryIDs())
		})
	}
}

func TestContacts_UserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []string
	}{
		{
			name: "success",
			contacts: Contacts{
				{
					ID:          "contact-id",
					Title:       "お問い合わせ件名",
					Content:     "お問い合わせ内容",
					Username:    "お問い合わせ氏名",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      ContactStatusUnknown,
					Note:        "",
					CategoryID:  "category-id",
					UserID:      "user-id",
					ResponderID: "responder-id",
				},
			},
			expect: []string{"user-id"},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.expect, tt.contacts.UserIDs())
		})
	}
}

func TestContacts_ResponderIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []string
	}{
		{
			name: "success",
			contacts: Contacts{
				{
					ID:          "contact-id",
					Title:       "お問い合わせ件名",
					Content:     "お問い合わせ内容",
					Username:    "お問い合わせ氏名",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      ContactStatusUnknown,
					Note:        "",
					CategoryID:  "category-id",
					UserID:      "user-id",
					ResponderID: "responder-id",
				},
			},
			expect: []string{"responder-id"},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.expect, tt.contacts.ResponderIDs())
		})
	}
}
