package types

import "encoding/xml"

// ProductStatus - 商品販売状況
type ProductStatus int32

const (
	ProductStatusUnknown   ProductStatus = 0
	ProductStatusPresale   ProductStatus = 1 // 予約受付中
	ProductStatusForSale   ProductStatus = 2 // 販売中
	ProductStatusOutOfSale ProductStatus = 3 // 販売期間外
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
	ID                string            `json:"id"`                // 商品ID
	CoordinatorID     string            `json:"coordinatorId"`     // コーディネータID
	ProducerID        string            `json:"producerId"`        // 生産者ID
	CategoryID        string            `json:"categoryId"`        // 商品種別ID
	ProductTypeID     string            `json:"productTypeId"`     // 品目ID
	ProductTagIDs     []string          `json:"productTagIds"`     // 商品タグID一覧
	Name              string            `json:"name"`              // 商品名
	Description       string            `json:"description"`       // 商品説明
	Status            ProductStatus     `json:"status"`            // 販売状況
	Inventory         int64             `json:"inventory"`         // 在庫数
	Weight            float64           `json:"weight"`            // 重量(kg,少数第一位まで)
	ItemUnit          string            `json:"itemUnit"`          // 数量単位
	ItemDescription   string            `json:"itemDescription"`   // 数量単位説明
	ThumbnailURL      string            `json:"thumbnailUrl"`      // サムネイルURL
	Media             []*ProductMedia   `json:"media"`             // メディア一覧
	Price             int64             `json:"price"`             // 販売価格(税込)
	ExpirationDate    int64             `json:"expirationDate"`    // 賞味期限(単位:日)
	RecommendedPoint1 string            `json:"recommendedPoint1"` // おすすめポイント1
	RecommendedPoint2 string            `json:"recommendedPoint2"` // おすすめポイント2
	RecommendedPoint3 string            `json:"recommendedPoint3"` // おすすめポイント3
	StorageMethodType StorageMethodType `json:"storageMethodType"` // 保存方法
	DeliveryType      DeliveryType      `json:"deliveryType"`      // 配送方法
	Box60Rate         int64             `json:"box60Rate"`         // 箱の占有率(サイズ:60)
	Box80Rate         int64             `json:"box80Rate"`         // 箱の占有率(サイズ:80)
	Box100Rate        int64             `json:"box100Rate"`        // 箱の占有率(サイズ:100)
	OriginPrefecture  string            `json:"originPrefecture"`  // 原産地(都道府県)
	OriginCity        string            `json:"originCity"`        // 原産地(市区町村)
	Rate              *ProductRate      `json:"rate"`              // 商品評価
	StartAt           int64             `json:"startAt"`           // 販売開始日時
	EndAt             int64             `json:"endAt"`             // 販売終了日時
}

// ProductMedia - 商品メディア情報
type ProductMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

