package response

// Spot - スポット情報
type Spot struct {
	ID           string  `json:"id"`           // スポットID
	UserType     int32   `json:"userType"`     // 投稿者の種別
	UserID       string  `json:"userId"`       // ユーザーID
	Name         string  `json:"name"`         // スポット名
	Description  string  `json:"description"`  // 説明
	ThumbnailURL string  `json:"thumbnailUrl"` // サムネイル画像URL
	Longitude    float64 `json:"longitude"`    // 座標情報:経度
	Latitude     float64 `json:"latitude"`     // 座標情報:緯度
	Approved     bool    `json:"approved"`     // 承認フラグ
	CreatedAt    int64   `json:"createdAt"`    // 登録日時
	UpdatedAt    int64   `json:"updatedAt"`    // 更新日時
}

type SpotResponse struct {
	Spot        *Spot        `json:"spot"`        // スポット情報
	User        *User        `json:"user"`        // ユーザ情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Producer    *Producer    `json:"producer"`    // 生産者情報
}

type SpotsResponse struct {
	Spots        []*Spot        `json:"spots"`        // スポット一覧
	Users        []*User        `json:"users"`        // ユーザ一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Producers    []*Producer    `json:"producers"`    // 生産者一覧
	Total        int64          `json:"total"`        // 合計数
}
