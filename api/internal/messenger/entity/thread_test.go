package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThread(t *testing.T) {
	t.Parallel()
	type params struct {
		contactID string
		userType  int32
		content   string
	}
	tests := []struct {
		name   string
		params *NewThreadParams
		expect *Thread
	}{
		{
			name: "success guest",
			params: &NewThreadParams{
				ContactID: "contact-id",
				UserType:  0,
				Content:   "content",
			},
			expect: &Thread{
				ContactID: "contact-id",
				UserType:  0,
				Content:   "content",
			},
		},
		{
			name: "success admin",
			params: &NewThreadParams{
				ContactID: "contact-id",
				UserType:  1,
				Content:   "content",
			},
			expect: &Thread{
				ContactID: "contact-id",
				UserType:  1,
				Content:   "content",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewThread(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestThread_Fill(t *testing.T) {
	t.Parallel()
	thread := &Thread{
		ContactID: "contact-id",
		UserType:  1,
		Content:   "content",
	}
	tests := []struct {
		name   string
		thread *Thread
		userID string
		expect *Thread
	}{
		{
			name:   "success",
			thread: thread,
			userID: "user-id",
			expect: &Thread{
				ContactID: "contact-id",
				UserType:  1,
				Content:   "content",
				UserID:    "user-id",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.thread.Fill(tt.userID)
			assert.Equal(t, tt.expect, tt.thread)
		})
	}
}
