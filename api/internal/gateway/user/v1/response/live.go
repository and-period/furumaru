package response

// Live - ライブ配信情報
type Live struct {
	ScheduleID string   `json:"scheduleId"` // マルシェ開催スケジュールID
	ProducerID string   `json:"producerId"` // 生産者ID
	ProductIDs []string `json:"productIds"` // 商品ID一覧
	Comment    string   `json:"comment"`    // コメント
	StartAt    int64    `json:"startAt"`    // ライブ配信開始日時
	EndAt      int64    `json:"endAt"`      // ライブ配信終了日時
}

// LiveSummary - ライブ配信の概要
type LiveSummary struct {
	ScheduleID    string         `json:"scheduleId"`    // 開催スケジュールID
	CoordinatorID string         `json:"coordinatorId"` // コーディネータID
	Status        int32          `json:"status"`        // 開催状況
	Title         string         `json:"title"`         // タイトル
	ThumbnailURL  string         `json:"thumbnailUrl"`  // サムネイルURL
	Thumbnails    []*Image       `json:"thumbnails"`    // サムネイル一覧(リサイズ済み)
	StartAt       int64          `json:"startAt"`       // 開催開始日時
	EndAt         int64          `json:"endAt"`         // 開催終了日時
	Products      []*LiveProduct `json:"products"`      // 販売商品一覧
}

// LiveProduct - ライブ配信で取り上げる商品情報
type LiveProduct struct {
	ProductID    string   `json:"id"`           // 商品ID
	Name         string   `json:"name"`         // 商品名
	Price        int64    `json:"price"`        // 販売価格(税込)
	Inventory    int64    `json:"inventory"`    // 在庫数
	ThumbnailURL string   `json:"thumbnailUrl"` // サムネイルURL
	Thumbnails   []*Image `json:"thumbnails"`   // サムネイル一覧(リサイズ済み)
}

// ArchiveSummary - オンデマンド配信の概要
type ArchiveSummary struct {
	ScheduleID    string   `json:"scheduleId"`    // 開催スケジュールID
	CoordinatorID string   `json:"coordinatorId"` // コーディネータID
	Title         string   `json:"title"`         // タイトル
	StartAt       int64    `json:"startAt"`       // 開催開始日時
	EndAt         int64    `json:"endAt"`         // 開催終了日時
	ThumbnailURL  string   `json:"thumbnailUrl"`  // サムネイルURL
	Thumbnails    []*Image `json:"thumbnails"`    // サムネイル一覧(リサイズ済み)
}
