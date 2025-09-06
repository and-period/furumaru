package request

type CreateSpotRequest struct {
	TypeID       string  `json:"spotTypeId" binding:"required"`                 // スポット種別ID
	Name         string  `json:"name" binding:"required,max=64"`                // スポット名
	Description  string  `json:"description" binding:"omitempty,max=2000"`      // 説明
	ThumbnailURL string  `json:"thumbnailUrl" binding:"omitempty,url"`          // サムネイルURL
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`    // 緯度
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"` // 経度
}

type UpdateSpotRequest struct {
	TypeID       string  `json:"spotTypeId" binding:"required"`                 // スポット種別ID
	Name         string  `json:"name" binding:"required,max=64"`                // スポット名
	Description  string  `json:"description" binding:"omitempty,max=2000"`      // 説明
	ThumbnailURL string  `json:"thumbnailUrl" binding:"omitempty,url"`          // サムネイルURL
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`    // 緯度
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"` // 経度
}
