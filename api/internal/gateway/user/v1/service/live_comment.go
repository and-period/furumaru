package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

type LiveComment struct {
	response.LiveComment
}

type LiveComments []*LiveComment

func NewLiveComment(comment *mentity.BroadcastComment, user *uentity.User) *LiveComment {
	res := &LiveComment{
		LiveComment: response.LiveComment{
			Comment:     comment.Content,
			PublishedAt: comment.CreatedAt.Unix(),
		},
	}
	if user == nil || user.Status != uentity.UserStatusVerified {
		return res
	}
	res.UserID = user.ID
	res.Username = user.Username
	res.AccountID = user.AccountID
	res.ThumbnailURL = user.ThumbnailURL
	return res
}

func (c *LiveComment) Response() *response.LiveComment {
	return &c.LiveComment
}

func NewLiveComments(comments mentity.BroadcastComments, users map[string]*uentity.User) LiveComments {
	res := make(LiveComments, 0, len(comments))
	for _, comment := range comments {
		if comment.Disabled {
			continue
		}
		res = append(res, NewLiveComment(comment, users[comment.UserID]))
	}
	return res
}

func (cs LiveComments) Response() []*response.LiveComment {
	res := make([]*response.LiveComment, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
