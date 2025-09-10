package types

// ExperienceStatus - 体験受付状況
type ExperienceStatus int32

const (
	ExperienceStatusUnknown   ExperienceStatus = 0
	ExperienceStatusWaiting   ExperienceStatus = 1 // 販売開始前
	ExperienceStatusAccepting ExperienceStatus = 2 // 体験受付中
	ExperienceStatusSoldOut   ExperienceStatus = 3 // 体験受付終了
	ExperienceStatusFinished  ExperienceStatus = 4 // 販売終了
)

// Experience - 体験情報
type Experience struct {
	ID                    string             `json:"id"`                    // 体験ID
	CoordinatorID         string             `json:"coordinatorId"`         // コーディネータID
	ProducerID            string             `json:"producerId"`            // プロデューサーID
	ExperienceTypeID      string             `json:"experienceTypeId"`      // 体験種別ID
	Title                 string             `json:"title"`                 // タイトル
	Description           string             `json:"description"`           // 説明
	Status                ExperienceStatus   `json:"status"`                // 販売状況
	ThumbnailURL          string             `json:"thumbnailUrl"`          // サムネイルURL
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
	HostPrefecture        string             `json:"hostPrefecture"`        // 開催場所(都道府県)
	HostCity              string             `json:"hostCity"`              // 開催場所(市区町村)
	HostAddressLine1      string             `json:"hostAddressLine1"`      // 開催場所(住所1)
	HostAddressLine2      string             `json:"hostAddressLine2"`      // 開催場所(住所2)
	HostLongitude         float64            `json:"hostLongitude"`         // 開催場所(座標情報:経度)
	HostLatitude          float64            `json:"hostLatitude"`          // 開催場所(座標情報:緯度)
	Rate                  *ExperienceRate    `json:"rate"`                  // 体験評価
	StartAt               int64              `json:"startAt"`               // 募集開始日時
	EndAt                 int64              `json:"endAt"`                 // 募集終了日時
}

// ExperienceMedia - 体験メディア情報
type ExperienceMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

// ExperienceRate - 体験評価情報
type ExperienceRate struct {
	Average float64         `json:"average"` // 平均評価
	Count   int64           `json:"count"`   // 合計評価数
	Detail  map[int64]int64 `json:"detail"`  // 評価詳細
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
