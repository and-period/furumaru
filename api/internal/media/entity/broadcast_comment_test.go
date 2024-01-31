package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastComment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *BroadcastCommentParams
		expect *BroadcastComment
	}{
		{
			name: "success",
			params: &BroadcastCommentParams{
				BroadcastID: "broadcast-id",
				UserID:      "user-id",
				Content:     "こんにちは",
			},
			expect: &BroadcastComment{
				ID:          "",
				BroadcastID: "broadcast-id",
				UserID:      "user-id",
				Content:     "こんにちは",
				Disabled:    false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcastComment(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestBroadcastComments_UserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		comments BroadcastComments
		expect   []string
	}{
		{
			name: "success",
			comments: BroadcastComments{
				{
					ID:          "comment-id",
					BroadcastID: "broadcast-id",
					UserID:      "user-id",
					Content:     "こんにちは",
					Disabled:    false,
				},
			},
			expect: []string{"user-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.comments.UserIDs())
		})
	}
}
