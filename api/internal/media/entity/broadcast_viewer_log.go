package entity

import "time"

type AggregateBroadcastViewerLogInterval string

const (
	AggregateBroadcastViewerLogIntervalSecond AggregateBroadcastViewerLogInterval = "%Y-%m-%d %H:%i:%s"
	AggregateBroadcastViewerLogIntervalMinute AggregateBroadcastViewerLogInterval = "%Y-%m-%d %H:%i:00"
	AggregateBroadcastViewerLogIntervalHour   AggregateBroadcastViewerLogInterval = "%Y-%m-%d %H:00:00"
)

var ExcludeUserAgentLogs = []string{
	"node",
	"Google-Safety",
}

// BroadcastViewerLog - ライブ視聴履歴情報
type BroadcastViewerLog struct {
	BroadcastID string    `gorm:"primaryKey;<-:create"` // ライブ配信ID
	SessionID   string    `gorm:"primaryKey;<-:create"` // セッションID
	CreatedAt   time.Time `gorm:"primaryKey;<-:create"` // 登録日時
	UserID      string    `gorm:"default:null"`         // ユーザーID
	UserAgent   string    `gorm:""`                     // ユーザーエージェント
	ClientIP    string    `gorm:""`                     // 接続元IPアドレス
	UpdatedAt   time.Time `gorm:""`                     // 更新日時
}

type BroadcastViewerLogs []*BroadcastViewerLog

// AggregatedBroadcastViewerLog - ライブ視聴履歴集計情報
type AggregatedBroadcastViewerLog struct {
	BroadcastID string    // ライブ配信ID
	Timestamp   time.Time // 集計日時
	Total       int64     // 視聴合計回数
}

type AggregatedBroadcastViewerLogs []*AggregatedBroadcastViewerLog

type BroadcastViewerLogParams struct {
	BroadcastID string
	SessionID   string
	UserID      string
	UserAgent   string
	ClientIP    string
}

func NewBroadcastViewerLog(params *BroadcastViewerLogParams) *BroadcastViewerLog {
	return &BroadcastViewerLog{
		BroadcastID: params.BroadcastID,
		SessionID:   params.SessionID,
		UserID:      params.UserID,
		UserAgent:   params.UserAgent,
		ClientIP:    params.ClientIP,
	}
}

func (ls AggregatedBroadcastViewerLogs) MapByTimestamp() map[time.Time]*AggregatedBroadcastViewerLog {
	res := make(map[time.Time]*AggregatedBroadcastViewerLog, len(ls))
	for _, l := range ls {
		res[l.Timestamp] = l
	}
	return res
}

func (ls AggregatedBroadcastViewerLogs) GroupByBroadcastID() map[string]AggregatedBroadcastViewerLogs {
	res := make(map[string]AggregatedBroadcastViewerLogs, len(ls))
	for _, l := range ls {
		if _, ok := res[l.BroadcastID]; !ok {
			res[l.BroadcastID] = make(AggregatedBroadcastViewerLogs, 0, len(ls))
		}
		res[l.BroadcastID] = append(res[l.BroadcastID], l)
	}
	return res
}
