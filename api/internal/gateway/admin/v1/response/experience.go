package response

// Experience - 体験情報
type Experience struct {
	ID                    string             `json:"id"`                    // 体験ID
	CoordinatorID         string             `json:"coordinatorId"`         // コーディネータID
	ProducerID            string             `json:"producerId"`            // プロデューサーID
	ExperienceTypeID      string             `json:"experienceTypeId"`      // 体験種別ID
	Title                 string             `json:"title"`                 // タイトル
	Description           string             `json:"description"`           // 説明
	Public                bool               `json:"public"`                // 公開設定
	SoldOut               bool               `json:"soldOut"`               // 定員オーバー設定
	Status                int32              `json:"status"`                // 販売状況
	Media                 []*ExperienceMedia `json:"media"`                 // メディア一覧
	PriceAdult            int64              `json:"priceAdult"`            // 大人料金
	PriceJuniorHighSchool int64              `json:"priceJuniorHighSchool"` // 中学生料金
	PriceElementarySchool int64              `json:"priceElementarySchool"` // 小学生料金
	PricePreschool        int64              `json:"pricePreschool"`        // 幼児料金
	PriceSenior           int64              `json:"priceSenior"`           // シニア料金
	RecommendedPoint1     string             `json:"recommendedPoint1"`     // おすすめポイント1
	RecommendedPoint2     string             `json:"recommendedPoint2"`     // おすすめポイント2
	RecommendedPoint3     string             `json:"recommendedPoint3"`     // おすすめポイント3
	PromotionVideoURL     string             `json:"promotionVideoUrl"`     // 紹介動画URL
	Duration              int64              `json:"duration"`              // 体験時間(分)
	Direction             string             `json:"direction"`             // アクセス方法
	BusinessOpenTime      string             `json:"businessOpenTime"`      // 営業開始時間
	BusinessCloseTime     string             `json:"businessCloseTime"`     // 営業終了時間
	HostPostalCode        string             `json:"hostPostalCode"`        // 開催場所(郵便番号)
	HostPrefectureCode    int32              `json:"hostPrefectureCode"`    // 開催場所(都道府県コード)
	HostCity              string             `json:"hostCity"`              // 開催場所(市区町村)
	HostAddressLine1      string             `json:"hostAddressLine1"`      // 開催場所(住所1)
	HostAddressLine2      string             `json:"hostAddressLine2"`      // 開催場所(住所2)
	StartAt               int64              `json:"startAt"`               // 募集開始日時
	EndAt                 int64              `json:"endAt"`                 // 募集終了日時
	CreatedAt             int64              `json:"createdAt"`             // 作成日時
	UpdatedAt             int64              `json:"updatedAt"`             // 更新日時
}

// ExperiencesMedia - 体験メディア情報
type ExperienceMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

type ExperienceResponse struct {
	Experience     *Experience     `json:"experience"`     // 体験情報
	Coordinator    *Coordinator    `json:"coordinator"`    // コーディネータ情報
	Producer       *Producer       `json:"producer"`       // 生産者情報
	ExperienceType *ExperienceType `json:"experienceType"` // 体験種別情報
}

type ExperiencesResponse struct {
	Experiences     []*Experience     `json:"experiences"`     // 体験一覧
	Coordinators    []*Coordinator    `json:"coordinators"`    // コーディネータ一覧
	Producers       []*Producer       `json:"producers"`       // 生産者一覧
	ExperienceTypes []*ExperienceType `json:"experienceTypes"` // 体験種別一覧
	Total           int64             `json:"total"`           // 体験合計数
}
