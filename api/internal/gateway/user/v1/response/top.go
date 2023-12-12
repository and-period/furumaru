package response

type TopCommonResponse struct {
	Lives        []*TopCommonLive    `json:"lives"`        // 配信中・配信予定のマルシェ一覧
	Archives     []*TopCommonArchive `json:"archives"`     // 過去のマルシェ一覧
	Coordinators []*Coordinator      `json:"coordinators"` // コーディネータ一覧
}

// TopCommonLive - 配信中・配信予定のマルシェ
type TopCommonLive struct {
	ScheduleID    string                  `json:"scheduleId"`    // 開催スケジュールID
	CoordinatorID string                  `json:"coordinatorId"` // コーディネータID
	Status        int32                   `json:"status"`        // 開催状況
	Title         string                  `json:"title"`         // タイトル
	ThumbnailURL  string                  `json:"thumbnailUrl"`  // サムネイルURL
	Thumbnails    []*Image                `json:"thumbnails"`    // サムネイル一覧(リサイズ済み)
	StartAt       int64                   `json:"startAt"`       // 開催開始日時
	EndAt         int64                   `json:"endAt"`         // 開催終了日時
	Products      []*TopCommonLiveProduct `json:"products"`      // 販売商品一覧
}

type TopCommonLiveProduct struct {
	ProductID    string   `json:"id"`           // 商品ID
	Name         string   `json:"name"`         // 商品名
	Price        int64    `json:"price"`        // 販売価格
	Inventory    int64    `json:"inventory"`    // 在庫数
	ThumbnailURL string   `json:"thumbnailUrl"` // サムネイルURL
	Thumbnails   []*Image `json:"thumbnails"`   // サムネイル一覧(リサイズ済み)
}

// TopCommonArchive - 過去のマルシェ
type TopCommonArchive struct {
	ScheduleID    string   `json:"scheduleId"`    // 開催スケジュールID
	CoordinatorID string   `json:"coordinatorId"` // コーディネータID
	Title         string   `json:"title"`         // タイトル
	ThumbnailURL  string   `json:"thumbnailUrl"`  // サムネイルURL
	Thumbnails    []*Image `json:"thumbnails"`    // サムネイル一覧(リサイズ済み)
}
