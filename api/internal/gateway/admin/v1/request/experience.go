package request

type CreateExperienceRequest struct {
	Title                 string                   `json:"title"`                 // 体験名
	Description           string                   `json:"description"`           // 説明
	Public                bool                     `json:"public"`                // 公開設定
	SoldOut               bool                     `json:"soldOut"`               // 定員オーバーフラグ
	CoordinatorID         string                   `json:"coordinatorId"`         // コーディネータID
	ProducerID            string                   `json:"producerId"`            // 生産者ID
	TypeID                string                   `json:"experienceTypeId"`      // 体験種別ID
	Media                 []*CreateExperienceMedia `json:"media"`                 // メディア一覧
	PriceAdult            int64                    `json:"priceAdult"`            // 大人料金
	PriceJuniorHighSchool int64                    `json:"priceJuniorHighSchool"` // 中学生料金
	PriceElementarySchool int64                    `json:"priceElementarySchool"` // 小学生料金
	PricePreschool        int64                    `json:"pricePreschool"`        // 幼児料金
	PriceSenior           int64                    `json:"priceSenior"`           // シニア料金
	RecommendedPoint1     string                   `json:"recommendedPoint1"`     // おすすめポイント1
	RecommendedPoint2     string                   `json:"recommendedPoint2"`     // おすすめポイント2
	RecommendedPoint3     string                   `json:"recommendedPoint3"`     // おすすめポイント3
	PromotionVideoURL     string                   `json:"promotionVideoUrl"`     // 紹介動画URL
	Duration              int64                    `json:"duration"`              // 体験時間(分)
	Direction             string                   `json:"direction"`             // アクセス方法
	BusinessOpenTime      string                   `json:"businessOpenTime"`      // 営業開始時間
	BusinessCloseTime     string                   `json:"businessCloseTime"`     // 営業終了時間
	HostPostalCode        string                   `json:"hostPostalCode"`        // 開催場所(郵便番号)
	HostPrefectureCode    int32                    `json:"hostPrefectureCode"`    // 開催場所(都道府県コード)
	HostCity              string                   `json:"hostCity"`              // 開催場所(市区町村)
	HostAddressLine1      string                   `json:"hostAddressLine1"`      // 開催場所(住所1)
	HostAddressLine2      string                   `json:"hostAddressLine2"`      // 開催場所(住所2)
	StartAt               int64                    `json:"startAt"`               // 募集開始日時
	EndAt                 int64                    `json:"endAt"`                 // 募集終了日時
}

type CreateExperienceMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

type UpdateExperienceRequest struct {
	Title                 string                   `json:"title"`                 // 体験名
	Description           string                   `json:"description"`           // 説明
	Public                bool                     `json:"public"`                // 公開設定
	SoldOut               bool                     `json:"soldOut"`               // 定員オーバーフラグ
	TypeID                string                   `json:"experienceTypeId"`      // 体験種別ID
	Media                 []*UpdateExperienceMedia `json:"media"`                 // メディア一覧
	PriceAdult            int64                    `json:"priceAdult"`            // 大人料金
	PriceJuniorHighSchool int64                    `json:"priceJuniorHighSchool"` // 中学生料金
	PriceElementarySchool int64                    `json:"priceElementarySchool"` // 小学生料金
	PricePreschool        int64                    `json:"pricePreschool"`        // 幼児料金
	PriceSenior           int64                    `json:"priceSenior"`           // シニア料金
	RecommendedPoint1     string                   `json:"recommendedPoint1"`     // おすすめポイント1
	RecommendedPoint2     string                   `json:"recommendedPoint2"`     // おすすめポイント2
	RecommendedPoint3     string                   `json:"recommendedPoint3"`     // おすすめポイント3
	PromotionVideoURL     string                   `json:"promotionVideoUrl"`     // 紹介動画URL
	Duration              int64                    `json:"duration"`              // 体験時間(分)
	Direction             string                   `json:"direction"`             // アクセス方法
	BusinessOpenTime      string                   `json:"businessOpenTime"`      // 営業開始時間
	BusinessCloseTime     string                   `json:"businessCloseTime"`     // 営業終了時間
	HostPostalCode        string                   `json:"hostPostalCode"`        // 開催場所(郵便番号)
	HostPrefectureCode    int32                    `json:"hostPrefectureCode"`    // 開催場所(都道府県コード)
	HostCity              string                   `json:"hostCity"`              // 開催場所(市区町村)
	HostAddressLine1      string                   `json:"hostAddressLine1"`      // 開催場所(住所1)
	HostAddressLine2      string                   `json:"hostAddressLine2"`      // 開催場所(住所2)
	StartAt               int64                    `json:"startAt"`               // 募集開始日時
	EndAt                 int64                    `json:"endAt"`                 // 募集終了日時
}

type UpdateExperienceMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}
