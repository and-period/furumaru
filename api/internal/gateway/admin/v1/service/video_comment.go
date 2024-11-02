package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
)

type VideoComment struct {
	response.VideoComment
}

type VideoComments []*VideoComment

func NewVideoComment(comment *mentity.VideoComment, user *User) *VideoComment {
	res := &VideoComment{
		VideoComment: response.VideoComment{
			ID:          comment.ID,
			Comment:     comment.Content,
			Disabled:    comment.Disabled,
			PublishedAt: comment.CreatedAt.Unix(),
		},
	}
	if user == nil || UserStatus(user.Status) != UserStatusVerified {
		return res
	}
	res.UserID = user.ID
	res.Username = user.Username
	res.AccountID = user.AccountID
	res.ThumbnailURL = user.ThumbnailURL
	return res
}

func (c *VideoComment) Response() *response.VideoComment {
	return &c.VideoComment
}

func NewVideoComments(comments mentity.VideoComments, users map[string]*User) VideoComments {
	res := make(VideoComments, len(comments))
	for i := range comments {
		res[i] = NewVideoComment(comments[i], users[comments[i].UserID])
	}
	return res
}

func (cs VideoComments) Response() []*response.VideoComment {
	res := make([]*response.VideoComment, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
