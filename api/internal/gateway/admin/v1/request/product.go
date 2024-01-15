package request

type CreateProductRequest struct {
	Name                 string                `json:"name,omitempty"`                 // 商品名
	Description          string                `json:"description,omitempty"`          // 商品説明
	Public               bool                  `json:"public,omitempty"`               // 公開フラグ
	CoordinatorID        string                `json:"coordinatorId,omitempty"`        // コーディネータID
	ProducerID           string                `json:"producerId,omitempty"`           // 生産者ID
	TypeID               string                `json:"productTypeId,omitempty"`        // 品目ID
	TagIDs               []string              `json:"productTagIds,omitempty"`        // 商品タグID一覧
	Inventory            int64                 `json:"inventory,omitempty"`            // 在庫数
	Weight               float64               `json:"weight,omitempty"`               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit,omitempty"`             // 数量単位
	ItemDescription      string                `json:"itemDescription,omitempty"`      // 数量単位説明
	Media                []*CreateProductMedia `json:"media,omitempty"`                // メディア一覧
	Price                int64                 `json:"price,omitempty"`                // 販売価格(税込)
	Cost                 int64                 `json:"cost,omitempty"`                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate,omitempty"`       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1,omitempty"`    // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2,omitempty"`    // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3,omitempty"`    // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType,omitempty"`    // 保存方法
	DeliveryType         int32                 `json:"deliveryType,omitempty"`         // 配送方法
	Box60Rate            int64                 `json:"box60Rate,omitempty"`            // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate,omitempty"`            // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate,omitempty"`           // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode,omitempty"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity,omitempty"`           // 原産地(市区町村)
	StartAt              int64                 `json:"startAt,omitempty"`              // 販売開始日時
	EndAt                int64                 `json:"endAt,omitempty"`                // 販売終了日時
}

type CreateProductMedia struct {
	URL         string `json:"url,omitempty"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail,omitempty"` // サムネイルとして使用
}

type UpdateProductRequest struct {
	Name                 string                `json:"name,omitempty"`                 // 商品名
	Description          string                `json:"description,omitempty"`          // 商品説明
	Public               bool                  `json:"public,omitempty"`               // 公開フラグ
	TypeID               string                `json:"productTypeId,omitempty"`        // 品目ID
	TagIDs               []string              `json:"productTagIds,omitempty"`        // 商品タグID一覧
	Inventory            int64                 `json:"inventory,omitempty"`            // 在庫数
	Weight               float64               `json:"weight,omitempty"`               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit,omitempty"`             // 数量単位
	ItemDescription      string                `json:"itemDescription,omitempty"`      // 数量単位説明
	Media                []*UpdateProductMedia `json:"media,omitempty"`                // メディア一覧
	Price                int64                 `json:"price,omitempty"`                // 販売価格(税込)
	Cost                 int64                 `json:"cost,omitempty"`                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate,omitempty"`       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1,omitempty"`    // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2,omitempty"`    // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3,omitempty"`    // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType,omitempty"`    // 保存方法
	DeliveryType         int32                 `json:"deliveryType,omitempty"`         // 配送方法
	Box60Rate            int64                 `json:"box60Rate,omitempty"`            // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate,omitempty"`            // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate,omitempty"`           // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode,omitempty"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity,omitempty"`           // 原産地(市区町村)
	StartAt              int64                 `json:"startAt,omitempty"`              // 販売開始日時
	EndAt                int64                 `json:"endAt,omitempty"`                // 販売終了日時
}

type UpdateProductMedia struct {
	URL         string `json:"url,omitempty"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail,omitempty"` // サムネイルとして使用
}
