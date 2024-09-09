package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoComment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewVideoCommentParams
		expect *VideoComment
	}{
		{
			name: "success",
			params: &NewVideoCommentParams{
				VideoID: "video-id",
				UserID:  "user-id",
				Content: "とても面白いですね",
			},
			expect: &VideoComment{
				VideoID: "video-id",
				UserID:  "user-id",
				Content: "とても面白いですね",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideoComment(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoComments_UserIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		comments VideoComments
		expect   []string
	}{
		{
			name: "success",
			comments: VideoComments{
				&VideoComment{ID: "comment-id01", UserID: ""},
				&VideoComment{ID: "comment-id02", UserID: "user-id"},
			},
			expect: []string{"user-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.comments.UserIDs()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