// ProductRate - 商品評価情報
type ProductRate struct {
	Average float64         `json:"average"` // 平均評価
	Count   int64           `json:"count"`   // 合計評価数
	Detail  map[int64]int64 `json:"detail"`  // 評価詳細
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

// MerchantCenterFeedResponse - Google Merchant Center用レスポンス
type MerchantCenterFeedResponse struct {
	XMLName xml.Name               `xml:"rss"`
	Version string                 `xml:"version,attr"`
	Xmlns   string                 `xml:"xmlns:g,attr"`
	Channel *MerchantCenterChannel `xml:"channel"`
}

type MerchantCenterChannel struct {
	Title       string                `xml:"title"`
	Link        string                `xml:"link"`
	Description string                `xml:"description"`
	Items       []*MerchantCenterItem `xml:"item"`
}

// Google Merchant Center 商品データ仕様
// https://support.google.com/merchants/answer/7052112?visit_id=638891829957746964-3945530472&hl=ja
type MerchantCenterItem struct {
	// 商品基本情報
	ID                   string   `xml:"g:id"`                              // 商品ID
	Title                string   `xml:"g:title"`                           // タイトル
	Description          string   `xml:"g:description"`                     // 商品説明
	Link                 string   `xml:"g:link"`                            // 商品リンク
	ImageLink            string   `xml:"g:image_link"`                      // 商品画像リンク
	AdditionalImageLinks []string `xml:"g:additional_image_link,omitempty"` // 追加の商品画像リンク
	// 価格と在庫状況
	Availability           string `xml:"g:availability"`                        // 在庫状況
	AvailabilityDate       string `xml:"g:availability_date,omitempty"`         // 入荷予定日
	CostOfGoodsSold        string `xml:"g:cost_of_goods_sold,omitempty"`        // 商品原価
	ExpirationDate         string `xml:"g:expiration_date,omitempty"`           // 有効期限
	Price                  string `xml:"g:price"`                               // 価格
	SalePrice              string `xml:"g:sale_price,omitempty"`                // セール価格
	SalePriceEffectiveDate string `xml:"g:sale_price_effective_date,omitempty"` // セール価格有効期間
	UnitPricingMeasure     string `xml:"g:unit_pricing_measure,omitempty"`      // 価格の計量単位
	UnitPricingBaseMeasure string `xml:"g:unit_pricing_base_measure,omitempty"` // 価格の基準計量単位
	SubscriptionCost       string `xml:"g:subscription_cost,omitempty"`         // 定期購入の費用
	// 商品カテゴリ
	GoogleProductCategory string `xml:"g:google_product_category,omitempty"` // Google商品カテゴリ
	ProductType           string `xml:"g:product_type,omitempty"`            // 商品カテゴリ
	// 商品ID
	Brand            string `xml:"g:brand"`                       // ブランド
	GTIN             string `xml:"g:gtin,omitempty"`              // GTIN
	MPN              string `xml:"g:mpn,omitempty"`               // MPN
	IdentifierExists string `xml:"g:identifier_exists,omitempty"` // IDの存在
	// 詳細な商品説明
	Condition        string `xml:"g:condition"`                   // 状態
	IsBundle         bool   `xml:"g:is_bundle,omitempty"`         // 一括販売商品
	AgeGroup         string `xml:"g:age_group,omitempty"`         // 年齢層
	Gender           string `xml:"g:gender,omitempty"`            // 性別
	ItemGroupID      string `xml:"g:item_group_id,omitempty"`     // 商品グループID
	ProductLength    string `xml:"g:product_length,omitempty"`    // 商品の奥行
	ProductWidth     string `xml:"g:product_width,omitempty"`     // 商品の幅
	ProductHeight    string `xml:"g:product_height,omitempty"`    // 商品の高さ
	ProductWeight    string `xml:"g:product_weight,omitempty"`    // 商品の重量
	ProductDetail    string `xml:"g:product_detail,omitempty"`    // 商品の詳細
	ProductHighlight string `xml:"g:product_highlight,omitempty"` // 商品に関する情報
	// ショッピング キャンペーンなどの設定
	CustomLabel0 string `xml:"g:custom_label_0,omitempty"` // カスタムラベル0
	CustomLabel1 string `xml:"g:custom_label_1,omitempty"` // カスタムラベル1
	CustomLabel2 string `xml:"g:custom_label_2,omitempty"` // カスタムラベル2
	CustomLabel3 string `xml:"g:custom_label_3,omitempty"` // カスタムラベル3
	CustomLabel4 string `xml:"g:custom_label_4,omitempty"` // カスタムラベル4
	PromotionID  string `xml:"g:promotion_id,omitempty"`   // プロモーションID
	// 掲載先
	ExcludedDestination string `xml:"g:excluded_destination,omitempty"` // 非掲載先
	IncludedDestination string `xml:"g:included_destination,omitempty"` // 掲載先
	// 送料
	Shipping              string `xml:"g:shipping,omitempty"`                // 送料
	ShippingLabel         string `xml:"g:shipping_label,omitempty"`          // 送料ラベル
	ShippingWeight        string `xml:"g:shipping_weight,omitempty"`         // 搬送重量
	ShipsFromCountry      string `xml:"g:ships_from_country,omitempty"`      // 発送元の国
	MaxHandlingTime       string `xml:"g:max_handling_time,omitempty"`       // 最長発送準備期間
	MinHandlingTime       string `xml:"g:min_handling_time,omitempty"`       // 最短発送準備期間
	FreeShippingThreshold string `xml:"g:free_shipping_threshold,omitempty"` // 無料配送の基準額
	// 税金
	Tax         string `xml:"g:tax,omitempty"`          // 税金
	TaxCategory string `xml:"g:tax_category,omitempty"` // 税金カテゴリ
}
