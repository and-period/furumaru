package request

type CreateVideoRequest struct {
	Title             string   `json:"title"`             // タイトル
	Description       string   `json:"description"`       // 説明
	CoordinatorID     string   `json:"coordinatorId"`     // コーディネータID
	ProductIDs        []string `json:"productIds"`        // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds"`     // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl"`      // サムネイルURL
	VideoURL          string   `json:"videoUrl"`          // 動画URL
	Public            bool     `json:"public"`            // 公開設定
	Limited           bool     `json:"limited"`           // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct"`    // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience"` // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt"`       // 公開日時
}

type UpdateVideoRequest struct {
	Title             string   `json:"title"`             // タイトル
	Description       string   `json:"description"`       // 説明
	CategoryIDs       []string `json:"categoryIds"`       // カテゴリID一覧
	ProductIDs        []string `json:"productIds"`        // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds"`     // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl"`      // サムネイルURL
	VideoURL          string   `json:"videoUrl"`          // 動画URL
	Public            bool     `json:"public"`            // 公開設定
	Limited           bool     `json:"limited"`           // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct"`    // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience"` // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt"`       // 公開日時
}
