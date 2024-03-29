package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
)

type LiveComment struct {
	response.LiveComment
}

type LiveComments []*LiveComment

func NewLiveComment(comment *mentity.BroadcastComment, user *User) *LiveComment {
	res := &LiveComment{
		LiveComment: response.LiveComment{
			ID:          comment.ID,
			Thumbnails:  []*response.Image{},
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
	res.Thumbnails = user.Thumbnails
	return res
}

func (c *LiveComment) Response() *response.LiveComment {
	return &c.LiveComment
}

func NewLiveComments(comments mentity.BroadcastComments, users map[string]*User) LiveComments {
	res := make(LiveComments, len(comments))
	for i := range comments {
		res[i] = NewLiveComment(comments[i], users[comments[i].UserID])
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
