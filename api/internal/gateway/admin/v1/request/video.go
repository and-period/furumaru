package request

type CreateVideoRequest struct {
	Title             string   `json:"title" binding:"required,max=128"`        // タイトル
	Description       string   `json:"description" binding:"required,max=2000"` // 説明
	CoordinatorID     string   `json:"coordinatorId" binding:"required"`        // コーディネータID
	ProductIDs        []string `json:"productIds" binding:"dive,required"`      // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds" binding:"dive,required"`   // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl" binding:"required,url"`     // サムネイルURL
	VideoURL          string   `json:"videoUrl" binding:"required,url"`         // 動画URL
	Public            bool     `json:"public" binding:""`                       // 公開設定
	Limited           bool     `json:"limited" binding:""`                      // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct" binding:""`               // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience" binding:""`            // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt" binding:"required"`          // 公開日時
}

type UpdateVideoRequest struct {
	Title             string   `json:"title" binding:"required,max=128"`        // タイトル
	Description       string   `json:"description" binding:"required,max=2000"` // 説明
	CoordinatorID     string   `json:"coordinatorId" binding:"required"`        // コーディネータID
	ProductIDs        []string `json:"productIds" binding:"dive,required"`      // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds" binding:"dive,required"`   // 体験ID一覧
	ThumbnailURL      string   `json:"thumbnailUrl" binding:"required,url"`     // サムネイルURL
	VideoURL          string   `json:"videoUrl" binding:"required,url"`         // 動画URL
	Public            bool     `json:"public" binding:""`                       // 公開設定
	Limited           bool     `json:"limited" binding:""`                      // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct" binding:""`               // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience" binding:""`            // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt" binding:"required"`          // 公開日時
}
