package response

// Spot - スポット情報
type Spot struct {
	ID               string  `json:"id"`               // スポットID
	Name             string  `json:"name"`             // スポット名
	Description      string  `json:"description"`      // 説明
	ThumbnailURL     string  `json:"thumbnailUrl"`     // サムネイル画像URL
	Longitude        float64 `json:"longitude"`        // 座標情報:経度
	Latitude         float64 `json:"latitude"`         // 座標情報:緯度
	UserType         int32   `json:"userType"`         // 投稿者の種別
	UserID           string  `json:"userId"`           // 投稿者のユーザーID
	Username         string  `json:"userName"`         // 投稿者名
	UserThumbnailURL string  `json:"userThumbnailUrl"` // 投稿者のサムネイルURL
	CreatedAt        int64   `json:"createdAt"`        // 登録日時
	UpdatedAt        int64   `json:"updatedAt"`        // 更新日時
}

type SpotResponse struct {
	Spot *Spot `json:"spot"` // スポット情報
}

type SpotsResponse struct {
	Spots []*Spot `json:"spots"` // スポット一覧
}
