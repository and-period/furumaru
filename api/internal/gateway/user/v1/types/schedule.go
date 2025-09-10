package types

// ScheduleStatus - 開催状況
type ScheduleStatus int32

const (
	ScheduleStatusUnknown  ScheduleStatus = 0
	ScheduleStatusWaiting  ScheduleStatus = 1 // ライブ配信前
	ScheduleStatusLive     ScheduleStatus = 2 // ライブ配信中
	ScheduleStatusClosed   ScheduleStatus = 3 // ライブ配信終了
	ScheduleStatusArchived ScheduleStatus = 4 // アーカイブ配信
)

// Schedule - マルシェ開催情報
type Schedule struct {
	ID                   string                        `json:"id"`                   // スケジュールID
	CoordinatorID        string                        `json:"coordinatorId"`        // コーディネータID
	Status               ScheduleStatus                `json:"status"`               // 開催状況
	Title                string                        `json:"title"`                // タイトル
	Description          string                        `json:"description"`          // 説明
	ThumbnailURL         string                        `json:"thumbnailUrl"`         // サムネイルURL
	DistributionURL      string                        `json:"distributionUrl"`      // 映像配信URL
	DistributionMetadata *ScheduleDistributionMetadata `json:"distributionMedatada"` // 映像メタデータ
	StartAt              int64                         `json:"startAt"`              // 配信開始日時
	EndAt                int64                         `json:"endAt"`                // 配信終了日時
}

type ScheduleDistributionMetadata struct {
	Subtitles map[string]string `json:"subtitles"` // 字幕情報
}

type ScheduleResponse struct {
	Schedule    *Schedule    `json:"schedule"`    // マルシェ開催情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Lives       []*Live      `json:"lives"`       // ライブ配信一覧
	Producers   []*Producer  `json:"producers"`   // 生産者一覧
	Products    []*Product   `json:"products"`    // 商品一覧
}

type LiveSchedulesResponse struct {
	Lives        []*LiveSummary `json:"lives"`        // 配信中・配信予定のマルシェ一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Total        int64          `json:"total"`        // 配信中・配信予定のマルシェ合計数
}

type ArchiveSchedulesResponse struct {
	Archives     []*ArchiveSummary `json:"archives"`     // 過去のマルシェ一覧
	Coordinators []*Coordinator    `json:"coordinators"` // コーディネータ一覧
	Total        int64             `json:"total"`        // 過去のマルシェ合計数
}
