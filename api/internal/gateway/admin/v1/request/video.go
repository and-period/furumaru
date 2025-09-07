package request

type CreateVideoRequest struct {
	Title             string   `json:"title" validate:"required,max=128"`        // タイトル
	Description       string   `json:"description" validate:"required,max=2000"` // 説明
	CoordinatorID     string   `json:"coordinatorId" validate:"required"`        // コーディネータID
	ProductIDs        []string `json:"productIds" validate:"dive,required"`      // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds" validate:"dive,required"`   // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl" validate:"required,url"`     // サムネイルURL
	VideoURL          string   `json:"videoUrl" validate:"required,url"`         // 動画URL
	Public            bool     `json:"public" validate:""`                       // 公開設定
	Limited           bool     `json:"limited" validate:""`                      // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct" validate:""`               // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience" validate:""`            // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt" validate:"required"`          // 公開日時
}

type UpdateVideoRequest struct {
	Title             string   `json:"title" validate:"required,max=128"`        // タイトル
	Description       string   `json:"description" validate:"required,max=2000"` // 説明
	CoordinatorID     string   `json:"coordinatorId" validate:"required"`        // コーディネータID
	ProductIDs        []string `json:"productIds" validate:"dive,required"`      // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds" validate:"dive,required"`   // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl" validate:"required,url"`     // サムネイルURL
	VideoURL          string   `json:"videoUrl" validate:"required,url"`         // 動画URL
	Public            bool     `json:"public" validate:""`                       // 公開設定
	Limited           bool     `json:"limited" validate:""`                      // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct" validate:""`               // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience" validate:""`            // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt" validate:"required"`          // 公開日時
}
