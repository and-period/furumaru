package request

type CreateProductRequest struct {
	Name                 string                `json:"name" validate:"required,max=64"`                       // 商品名
	Description          string                `json:"description" validate:"required,max=2000"`              // 商品説明
	Public               bool                  `json:"public"`                                                // 公開フラグ
	CoordinatorID        string                `json:"coordinatorId" validate:"required"`                     // コーディネータID
	ProducerID           string                `json:"producerId" validate:"required"`                        // 生産者ID
	TypeID               string                `json:"productTypeId" validate:"required"`                     // 品目ID
	TagIDs               []string              `json:"productTagIds" validate:"max=8,dive,required"`          // 商品タグID一覧
	Inventory            int64                 `json:"inventory" validate:"min=0"`                            // 在庫数
	Weight               float64               `json:"weight" validate:"min=0"`                               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit" validate:"required,max=16"`                   // 数量単位
	ItemDescription      string                `json:"itemDescription" validate:"required,max=64"`            // 数量単位説明
	Media                []*CreateProductMedia `json:"media" validate:"max=8,dive"`                           // メディア一覧
	Price                int64                 `json:"price" validate:"min=0"`                                // 販売価格(税込)
	Cost                 int64                 `json:"cost" validate:"min=0"`                                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate" validate:"min=0"`                       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1" validate:"omitempty,max=128"`        // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2" validate:"omitempty,max=128"`        // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3" validate:"omitempty,max=128"`        // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType" validate:"required"`                 // 保存方法
	DeliveryType         int32                 `json:"deliveryType" validate:"required"`                      // 配送方法
	Box60Rate            int64                 `json:"box60Rate" validate:"min=0,max=600"`                    // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate" validate:"min=0,max=250"`                    // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate" validate:"min=0,max=100"`                   // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode" validate:"required,min=1,max=47"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity" validate:"max=32"`                          // 原産地(市区町村)
	StartAt              int64                 `json:"startAt" validate:"required"`                           // 販売開始日時
	EndAt                int64                 `json:"endAt" validate:"required,gtfield=StartAt"`             // 販売終了日時
}

type CreateProductMedia struct {
	URL         string `json:"url" validate:"required,url"` // メディアURL
	IsThumbnail bool   `json:"isThumbnail"`                 // サムネイルとして使用
}

type UpdateProductRequest struct {
	Name                 string                `json:"name" validate:"required,max=64"`                       // 商品名
	Description          string                `json:"description" validate:"required,max=2000"`              // 商品説明
	Public               bool                  `json:"public"`                                                // 公開フラグ
	TypeID               string                `json:"productTypeId" validate:"required"`                     // 品目ID
	TagIDs               []string              `json:"productTagIds" validate:"max=8,dive,required"`          // 商品タグID一覧
	Inventory            int64                 `json:"inventory" validate:"min=0"`                            // 在庫数
	Weight               float64               `json:"weight" validate:"min=0"`                               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit" validate:"required,max=16"`                   // 数量単位
	ItemDescription      string                `json:"itemDescription" validate:"required,max=64"`            // 数量単位説明
	Media                []*UpdateProductMedia `json:"media" validate:"max=8,dive"`                           // メディア一覧
	Price                int64                 `json:"price" validate:"min=0"`                                // 販売価格(税込)
	Cost                 int64                 `json:"cost" validate:"min=0"`                                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate" validate:"min=0"`                       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1" validate:"omitempty,max=128"`        // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2" validate:"omitempty,max=128"`        // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3" validate:"omitempty,max=128"`        // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType" validate:"required"`                 // 保存方法
	DeliveryType         int32                 `json:"deliveryType" validate:"required"`                      // 配送方法
	Box60Rate            int64                 `json:"box60Rate" validate:"min=0,max=600"`                    // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate" validate:"min=0,max=250"`                    // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate" validate:"min=0,max=100"`                   // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode" validate:"required,min=1,max=47"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity" validate:"max=32"`                          // 原産地(市区町村)
	StartAt              int64                 `json:"startAt" validate:"required"`                           // 販売開始日時
	EndAt                int64                 `json:"endAt" validate:"required,gtfield=StartAt"`             // 販売終了日時
}

type UpdateProductMedia struct {
	URL         string `json:"url" validate:"required,url"` // メディアURL
	IsThumbnail bool   `json:"isThumbnail"`                 // サムネイルとして使用
}
