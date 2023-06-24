package response

// Product - 商品情報
type Product struct {
	ID                string          `json:"id"`                 // 商品ID
	ProducerID        string          `json:"producerId"`         // 生産者ID
	ProducerName      string          `json:"producerName"`       // 農家名
	CategoryID        string          `json:"categoryId"`         // 商品種別ID
	CategoryName      string          `json:"categoryName"`       // 商品種別名
	TypeID            string          `json:"productTypeId"`      // 品目ID
	TypeName          string          `json:"productTypeName"`    // 品目名
	TypeIconURL       string          `json:"productTypeIconUrl"` // 品目アイコンURL
	TagIDs            []string        `json:"productTagIds"`      // 商品タグID一覧
	Name              string          `json:"name"`               // 商品名
	Description       string          `json:"description"`        // 商品説明
	Public            bool            `json:"public"`             // 公開フラグ
	Inventory         int64           `json:"inventory"`          // 在庫数
	Weight            float64         `json:"weight"`             // 重量(kg,少数第一位まで)
	ItemUnit          string          `json:"itemUnit"`           // 数量単位
	ItemDescription   string          `json:"itemDescription"`    // 数量単位説明
	Media             []*ProductMedia `json:"media"`              // メディア一覧
	Price             int64           `json:"price"`              // 販売価格
	Cost              int64           `json:"cost"`               // 原価
	ExpirationDate    int64           `json:"expirationDate"`     // 賞味期限(単位:日)
	RecommendedPoint1 string          `json:"recommendedPoint1"`  // おすすめポイント1
	RecommendedPoint2 string          `json:"recommendedPoint2"`  // おすすめポイント2
	RecommendedPoint3 string          `json:"recommendedPoint3"`  // おすすめポイント3
	StorageMethodType int32           `json:"storageMethodType"`  // 保存方法
	DeliveryType      int32           `json:"deliveryType"`       // 配送方法
	Box60Rate         int64           `json:"box60Rate"`          // 箱の占有率(サイズ:60)
	Box80Rate         int64           `json:"box80Rate"`          // 箱の占有率(サイズ:80)
	Box100Rate        int64           `json:"box100Rate"`         // 箱の占有率(サイズ:100)
	OriginPrefecture  string          `json:"originPrefecture"`   // 原産地(都道府県)
	OriginCity        string          `json:"originCity"`         // 原産地(市区町村)
	CreatedAt         int64           `json:"createdAt"`          // 登録日時
	UpdatedAt         int64           `json:"updatedAt"`          // 更新日時
}

// ProductMedia - 商品メディア情報
type ProductMedia struct {
	URL         string   `json:"url"`         // メディアURL
	IsThumbnail bool     `json:"isThumbnail"` // サムネイルとして使用
	Images      []*Image `json:"images"`      // リサイズ済み画像一覧
}

type ProductResponse struct {
	*Product
}

type ProductsResponse struct {
	Products []*Product `json:"products"` // 商品情報一覧
	Total    int64      `json:"total"`    // 商品合計数
}
