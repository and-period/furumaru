package request

type CreateExperienceRequest struct {
	Title                 string                   `json:"title,omitempty"`                 // 体験名
	Description           string                   `json:"description,omitempty"`           // 説明
	Public                bool                     `json:"public,omitempty"`                // 公開設定
	SoldOut               bool                     `json:"soldOut,omitempty"`               // 定員オーバーフラグ
	CoordinatorID         string                   `json:"coordinatorId,omitempty"`         // コーディネータID
	ProducerID            string                   `json:"producerId,omitempty"`            // 生産者ID
	TypeID                string                   `json:"experienceTypeId,omitempty"`      // 体験種別ID
	Media                 []*CreateExperienceMedia `json:"media,omitempty"`                 // メディア一覧
	PriceAdult            int64                    `json:"priceAdult,omitempty"`            // 大人料金
	PriceJuniorHighSchool int64                    `json:"priceJuniorHighSchool,omitempty"` // 中学生料金
	PriceElementarySchool int64                    `json:"priceElementarySchool,omitempty"` // 小学生料金
	PricePreschooler      int64                    `json:"pricePreschooler,omitempty"`      // 幼児料金
	PriceSenior           int64                    `json:"priceSenior,omitempty"`           // シニア料金
	RecommendedPoint1     string                   `json:"recommendedPoint1,omitempty"`     // おすすめポイント1
	RecommendedPoint2     string                   `json:"recommendedPoint2,omitempty"`     // おすすめポイント2
	RecommendedPoint3     string                   `json:"recommendedPoint3,omitempty"`     // おすすめポイント3
	PromotionVideoURL     string                   `json:"promotionVideoUrl,omitempty"`     // 紹介動画URL
	HostPrefectureCode    int32                    `json:"hostPrefectureCode,omitempty"`    // 開催場所(都道府県コード)
	HostCity              string                   `json:"hostCity,omitempty"`              // 開催場所(市区町村)
	StartAt               int64                    `json:"startAt,omitempty"`               // 募集開始日時
	EndAt                 int64                    `json:"endAt,omitempty"`                 // 募集終了日時
}

type CreateExperienceMedia struct {
	URL         string `json:"url,omitempty"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail,omitempty"` // サムネイルとして使用
}

type UpdateExperienceRequest struct {
	Title                 string                   `json:"title,omitempty"`                 // 体験名
	Description           string                   `json:"description,omitempty"`           // 説明
	Public                bool                     `json:"public,omitempty"`                // 公開設定
	SoldOut               bool                     `json:"soldOut,omitempty"`               // 定員オーバーフラグ
	CoordinatorID         string                   `json:"coordinatorId,omitempty"`         // コーディネータID
	ProducerID            string                   `json:"producerId,omitempty"`            // 生産者ID
	TypeID                string                   `json:"experienceTypeId,omitempty"`      // 体験種別ID
	Media                 []*UpdateExperienceMedia `json:"media,omitempty"`                 // メディア一覧
	PriceAdult            int64                    `json:"priceAdult,omitempty"`            // 大人料金
	PriceJuniorHighSchool int64                    `json:"priceJuniorHighSchool,omitempty"` // 中学生料金
	PriceElementarySchool int64                    `json:"priceElementarySchool,omitempty"` // 小学生料金
	PricePreschooler      int64                    `json:"pricePreschooler,omitempty"`      // 幼児料金
	PriceSenior           int64                    `json:"priceSenior,omitempty"`           // シニア料金
	RecommendedPoint1     string                   `json:"recommendedPoint1,omitempty"`     // おすすめポイント1
	RecommendedPoint2     string                   `json:"recommendedPoint2,omitempty"`     // おすすめポイント2
	RecommendedPoint3     string                   `json:"recommendedPoint3,omitempty"`     // おすすめポイント3
	PromotionVideoURL     string                   `json:"promotionVideoUrl,omitempty"`     // 紹介動画URL
	HostPrefectureCode    int32                    `json:"hostPrefectureCode,omitempty"`    // 開催場所(都道府県コード)
	HostCity              string                   `json:"hostCity,omitempty"`              // 開催場所(市区町村)
	StartAt               int64                    `json:"startAt,omitempty"`               // 募集開始日時
	EndAt                 int64                    `json:"endAt,omitempty"`                 // 募集終了日時
}

type UpdateExperienceMedia struct {
	URL         string `json:"url,omitempty"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail,omitempty"` // サムネイルとして使用
}
