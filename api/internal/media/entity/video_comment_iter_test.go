package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoComments_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		comments VideoComments
	}{
		{
			name: "success",
			comments: VideoComments{
				{ID: "comment-id01", VideoID: "video-id01", UserID: "user-id01"},
				{ID: "comment-id02", VideoID: "video-id01", UserID: "user-id02"},
			},
		},
		{
			name:     "empty",
			comments: VideoComments{},
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

func TestVideoComments_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	comments := VideoComments{
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

func TestVideoComments_IterUserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		comments VideoComments
		expect   []string
	}{
		{
			name: "success",
			comments: VideoComments{
				{ID: "comment-id01", UserID: "user-id01"},
				{ID: "comment-id02", UserID: "user-id02"},
			},
			expect: []string{"user-id01", "user-id02"},
		},
		{
			name: "skip empty user id",
			comments: VideoComments{
				{ID: "comment-id01", UserID: "user-id01"},
				{ID: "comment-id02", UserID: ""},
				{ID: "comment-id03", UserID: "user-id03"},
			},
			expect: []string{"user-id01", "user-id03"},
		},
		{
			name:     "empty",
			comments: VideoComments{},
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
