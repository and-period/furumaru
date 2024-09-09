package entity

import "time"

type AggregateVideoViewerLogInterval string

const (
	AggregateVideoViewerLogIntervalSecond AggregateVideoViewerLogInterval = "%Y-%m-%d %H:%i:%s"
	AggregateVideoViewerLogIntervalMinute AggregateVideoViewerLogInterval = "%Y-%m-%d %H:%i:00"
	AggregateVideoViewerLogIntervalHour   AggregateVideoViewerLogInterval = "%Y-%m-%d %H:00:00"
)

// VideoViewerLog - オンデマンド配信視聴履歴情報
type VideoViewerLog struct {
	VideoID   string    `gorm:"primaryKey;<-:create"` // オンデマンド配信ID
	SessionID string    `gorm:"primaryKey;<-:create"` // セッションID
	CreatedAt time.Time `gorm:"primaryKey;<-:create"` // 登録日時
	UserID    string    `gorm:"default:null"`         // ユーザーID
	UserAgent string    `gorm:""`                     // ユーザーエージェント
	ClientIP  string    `gorm:""`                     // 接続元IPアドレス
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type VideoViewerLogs []*VideoViewerLog

// AggregatedVideoViewerLog - オンデマンド配信視聴履歴集計情報
type AggregatedVideoViewerLog struct {
	VideoID    string    // オンデマンド配信ID
	ReportedAt time.Time // 集計日時
	Total      int64     // 視聴合計回数
}

type AggregatedVideoViewerLogs []*AggregatedVideoViewerLog

type NewVideoViewerLogParams struct {
	VideoID   string
	SessionID string
	UserID    string
	UserAgent string
	ClientIP  string
}

func NewVideoViewerLog(params *NewVideoViewerLogParams) *VideoViewerLog {
	return &VideoViewerLog{
		VideoID:   params.VideoID,
		SessionID: params.SessionID,
		UserID:    params.UserID,
		UserAgent: params.UserAgent,
		ClientIP:  params.ClientIP,
	}
}

func (ls AggregatedVideoViewerLogs) MapByReportedAt() map[time.Time]*AggregatedVideoViewerLog {
	res := make(map[time.Time]*AggregatedVideoViewerLog)
	for _, l := range ls {
		res[l.ReportedAt] = l
	}
	return res
}

func (ls AggregatedVideoViewerLogs) GroupByVideoID() map[string]AggregatedVideoViewerLogs {
	res := make(map[string]AggregatedVideoViewerLogs)
	for _, l := range ls {
		if _, ok := res[l.VideoID]; !ok {
			res[l.VideoID] = make(AggregatedVideoViewerLogs, 0, len(ls))
		}
		res[l.VideoID] = append(res[l.VideoID], l)
	}
	return res
}
