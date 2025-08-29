package request

type CreateSpotRequest struct {
	TypeID       string  `json:"spotTypeId"`   // スポット種別ID
	Name         string  `json:"name"`         // スポット名
	Description  string  `json:"description"`  // 説明
	ThumbnailURL string  `json:"thumbnailUrl"` // サムネイルURL
	Latitude     float64 `json:"latitude"`     // 緯度
	Longitude    float64 `json:"longitude"`    // 経度
}

type UpdateSpotRequest struct {
	TypeID       string  `json:"spotTypeId"`   // スポット種別ID
	Name         string  `json:"name"`         // スポット名
	Description  string  `json:"description"`  // 説明
	ThumbnailURL string  `json:"thumbnailUrl"` // サムネイルURL
	Latitude     float64 `json:"latitude"`     // 緯度
	Longitude    float64 `json:"longitude"`    // 経度
}
