package types

// ProductStatus - 商品販売状況
type ProductStatus int32

const (
	ProductStatusUnknown   ProductStatus = 0
	ProductStatusPrivate   ProductStatus = 1 // 非公開
	ProductStatusPresale   ProductStatus = 2 // 予約受付中
	ProductStatusForSale   ProductStatus = 3 // 販売中
	ProductStatusOutOfSale ProductStatus = 4 // 販売期間外
	ProductStatusArchived  ProductStatus = 5 // アーカイブ済み
)

// StorageMethodType - 保存方法
type StorageMethodType int32

const (
	StorageMethodTypeUnknown      StorageMethodType = 0
	StorageMethodTypeNormal       StorageMethodType = 1 // 常温保存
	StorageMethodTypeCoolDark     StorageMethodType = 2 // 冷暗所保存
	StorageMethodTypeRefrigerated StorageMethodType = 3 // 冷蔵保存
	StorageMethodTypeFrozen       StorageMethodType = 4 // 冷凍保存
)

// DeliveryType - 配送方法
type DeliveryType int32

const (
	DeliveryTypeUnknown      DeliveryType = 0
	DeliveryTypeNormal       DeliveryType = 1 // 常温便
	DeliveryTypeRefrigerated DeliveryType = 2 // 冷蔵便
	DeliveryTypeFrozen       DeliveryType = 3 // 冷凍便
)

// Product - 商品情報
type Product struct {
	ID                   string            `json:"id"`                   // 商品ID
	CoordinatorID        string            `json:"coordinatorId"`        // コーディネータID
	ProducerID           string            `json:"producerId"`           // 生産者ID
	CategoryID           string            `json:"categoryId"`           // 商品種別ID
	ProductTypeID        string            `json:"productTypeId"`        // 品目ID
	ProductTagIDs        []string          `json:"productTagIds"`        // 商品タグID一覧
	Name                 string            `json:"name"`                 // 商品名
	Description          string            `json:"description"`          // 商品説明
	Public               bool              `json:"public"`               // 公開フラグ
	Status               ProductStatus     `json:"status"`               // 販売状況
	Inventory            int64             `json:"inventory"`            // 在庫数
	Weight               float64           `json:"weight"`               // 重量(kg,少数第一位まで)
	ItemUnit             string            `json:"itemUnit"`             // 数量単位
	ItemDescription      string            `json:"itemDescription"`      // 数量単位説明
	Media                []*ProductMedia   `json:"media"`                // メディア一覧
	Price                int64             `json:"price"`                // 販売価格(税込)
	Cost                 int64             `json:"cost"`                 // 原価
	ExpirationDate       int64             `json:"expirationDate"`       // 賞味期限(単位:日)
	RecommendedPoint1    string            `json:"recommendedPoint1"`    // おすすめポイント1
	RecommendedPoint2    string            `json:"recommendedPoint2"`    // おすすめポイント2
	RecommendedPoint3    string            `json:"recommendedPoint3"`    // おすすめポイント3
	StorageMethodType    StorageMethodType `json:"storageMethodType"`    // 保存方法
	DeliveryType         DeliveryType      `json:"deliveryType"`         // 配送方法
	Box60Rate            int64             `json:"box60Rate"`            // 箱の占有率(サイズ:60)
	Box80Rate            int64             `json:"box80Rate"`            // 箱の占有率(サイズ:80)
	Box100Rate           int64             `json:"box100Rate"`           // 箱の占有率(サイズ:100)
	OriginPrefectureCode int32             `json:"originPrefectureCode"` // 原産地(都道府県)
	OriginCity           string            `json:"originCity"`           // 原産地(市区町村)
	StartAt              int64             `json:"startAt"`              // 販売開始日時
	EndAt                int64             `json:"endAt"`                // 販売終了日時
	CreatedAt            int64             `json:"createdAt"`            // 登録日時
	UpdatedAt            int64             `json:"updatedAt"`            // 更新日時
}

// ProductMedia - 商品メディア情報
type ProductMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

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
