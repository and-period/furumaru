package request

type CreateSpotRequest struct {
	TypeID       string  `json:"spotTypeId" validate:"required"`                 // スポット種別ID
	Name         string  `json:"name" validate:"required,max=64"`                // スポット名
	Description  string  `json:"description" validate:"omitempty,max=2000"`      // 説明
	ThumbnailURL string  `json:"thumbnailUrl" validate:"omitempty,url"`          // サムネイルURL
	Latitude     float64 `json:"latitude" validate:"required,min=-90,max=90"`    // 緯度
	Longitude    float64 `json:"longitude" validate:"required,min=-180,max=180"` // 経度
}

type UpdateSpotRequest struct {
	TypeID       string  `json:"spotTypeId" validate:"required"`                 // スポット種別ID
	Name         string  `json:"name" validate:"required,max=64"`                // スポット名
	Description  string  `json:"description" validate:"omitempty,max=2000"`      // 説明
	ThumbnailURL string  `json:"thumbnailUrl" validate:"omitempty,url"`          // サムネイルURL
	Latitude     float64 `json:"latitude" validate:"required,min=-90,max=90"`    // 緯度
	Longitude    float64 `json:"longitude" validate:"required,min=-180,max=180"` // 経度
}
