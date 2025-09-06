package request

type CreateProductRequest struct {
	Name                 string                `json:"name" binding:"required,max=64"`                       // 商品名
	Description          string                `json:"description" binding:"required,max=2000"`              // 商品説明
	Public               bool                  `json:"public"`                                               // 公開フラグ
	CoordinatorID        string                `json:"coordinatorId" binding:"required"`                     // コーディネータID
	ProducerID           string                `json:"producerId" binding:"required"`                        // 生産者ID
	TypeID               string                `json:"productTypeId" binding:"required"`                     // 品目ID
	TagIDs               []string              `json:"productTagIds" binding:"max=8,dive,required"`          // 商品タグID一覧
	Inventory            int64                 `json:"inventory" binding:"min=0"`                            // 在庫数
	Weight               float64               `json:"weight" binding:"min=0"`                               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit" binding:"required,max=16"`                   // 数量単位
	ItemDescription      string                `json:"itemDescription" binding:"required,max=64"`            // 数量単位説明
	Media                []*CreateProductMedia `json:"media" binding:"max=8,dive"`                           // メディア一覧
	Price                int64                 `json:"price" binding:"min=0"`                                // 販売価格(税込)
	Cost                 int64                 `json:"cost" binding:"min=0"`                                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate" binding:"min=0"`                       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1" binding:"omitempty,max=128"`        // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2" binding:"omitempty,max=128"`        // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3" binding:"omitempty,max=128"`        // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType" binding:"required"`                 // 保存方法
	DeliveryType         int32                 `json:"deliveryType" binding:"required"`                      // 配送方法
	Box60Rate            int64                 `json:"box60Rate" binding:"min=0,max=600"`                    // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate" binding:"min=0,max=250"`                    // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate" binding:"min=0,max=100"`                   // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode" binding:"required,min=1,max=47"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity" binding:"max=32"`                          // 原産地(市区町村)
	StartAt              int64                 `json:"startAt" binding:"required"`                           // 販売開始日時
	EndAt                int64                 `json:"endAt" binding:"required,gtfield=StartAt"`             // 販売終了日時
}

type CreateProductMedia struct {
	URL         string `json:"url" binding:"required,url"` // メディアURL
	IsThumbnail bool   `json:"isThumbnail"`                // サムネイルとして使用
}

type UpdateProductRequest struct {
	Name                 string                `json:"name" binding:"required,max=64"`                       // 商品名
	Description          string                `json:"description" binding:"required,max=2000"`              // 商品説明
	Public               bool                  `json:"public"`                                               // 公開フラグ
	TypeID               string                `json:"productTypeId" binding:"required"`                     // 品目ID
	TagIDs               []string              `json:"productTagIds" binding:"max=8,dive,required"`          // 商品タグID一覧
	Inventory            int64                 `json:"inventory" binding:"min=0"`                            // 在庫数
	Weight               float64               `json:"weight" binding:"min=0"`                               // 重量(kg,少数第一位まで)
	ItemUnit             string                `json:"itemUnit" binding:"required,max=16"`                   // 数量単位
	ItemDescription      string                `json:"itemDescription" binding:"required,max=64"`            // 数量単位説明
	Media                []*UpdateProductMedia `json:"media" binding:"max=8,dive"`                           // メディア一覧
	Price                int64                 `json:"price" binding:"min=0"`                                // 販売価格(税込)
	Cost                 int64                 `json:"cost" binding:"min=0"`                                 // 原価(税込)
	ExpirationDate       int64                 `json:"expirationDate" binding:"min=0"`                       // 賞味期限(単位:日)
	RecommendedPoint1    string                `json:"recommendedPoint1" binding:"omitempty,max=128"`        // おすすめポイント1
	RecommendedPoint2    string                `json:"recommendedPoint2" binding:"omitempty,max=128"`        // おすすめポイント2
	RecommendedPoint3    string                `json:"recommendedPoint3" binding:"omitempty,max=128"`        // おすすめポイント3
	StorageMethodType    int32                 `json:"storageMethodType" binding:"required"`                 // 保存方法
	DeliveryType         int32                 `json:"deliveryType" binding:"required"`                      // 配送方法
	Box60Rate            int64                 `json:"box60Rate" binding:"min=0,max=600"`                    // 箱の占有率(サイズ:60)
	Box80Rate            int64                 `json:"box80Rate" binding:"min=0,max=250"`                    // 箱の占有率(サイズ:80)
	Box100Rate           int64                 `json:"box100Rate" binding:"min=0,max=100"`                   // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32                 `json:"originPrefectureCode" binding:"required,min=1,max=47"` // 原産地(都道府県)
	OriginCity           string                `json:"originCity" binding:"max=32"`                          // 原産地(市区町村)
	StartAt              int64                 `json:"startAt" binding:"required"`                           // 販売開始日時
	EndAt                int64                 `json:"endAt" binding:"required,gtfield=StartAt"`             // 販売終了日時
}

type UpdateProductMedia struct {
	URL         string `json:"url" binding:"required,url"` // メディアURL
	IsThumbnail bool   `json:"isThumbnail"`                // サムネイルとして使用
}
