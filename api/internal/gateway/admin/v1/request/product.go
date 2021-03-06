package request

type CreateProductRequest struct {
	Name             string                `json:"name,omitempty"`             // 商品名
	Description      string                `json:"description,omitempty"`      // 商品説明
	Public           bool                  `json:"public,omitempty"`           // 公開フラグ
	ProducerID       string                `json:"producerID,omitempty"`       // 生産者ID
	CategoryID       string                `json:"categoryID,omitempty"`       // 商品種別ID
	TypeID           string                `json:"typeID,omitempty"`           // 品目ID
	Inventory        int64                 `json:"inventory,omitempty"`        // 在庫数
	Weight           float64               `json:"weight,omitempty"`           // 重量(kg,少数第一位まで)
	ItemUnit         string                `json:"itemUnit,omitempty"`         // 数量単位
	ItemDescription  string                `json:"itemDescription,omitempty"`  // 数量単位説明
	Media            []*CreateProductMedia `json:"media,omitempty"`            // メディア一覧
	Price            int64                 `json:"price,omitempty"`            // 販売価格
	DeliveryType     int32                 `json:"deliveryType,omitempty"`     // 配送方法
	Box60Rate        int64                 `json:"box60Rate,omitempty"`        // 箱の占有率(サイズ:60)
	Box80Rate        int64                 `json:"box80Rate,omitempty"`        // 箱の占有率(サイズ:80)
	Box100Rate       int64                 `json:"box100Rate,omitempty"`       // 箱の占有率(サイズ:100)
	OriginPrefecture string                `json:"originPrefecture,omitempty"` // 原産地(都道府県)
	OriginCity       string                `json:"originCity,omitempty"`       // 原産地(市区町村)
}

type CreateProductMedia struct {
	URL         string `json:"url,omitempty"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail,omitempty"` // サムネイルとして使用
}

type UpdateProductRequest struct {
	Name             string                `json:"name,omitempty"`             // 商品名
	Description      string                `json:"description,omitempty"`      // 商品説明
	Public           bool                  `json:"public,omitempty"`           // 公開フラグ
	ProducerID       string                `json:"producerID,omitempty"`       // 生産者ID
	CategoryID       string                `json:"categoryID,omitempty"`       // 商品種別ID
	TypeID           string                `json:"typeID,omitempty"`           // 品目ID
	Inventory        int64                 `json:"inventory,omitempty"`        // 在庫数
	Weight           float64               `json:"weight,omitempty"`           // 重量(kg,少数第一位まで)
	ItemUnit         string                `json:"itemUnit,omitempty"`         // 数量単位
	ItemDescription  string                `json:"itemDescription,omitempty"`  // 数量単位説明
	Media            []*UpdateProductMedia `json:"media,omitempty"`            // メディア一覧
	Price            int64                 `json:"price,omitempty"`            // 販売価格
	DeliveryType     int32                 `json:"deliveryType,omitempty"`     // 配送方法
	Box60Rate        int64                 `json:"box60Rate,omitempty"`        // 箱の占有率(サイズ:60)
	Box80Rate        int64                 `json:"box80Rate,omitempty"`        // 箱の占有率(サイズ:80)
	Box100Rate       int64                 `json:"box100Rate,omitempty"`       // 箱の占有率(サイズ:100)
	OriginPrefecture string                `json:"originPrefecture,omitempty"` // 原産地(都道府県)
	OriginCity       string                `json:"originCity,omitempty"`       // 原産地(市区町村)
}

type UpdateProductMedia struct {
	URL         string `json:"url,omitempty"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail,omitempty"` // サムネイルとして使用
}
