package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactRead(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		params    *NewContactReadParams
		expect    *ContactRead
		expectErr error
	}{
		{
			name: "success guest",
			params: &NewContactReadParams{
				ContactID: "contact-id",
				UserType:  ContactUserTypeGuest,
				UserID:    "",
				Read:      false,
			},
			expect: &ContactRead{
				ID:        "", // ignore
				ContactID: "contact-id",
				UserType:  ContactUserTypeGuest,
				UserID:    "",
				Read:      false,
			},
			expectErr: nil,
		},
		{
			name: "success admin",
			params: &NewContactReadParams{
				ContactID: "contact-id",
				UserType:  ContactUserTypeAdmin,
				UserID:    "admin-id",
				Read:      false,
			},
			expect: &ContactRead{
				ID:        "", // ignore
				ContactID: "contact-id",
				UserType:  ContactUserTypeAdmin,
				UserID:    "admin-id",
				Read:      false,
			},
			expectErr: nil,
		},
		{
			name: "error invalid argument",
			params: &NewContactReadParams{
				ContactID: "contact-id",
				UserType:  ContactUserTypeAdmin,
				UserID:    "",
				Read:      false,
			},
			expect:    &ContactRead{},
			expectErr: errors.New("entity: failed to new contact read"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewContactRead(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectErr, err)
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.contactRead.Fill(tt.userID)
			assert.Equal(t, tt.expect, tt.contactRead)
		})
	}
}
