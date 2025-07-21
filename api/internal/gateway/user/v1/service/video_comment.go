package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

type VideoComment struct {
	response.VideoComment
}

type VideoComments []*VideoComment

func NewVideoComment(comment *mentity.VideoComment, user *uentity.User) *VideoComment {
	res := &VideoComment{
		VideoComment: response.VideoComment{
			Comment:     comment.Content,
			PublishedAt: comment.CreatedAt.Unix(),
		},
	}
	if user == nil || user.Status != uentity.UserStatusVerified {
		return res
	}
	res.UserID = user.ID
	res.Username = user.Username()
	res.AccountID = user.AccountID
	res.ThumbnailURL = user.ThumbnailURL
	return res
}

func (c *VideoComment) Response() *response.VideoComment {
	return &c.VideoComment
}

func NewVideoComments(
	comments mentity.VideoComments,
	users map[string]*uentity.User,
) VideoComments {
	res := make(VideoComments, 0, len(comments))
	for _, comment := range comments {
		if comment.Disabled {
			continue
		}
		res = append(res, NewVideoComment(comment, users[comment.UserID]))
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
