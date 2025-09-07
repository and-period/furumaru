package request

type CreateExperienceRequest struct {
	Title                 string                   `json:"title" validate:"required,max=128"`                      // 体験名
	Description           string                   `json:"description" validate:"required,max=20000"`              // 説明
	Public                bool                     `json:"public"`                                                 // 公開設定
	SoldOut               bool                     `json:"soldOut"`                                                // 定員オーバーフラグ
	CoordinatorID         string                   `json:"coordinatorId" validate:"required"`                      // コーディネータID
	ProducerID            string                   `json:"producerId" validate:"required"`                         // 生産者ID
	TypeID                string                   `json:"experienceTypeId" validate:"required"`                   // 体験種別ID
	Media                 []*CreateExperienceMedia `json:"media" validate:"required,dive"`                         // メディア一覧
	PriceAdult            int64                    `json:"priceAdult" validate:"required,min=0,max=99"`            // 大人料金
	PriceJuniorHighSchool int64                    `json:"priceJuniorHighSchool" validate:"required,min=0,max=99"` // 中学生料金
	PriceElementarySchool int64                    `json:"priceElementarySchool" validate:"required,min=0,max=99"` // 小学生料金
	PricePreschool        int64                    `json:"pricePreschool" validate:"required,min=0,max=99"`        // 幼児料金
	PriceSenior           int64                    `json:"priceSenior" validate:"required,min=0,max=99"`           // シニア料金
	RecommendedPoint1     string                   `json:"recommendedPoint1" validate:"omitempty,max=128"`         // おすすめポイント1
	RecommendedPoint2     string                   `json:"recommendedPoint2" validate:"omitempty,max=128"`         // おすすめポイント2
	RecommendedPoint3     string                   `json:"recommendedPoint3" validate:"omitempty,max=128"`         // おすすめポイント3
	PromotionVideoURL     string                   `json:"promotionVideoUrl" validate:"omitempty,url"`             // 紹介動画URL
	Duration              int64                    `json:"duration" validate:"min=0"`                              // 体験時間(分)
	Direction             string                   `json:"direction" validate:"max=2000"`                          // アクセス方法
	BusinessOpenTime      string                   `json:"businessOpenTime" validate:"time"`                       // 営業開始時間
	BusinessCloseTime     string                   `json:"businessCloseTime" validate:"time"`                      // 営業終了時間
	HostPostalCode        string                   `json:"hostPostalCode" validate:"required,numeric,max=16"`      // 開催場所(郵便番号)
	HostPrefectureCode    int32                    `json:"hostPrefectureCode" validate:"required,min=1,max=47"`    // 開催場所(都道府県コード)
	HostCity              string                   `json:"hostCity" validate:"required,max=32"`                    // 開催場所(市区町村)
	HostAddressLine1      string                   `json:"hostAddressLine1" validate:"required,max=64"`            // 開催場所(住所1)
	HostAddressLine2      string                   `json:"hostAddressLine2" validate:"omitempty,max=64"`           // 開催場所(住所2)
	StartAt               int64                    `json:"startAt" validate:"required"`                            // 募集開始日時
	EndAt                 int64                    `json:"endAt" validate:"required,gtfield=StartAt"`              // 募集終了日時
}

type CreateExperienceMedia struct {
	URL         string `json:"url" validate:"required,url"` // メディアURL
	IsThumbnail bool   `json:"isThumbnail"`                 // サムネイルとして使用
}

type UpdateExperienceRequest struct {
	Title                 string                   `json:"title" validate:"required,max=128"`                      // 体験名
	Description           string                   `json:"description" validate:"required,max=20000"`              // 説明
	Public                bool                     `json:"public"`                                                 // 公開設定
	SoldOut               bool                     `json:"soldOut"`                                                // 定員オーバーフラグ
	TypeID                string                   `json:"experienceTypeId" validate:"required"`                   // 体験種別ID
	Media                 []*UpdateExperienceMedia `json:"media" validate:"required,dive"`                         // メディア一覧
	PriceAdult            int64                    `json:"priceAdult" validate:"required,min=0,max=99"`            // 大人料金
	PriceJuniorHighSchool int64                    `json:"priceJuniorHighSchool" validate:"required,min=0,max=99"` // 中学生料金
	PriceElementarySchool int64                    `json:"priceElementarySchool" validate:"required,min=0,max=99"` // 小学生料金
	PricePreschool        int64                    `json:"pricePreschool" validate:"required,min=0,max=99"`        // 幼児料金
	PriceSenior           int64                    `json:"priceSenior" validate:"required,min=0,max=99"`           // シニア料金
	RecommendedPoint1     string                   `json:"recommendedPoint1" validate:"omitempty,max=128"`         // おすすめポイント1
	RecommendedPoint2     string                   `json:"recommendedPoint2" validate:"omitempty,max=128"`         // おすすめポイント2
	RecommendedPoint3     string                   `json:"recommendedPoint3" validate:"omitempty,max=128"`         // おすすめポイント3
	PromotionVideoURL     string                   `json:"promotionVideoUrl" validate:"omitempty,url"`             // 紹介動画URL
	Duration              int64                    `json:"duration" validate:"min=0"`                              // 体験時間(分)
	Direction             string                   `json:"direction" validate:"max=2000"`                          // アクセス方法
	BusinessOpenTime      string                   `json:"businessOpenTime" validate:"time"`                       // 営業開始時間
	BusinessCloseTime     string                   `json:"businessCloseTime" validate:"time"`                      // 営業終了時間
	HostPostalCode        string                   `json:"hostPostalCode" validate:"required,numeric,max=16"`      // 開催場所(郵便番号)
	HostPrefectureCode    int32                    `json:"hostPrefectureCode" validate:"required,min=1,max=47"`    // 開催場所(都道府県コード)
	HostCity              string                   `json:"hostCity" validate:"required,max=32"`                    // 開催場所(市区町村)
	HostAddressLine1      string                   `json:"hostAddressLine1" validate:"required,max=64"`            // 開催場所(住所1)
	HostAddressLine2      string                   `json:"hostAddressLine2" validate:"omitempty,max=64"`           // 開催場所(住所2)
	StartAt               int64                    `json:"startAt" validate:"required"`                            // 募集開始日時
	EndAt                 int64                    `json:"endAt" validate:"required,gtfield=StartAt"`              // 募集終了日時
}

type UpdateExperienceMedia struct {
	URL         string `json:"url" validate:"required,url"` // メディアURL
	IsThumbnail bool   `json:"isThumbnail"`                 // サムネイルとして使用
}
