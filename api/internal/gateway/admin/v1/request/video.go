package request

type CreateVideoRequest struct {
	Title             string   `json:"title,omitempty"`             // タイトル
	Description       string   `json:"description,omitempty"`       // 説明
	CoordinatorID     string   `json:"coordinatorId,omitempty"`     // コーディネータID
	ProductIDs        []string `json:"productIds,omitempty"`        // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds,omitempty"`     // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl,omitempty"`      // サムネイルURL
	VideoURL          string   `json:"videoUrl,omitempty"`          // 動画URL
	Public            bool     `json:"public,omitempty"`            // 公開設定
	Limited           bool     `json:"limited,omitempty"`           // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct,omitempty"`    // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience,omitempty"` // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt,omitempty"`       // 公開日時
}

type UpdateVideoRequest struct {
	Title             string   `json:"title,omitempty"`             // タイトル
	Description       string   `json:"description,omitempty"`       // 説明
	CategoryIDs       []string `json:"categoryIds,omitempty"`       // カテゴリID一覧
	ProductIDs        []string `json:"productIds,omitempty"`        // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds,omitempty"`     // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl,omitempty"`      // サムネイルURL
	VideoURL          string   `json:"videoUrl,omitempty"`          // 動画URL
	Public            bool     `json:"public,omitempty"`            // 公開設定
	Limited           bool     `json:"limited,omitempty"`           // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct,omitempty"`    // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience,omitempty"` // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt,omitempty"`       // 公開日時
}
