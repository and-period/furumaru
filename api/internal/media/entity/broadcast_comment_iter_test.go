package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastComments_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		comments BroadcastComments
	}{
		{
			name: "success",
			comments: BroadcastComments{
				{ID: "comment-id01", BroadcastID: "broadcast-id01", UserID: "user-id01"},
				{ID: "comment-id02", BroadcastID: "broadcast-id01", UserID: "user-id02"},
			},
		},
		{
			name:     "empty",
			comments: BroadcastComments{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, c := range tt.comments.All() {
				indices = append(indices, i)
				ids = append(ids, c.ID)
			}
			for i, c := range tt.comments {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, c.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.comments))
		})
	}
}

func TestBroadcastComments_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	comments := BroadcastComments{
		{ID: "comment-id01"},
		{ID: "comment-id02"},
		{ID: "comment-id03"},
	}
	var count int
	for range comments.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestBroadcastComments_IterUserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		comments BroadcastComments
		expect   []string
	}{
		{
			name: "success",
			comments: BroadcastComments{
				{ID: "comment-id01", UserID: "user-id01"},
				{ID: "comment-id02", UserID: "user-id02"},
			},
			expect: []string{"user-id01", "user-id02"},
		},
		{
			name: "skip empty user id",
			comments: BroadcastComments{
				{ID: "comment-id01", UserID: "user-id01"},
				{ID: "comment-id02", UserID: ""},
				{ID: "comment-id03", UserID: "user-id03"},
			},
			expect: []string{"user-id01", "user-id03"},
		},
		{
			name:     "empty",
			comments: BroadcastComments{},
			expect:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.comments.IterUserIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}
