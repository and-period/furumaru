package entity

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContacts_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
	}{
		{
			name: "success",
			contacts: Contacts{
				{
					ID:          "contact-id01",
					Title:       "お問い合わせ件名1",
					CategoryID:  "category-id01",
					UserID:      "user-id01",
					ResponderID: "responder-id01",
				},
				{
					ID:          "contact-id02",
					Title:       "お問い合わせ件名2",
					CategoryID:  "category-id02",
					UserID:      "user-id02",
					ResponderID: "responder-id02",
				},
			},
		},
		{
			name:     "empty",
			contacts: Contacts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, c := range tt.contacts.All() {
				indices = append(indices, i)
				ids = append(ids, c.ID)
			}
			for i, c := range tt.contacts {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, c.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.contacts))
		})
	}
}

func TestContacts_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	contacts := Contacts{
		{ID: "contact-id01"},
		{ID: "contact-id02"},
		{ID: "contact-id03"},
	}
	var count int
	for range contacts.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestContacts_IterIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []string
	}{
		{
			name: "success",
			contacts: Contacts{
				{ID: "contact-id01"},
				{ID: "contact-id02"},
				{ID: "contact-id03"},
			},
			expect: []string{"contact-id01", "contact-id02", "contact-id03"},
		},
		{
			name:     "empty",
			contacts: Contacts{},
			expect:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.contacts.IterIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestContacts_IterCategoryIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []string
	}{
		{
			name: "success",
			contacts: Contacts{
				{ID: "contact-id01", CategoryID: "category-id01"},
				{ID: "contact-id02", CategoryID: "category-id02"},
			},
			expect: []string{"category-id01", "category-id02"},
		},
		{
			name:     "empty",
			contacts: Contacts{},
			expect:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.contacts.IterCategoryIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestContacts_IterUserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []string
	}{
		{
			name: "success",
			contacts: Contacts{
				{ID: "contact-id01", UserID: "user-id01"},
				{ID: "contact-id02", UserID: "user-id02"},
			},
			expect: []string{"user-id01", "user-id02"},
		},
		{
			name:     "empty",
			contacts: Contacts{},
			expect:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.contacts.IterUserIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestContacts_IterResponderIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []string
	}{
		{
			name: "success",
			contacts: Contacts{
				{ID: "contact-id01", ResponderID: "responder-id01"},
				{ID: "contact-id02", ResponderID: "responder-id02"},
			},
			expect: []string{"responder-id01", "responder-id02"},
		},
		{
			name:     "empty",
			contacts: Contacts{},
			expect:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.contacts.IterResponderIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}
