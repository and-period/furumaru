package request

type CreateSpotRequest struct {
	TypeID       string  `json:"spotTypeId,omitempty"`   // スポット種別ID
	Name         string  `json:"name,omitempty"`         // スポット名
	Description  string  `json:"description,omitempty"`  // 説明
	ThumbnailURL string  `json:"thumbnailUrl,omitempty"` // サムネイルURL
	Latitude     float64 `json:"latitude,omitempty"`     // 緯度
	Longitude    float64 `json:"longitude,omitempty"`    // 経度
}

type UpdateSpotRequest struct {
	TypeID       string  `json:"spotTypeId,omitempty"`   // スポット種別ID
	Name         string  `json:"name,omitempty"`         // スポット名
	Description  string  `json:"description,omitempty"`  // 説明
	ThumbnailURL string  `json:"thumbnailUrl,omitempty"` // サムネイルURL
	Latitude     float64 `json:"latitude,omitempty"`     // 緯度
	Longitude    float64 `json:"longitude,omitempty"`    // 経度
}
