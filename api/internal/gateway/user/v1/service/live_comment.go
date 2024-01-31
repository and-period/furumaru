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
	return &LiveComment{
		LiveComment: response.LiveComment{
			UserID:       user.ID,
			Username:     user.Username,
			AccountID:    user.AccountID,
			ThumbnailURL: user.ThumbnailURL,
			Comment:      comment.Content,
			PublishedAt:  comment.CreatedAt.Unix(),
		},
	}
}

func (c *LiveComment) Response() *response.LiveComment {
	return &c.LiveComment
}

func NewLiveComments(comments mentity.BroadcastComments, users map[string]*uentity.User) LiveComments {
	res := make(LiveComments, 0, len(comments))
	for _, comment := range comments {
		u, ok := users[comment.UserID]
		if !ok || u.Status != uentity.UserStatusVerified {
			continue
		}
		res = append(res, NewLiveComment(comment, u))
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
