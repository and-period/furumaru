package entity

import "time"

// Rehearsal - リハーサルスケジュール
type Rehearsal struct {
	LiveID         string    `dynamodbav:"live_id"`             // ライブ配信ID
	ScheduleID     string    `dynamodbav:"schedule_id"`         // 開催スケジュールID
	ChannelName    string    `dyanmodbav:"channel_name"`        // IVS チャンネル名
	ChannelArn     string    `dynamodbav:"channel_arn"`         // IVS チャンネルARN
	IngestEndpoint string    `dynamodbav:"ingest_endpoint"`     // IVS 配信取り込みエンドポイント
	StreamKey      string    `dynamodbav:"stream_key"`          // IVS 配信用ストリームキー
	PlaybackURL    string    `dynamodbav:"playback_url"`        // IVS 再生用URL
	StartAt        time.Time `dynamodbav:"start_at"`            // 配信開始日時
	ExpiresAt      time.Time `dynamodbav:"expires_at,unixtime"` // リハーサルモード実施有効期限
	CreatedAt      time.Time `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt      time.Time `dynamodbav:"updated_at"`          // 更新日時
}

func (s *Rehearsal) TableName() string {
	return "rehearsals"
}

func (s *Rehearsal) PrimaryKey() map[string]interface{} {
	return map[string]interface{}{
		"live_id": s.LiveID,
	}
}
