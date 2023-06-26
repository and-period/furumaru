package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactRead(t *testing.T) {
	t.Parallel()
	type params struct {
		contactID string
		userType  int32
		content   string
	}
	tests := []struct {
		name   string
		params *NewContactReadParams
		expect *ContactRead
	}{
		{
			name: "success guest",
			params: &NewContactReadParams{
				ContactID: "contact-id",
				UserType:  ContactUserTypeGuest,
				Read:      false,
			},
			expect: &ContactRead{
				ContactID: "contact-id",
				UserType:  ContactUserTypeGuest,
				Read:      false,
			},
		},
		{
			name: "success admin",
			params: &NewContactReadParams{
				ContactID: "contact-id",
				UserType:  ContactUserTypeAdmin,
				Read:      false,
			},
			expect: &ContactRead{
				ContactID: "contact-id",
				UserType:  ContactUserTypeAdmin,
				Read:      false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewContactRead(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestContactRead_Fill(t *testing.T) {
	t.Parallel()
	contactRead := &ContactRead{
		ContactID: "contact-id",
		UserType:  ContactUserTypeAdmin,
		Read:      false,
	}
	tests := []struct {
		name        string
		contactRead *ContactRead
		userID      string
		expect      *ContactRead
	}{
		{
			name:        "success",
			contactRead: contactRead,
			userID:      "user-id",
			expect: &ContactRead{
				ContactID: "contact-id",
				UserType:  ContactUserTypeAdmin,
				UserID:    "user-id",
				Read:      false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.contactRead.Fill(tt.userID)
			assert.Equal(t, tt.expect, tt.contactRead)
		})
	}
}
