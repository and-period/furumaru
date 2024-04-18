package response

// Product - 商品情報
type Product struct {
	ID                string          `json:"id"`                // 商品ID
	CoordinatorID     string          `json:"coordinatorId"`     // コーディネータID
	ProducerID        string          `json:"producerId"`        // 生産者ID
	CategoryID        string          `json:"categoryId"`        // 商品種別ID
	ProductTypeID     string          `json:"productTypeId"`     // 品目ID
	ProductTagIDs     []string        `json:"productTagIds"`     // 商品タグID一覧
	Name              string          `json:"name"`              // 商品名
	Description       string          `json:"description"`       // 商品説明
	Status            int32           `json:"status"`            // 販売状況
	Inventory         int64           `json:"inventory"`         // 在庫数
	Weight            float64         `json:"weight"`            // 重量(kg,少数第一位まで)
	ItemUnit          string          `json:"itemUnit"`          // 数量単位
	ItemDescription   string          `json:"itemDescription"`   // 数量単位説明
	ThumbnailURL      string          `json:"thumbnailUrl"`      // サムネイルURL
	Media             []*ProductMedia `json:"media"`             // メディア一覧
	Price             int64           `json:"price"`             // 販売価格(税込)
	ExpirationDate    int64           `json:"expirationDate"`    // 賞味期限(単位:日)
	RecommendedPoint1 string          `json:"recommendedPoint1"` // おすすめポイント1
	RecommendedPoint2 string          `json:"recommendedPoint2"` // おすすめポイント2
	RecommendedPoint3 string          `json:"recommendedPoint3"` // おすすめポイント3
	StorageMethodType int32           `json:"storageMethodType"` // 保存方法
	DeliveryType      int32           `json:"deliveryType"`      // 配送方法
	Box60Rate         int64           `json:"box60Rate"`         // 箱の占有率(サイズ:60)
	Box80Rate         int64           `json:"box80Rate"`         // 箱の占有率(サイズ:80)
	Box100Rate        int64           `json:"box100Rate"`        // 箱の占有率(サイズ:100)
	OriginPrefecture  string          `json:"originPrefecture"`  // 原産地(都道府県)
	OriginCity        string          `json:"originCity"`        // 原産地(市区町村)
	StartAt           int64           `json:"startAt"`           // 販売開始日時
	EndAt             int64           `json:"endAt"`             // 販売終了日時
}

// ProductMedia - 商品メディア情報
type ProductMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

type ProductResponse struct {
	Product     *Product      `json:"product"`     // 商品情報
	Coordinator *Coordinator  `json:"coordinator"` // コーディネータ情報
	Producer    *Producer     `json:"producer"`    // 生産者情報
	Category    *Category     `json:"category"`    // 商品種別情報
	ProductType *ProductType  `json:"productType"` // 品目情報
	ProductTags []*ProductTag `json:"productTags"` // 商品タグ一覧
}

type ProductsResponse struct {
	Products     []*Product     `json:"products"`     // 商品一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Producers    []*Producer    `json:"producers"`    // 生産者一覧
	Categories   []*Category    `json:"categories"`   // 商品種別一覧
	ProductTypes []*ProductType `json:"productTypes"` // 品目一覧
	ProductTags  []*ProductTag  `json:"productTags"`  // 商品タグ一覧
	Total        int64          `json:"total"`        // 商品合計数
}
