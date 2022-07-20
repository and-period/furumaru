package request

import "time"

type CreateNotificationRequest struct {
	Title       string       `json:"title,omitempty"`       //タイトル
	Body        string       `json:"body,omitempty"`        //本文
	Targets     []TargetType `json:"targets,omitempty"`     //掲載対象一覧
	PublishedAt time.Time    `json:"publishedAt,omitempty"` //掲載開始日
	Public      bool         `json:"public,omitempty"`      //公開フラグ
}

type TargetType int32

const (
	PostTargetAll          TargetType = 0 // 全員対象
	PostTargetUsers        TargetType = 1 // ユーザー対象
	PostTargetProducers    TargetType = 2 // 生産者対象
	PostTargetCoordinators TargetType = 3 // コーディネーター対象
)
