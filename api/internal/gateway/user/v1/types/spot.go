package types

// SpotUserType - 投稿者の種別
type SpotUserType int32

const (
	SpotUserTypeUnknown     SpotUserType = 0
	SpotUserTypeUser        SpotUserType = 1 // ユーザー
	SpotUserTypeCoordinator SpotUserType = 2 // コーディネータ
	SpotUserTypeProducer    SpotUserType = 3 // 生産者
)

// Spot - スポット情報
type Spot struct {
	ID               string       `json:"id"`                   // スポットID
	TypeID           string       `json:"spotTypeId,omitempty"` // スポット種別ID
	Name             string       `json:"name"`                 // スポット名
	Description      string       `json:"description"`          // 説明
	ThumbnailURL     string       `json:"thumbnailUrl"`         // サムネイル画像URL
	Longitude        float64      `json:"longitude"`            // 座標情報:経度
	Latitude         float64      `json:"latitude"`             // 座標情報:緯度
	UserType         SpotUserType `json:"userType"`             // 投稿者の種別
	UserID           string       `json:"userId"`               // 投稿者のユーザーID
	Username         string       `json:"userName"`             // 投稿者名
	UserThumbnailURL string       `json:"userThumbnailUrl"`     // 投稿者のサムネイルURL
	CreatedAt        int64        `json:"createdAt"`            // 登録日時
	UpdatedAt        int64        `json:"updatedAt"`            // 更新日時
}

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

type SpotResponse struct {
	Spot     *Spot     `json:"spot"`     // スポット情報
	SpotType *SpotType `json:"spotType"` // スポット種別情報
}

type SpotsResponse struct {
	Spots     []*Spot     `json:"spots"`     // スポット一覧
	SpotTypes []*SpotType `json:"spotTypes"` // スポット種別一覧
}
