package request

type CreateProductRequest struct {
	Name                 string                `json:"name"`                 // 商品名
	Description          string                `json:"description"`          // 商品説明
	Public               bool                  `json:"public"`               // 公開フラグ
	CoordinatorID        string                `json:"coordinatorId"`        // コーディネータID
	ProducerID           string                `json:"producerId"`           // 生産者ID
	TypeID               string                `json:"productTypeId"`        // 品目ID
	TagIDs               []string              `json:"productTagIds"`        // 商品タグID一覧
	Inventory            int64                 `json:"inventory"`            // 在庫数
	Weight               float64               `json:"weight"`               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit"`             // 数量単位
	ItemDescription      string                `json:"itemDescription"`      // 数量単位説明
	Media                []*CreateProductMedia `json:"media"`                // メディア一覧
	Price                int64                 `json:"price"`                // 販売価格(税込)
	Cost                 int64                 `json:"cost"`                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate"`       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1"`    // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2"`    // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3"`    // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType"`    // 保存方法
	DeliveryType         int32                 `json:"deliveryType"`         // 配送方法
	Box60Rate            int64                 `json:"box60Rate"`            // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate"`            // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate"`           // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity"`           // 原産地(市区町村)
	StartAt              int64                 `json:"startAt"`              // 販売開始日時
	EndAt                int64                 `json:"endAt"`                // 販売終了日時
}

type CreateProductMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

type UpdateProductRequest struct {
	Name                 string                `json:"name"`                 // 商品名
	Description          string                `json:"description"`          // 商品説明
	Public               bool                  `json:"public"`               // 公開フラグ
	TypeID               string                `json:"productTypeId"`        // 品目ID
	TagIDs               []string              `json:"productTagIds"`        // 商品タグID一覧
	Inventory            int64                 `json:"inventory"`            // 在庫数
	Weight               float64               `json:"weight"`               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit"`             // 数量単位
	ItemDescription      string                `json:"itemDescription"`      // 数量単位説明
	Media                []*UpdateProductMedia `json:"media"`                // メディア一覧
	Price                int64                 `json:"price"`                // 販売価格(税込)
	Cost                 int64                 `json:"cost"`                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate"`       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1"`    // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2"`    // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3"`    // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType"`    // 保存方法
	DeliveryType         int32                 `json:"deliveryType"`         // 配送方法
	Box60Rate            int64                 `json:"box60Rate"`            // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate"`            // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate"`           // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity"`           // 原産地(市区町村)
	StartAt              int64                 `json:"startAt"`              // 販売開始日時
	EndAt                int64                 `json:"endAt"`                // 販売終了日時
}

type UpdateProductMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}
