package request

type CreateStoreRequest struct {
	Name string `json:"name,omitempty"` // 店舗名
}

type UpdateStoreRequest struct {
	Name         string `json:"name,omitempty"`         // 店舗名
	ThumbnailURL string `json:"thumbnailUrl,omitempty"` // サムネイルURL
}
