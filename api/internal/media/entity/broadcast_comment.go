package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

type BroadcastCommentOrderBy string

const (
	BroadcastCommentOrderByCreatedAt BroadcastCommentOrderBy = "created_at"
)

// BroadcastComment - ライブ配信コメント情報
type BroadcastComment struct {
	ID          string    `gorm:"primaryKey;<-:create"` // コメントID
	BroadcastID string    `gorm:""`                     // ライブ配信ID
	UserID      string    `gorm:""`                     // ユーザーID
	Content     string    `gorm:""`                     // コメント内容
	Disabled    bool      `gorm:""`                     // コメント無効フラグ
	CreatedAt   time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time `gorm:""`                     // 更新日時
}

type BroadcastComments []*BroadcastComment

type BroadcastCommentParams struct {
	BroadcastID string
	UserID      string
	Content     string
}

func NewBroadcastComment(params *BroadcastCommentParams) *BroadcastComment {
	return &BroadcastComment{
		ID:          uuid.Base58Encode(uuid.New()),
		BroadcastID: params.BroadcastID,
		UserID:      params.UserID,
		Content:     params.Content,
		Disabled:    false,
	}
}

func (cs BroadcastComments) UserIDs() []string {
	set := set.NewEmpty[string](len(cs))
	for i := range cs {
		if cs[i].UserID == "" {
			continue
		}
		set.Add(cs[i].UserID)
	}
	return set.Slice()
}
