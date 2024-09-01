package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

// VideoComment - オンデマンド配信コメント情報
type VideoComment struct {
	ID        string    `gorm:"primaryKey;<-:create"` // コメントID
	VideoID   string    `gorm:""`                     // オンデマンド配信ID
	UserID    string    `gorm:""`                     // ユーザーID
	Content   string    `gorm:""`                     // コメント内容
	Disabled  bool      `gorm:""`                     // コメント無効フラグ
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type VideoComments []*VideoComment

type NewVideoCommentParams struct {
	VideoID string
	UserID  string
	Content string
}

func NewVideoComment(params *NewVideoCommentParams) *VideoComment {
	return &VideoComment{
		ID:       uuid.Base58Encode(uuid.New()),
		VideoID:  params.VideoID,
		UserID:   params.UserID,
		Content:  params.Content,
		Disabled: false,
	}
}

func (cs VideoComments) UserIDs() []string {
	set := set.NewEmpty[string](len(cs))
	for i := range cs {
		if cs[i].UserID == "" {
			continue
		}
		set.Add(cs[i].UserID)
	}
	return set.Slice()
}
