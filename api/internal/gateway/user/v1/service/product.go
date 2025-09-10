package service

import (
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/format"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/shopspring/decimal"
)

// ProductStatus - 商品販売状況
type ProductStatus types.ProductStatus

func NewProductStatus(status entity.ProductStatus) ProductStatus {
	switch status {
	case entity.ProductStatusPrivate:
		return ProductStatus(types.ProductStatusUnknown)
	case entity.ProductStatusPresale:
		return ProductStatus(types.ProductStatusPresale)
	case entity.ProductStatusForSale:
		return ProductStatus(types.ProductStatusForSale)
	case entity.ProductStatusOutOfSale:
		return ProductStatus(types.ProductStatusOutOfSale)
	default:
		return ProductStatus(types.ProductStatusUnknown)
	}
}

func (s ProductStatus) Response() types.ProductStatus {
	return types.ProductStatus(s)
}

// StorageMethodType - 保存方法
type StorageMethodType types.StorageMethodType

func NewStorageMethodType(typ entity.StorageMethodType) StorageMethodType {
	switch typ {
	case entity.StorageMethodTypeNormal:
		return StorageMethodType(types.StorageMethodTypeNormal)
	case entity.StorageMethodTypeCoolDark:
		return StorageMethodType(types.StorageMethodTypeCoolDark)
	case entity.StorageMethodTypeRefrigerated:
		return StorageMethodType(types.StorageMethodTypeRefrigerated)
	case entity.StorageMethodTypeFrozen:
		return StorageMethodType(types.StorageMethodTypeFrozen)
	default:
		return StorageMethodType(types.StorageMethodTypeUnknown)
	}
}

func (t StorageMethodType) Response() types.StorageMethodType {
	return types.StorageMethodType(t)
}

// DeliveryType - 配送方法
type DeliveryType int32

func NewDeliveryType(typ entity.DeliveryType) DeliveryType {
	switch typ {
	case entity.DeliveryTypeNormal:
		return DeliveryType(types.DeliveryTypeNormal)
	case entity.DeliveryTypeFrozen:
		return DeliveryType(types.DeliveryTypeFrozen)
	case entity.DeliveryTypeRefrigerated:
		return DeliveryType(types.DeliveryTypeRefrigerated)
	default:
		return DeliveryType(types.DeliveryTypeUnknown)
	}
}

func (t DeliveryType) Response() types.DeliveryType {
	return types.DeliveryType(t)
}

func NewProductWeight(weight int64, unit entity.WeightUnit) float64 {
	const precision = 1
	var exp int32
	switch unit {
	case entity.WeightUnitGram:
		exp = 0
	case entity.WeightUnitKilogram:
		exp = 3 // 1kg = 1,000g
	default:
		return 0
	}
	gweight := decimal.New(weight, exp)                               // g単位に揃える
	kgweight := gweight.DivRound(decimal.NewFromInt(1000), precision) // 少数第一位までを取得
	fweight, _ := kgweight.Float64()
	return fweight
}

type Product struct {
	types.Product
	revisionID int64
	cost       int64
	status     ProductStatus
	media      MultiProductMedia
}

type Products []*Product

type ProductDetailsParams struct {
	Categories   map[string]*Category
	ProductTypes map[string]*ProductType
	ProductRates map[string]*ProductRate
}

func NewProduct(product *entity.Product, category *Category, rate *ProductRate) *Product {
	var (
		point1, point2, point3 string
		categoryID             string
	)
	if len(product.RecommendedPoints) > 0 {
		point1 = product.RecommendedPoints[0]
	}
	if len(product.RecommendedPoints) > 1 {
		point2 = product.RecommendedPoints[1]
	}
	if len(product.RecommendedPoints) > 2 {
		point3 = product.RecommendedPoints[2]
	}
	if category != nil {
		categoryID = category.ID
	}
	media := NewMultiProductMedia(product.Media)
	return &Product{
		Product: types.Product{
			ID:                product.ID,
			CoordinatorID:     product.CoordinatorID,
			ProducerID:        product.ProducerID,
			CategoryID:        categoryID,
			ProductTypeID:     product.TypeID,
			ProductTagIDs:     product.TagIDs,
			Name:              product.Name,
			Description:       product.Description,
			Status:            NewProductStatus(product.Status).Response(),
			Inventory:         product.Inventory,
			Weight:            NewProductWeight(product.Weight, product.WeightUnit),
			ItemUnit:          product.ItemUnit,
			ItemDescription:   product.ItemDescription,
			ThumbnailURL:      product.ThumbnailURL,
			Media:             media.Response(),
			Price:             product.Price,
			ExpirationDate:    product.ExpirationDate,
			RecommendedPoint1: point1,
			RecommendedPoint2: point2,
			RecommendedPoint3: point3,
			StorageMethodType: NewStorageMethodType(product.StorageMethodType).Response(),
			DeliveryType:      NewDeliveryType(product.DeliveryType).Response(),
			Box60Rate:         product.Box60Rate,
			Box80Rate:         product.Box80Rate,
			Box100Rate:        product.Box100Rate,
			OriginPrefecture:  product.OriginPrefecture,
			OriginCity:        product.OriginCity,
			Rate:              rate.Response(),
			StartAt:           product.StartAt.Unix(),
			EndAt:             product.EndAt.Unix(),
		},
		revisionID: product.ProductRevision.ID,
		status:     NewProductStatus(product.Status),
		cost:       product.ProductRevision.Cost,
		media:      media,
	}
}

func (p *Product) MediaURLs() []string {
	return p.media.URLs()
}

func (p *Product) MerchantCenterItemCondition() string {
	switch types.ProductStatus(p.status) {
	case types.ProductStatusPresale:
		if p.Inventory > 0 {
			return "preorder"
		}
		return "out_of_stock"
	case types.ProductStatusForSale:
		if p.Inventory > 0 {
			return "in_stock"
		}
		return "out_of_stock"
	case types.ProductStatusOutOfSale:
		return "out_of_stock"
	default:
		return "new"
	}
}

func (p *Product) Response() *types.Product {
	return &p.Product
}

func NewProducts(products entity.Products, params *ProductDetailsParams) Products {
	res := make(Products, len(products))
	for i, product := range products {
		var category *Category
		typ, ok := params.ProductTypes[product.TypeID]
		if ok {
			category = params.Categories[typ.CategoryID]
		}
		res[i] = NewProduct(product, category, params.ProductRates[product.ID])
	}
	return res
}

func (ps Products) IDs() []string {
	return set.UniqBy(ps, func(p *Product) string {
		return p.ID
	})
}

func (ps Products) MapByRevision() map[int64]*Product {
	res := make(map[int64]*Product, len(ps))
	for _, p := range ps {
		res[p.revisionID] = p
	}
	return res
}

func (ps Products) Response() []*types.Product {
	res := make([]*types.Product, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}

type ProductMedia struct {
	types.ProductMedia
}

type MultiProductMedia []*ProductMedia

func NewProductMedia(media *entity.ProductMedia) *ProductMedia {
	return &ProductMedia{
		ProductMedia: types.ProductMedia{
			URL:         media.URL,
			IsThumbnail: media.IsThumbnail,
		},
	}
}

func (m *ProductMedia) Response() *types.ProductMedia {
	return &m.ProductMedia
}

func NewMultiProductMedia(media entity.MultiProductMedia) MultiProductMedia {
	res := make(MultiProductMedia, len(media))
	for i := range media {
		res[i] = NewProductMedia(media[i])
	}
	return res
}

func (m MultiProductMedia) URLs() []string {
	res := make([]string, 0, len(m))
	for _, media := range m {
		res = append(res, media.URL)
	}
	return res
}

func (m MultiProductMedia) Response() []*types.ProductMedia {
	res := make([]*types.ProductMedia, len(m))
	for i := range m {
		res[i] = m[i].Response()
	}
	return res
}

type ProductRate struct {
	types.ProductRate
	productID string
}

type ProductRates []*ProductRate

func newProductRate(review *entity.AggregatedProductReview) *ProductRate {
	return &ProductRate{
		ProductRate: types.ProductRate{
			Count:   review.Count,
			Average: format.Round(review.Average, 1),
			Detail: map[int64]int64{
				1: review.Rate1,
				2: review.Rate2,
				3: review.Rate3,
				4: review.Rate4,
				5: review.Rate5,
			},
		},
		productID: review.ProductID,
	}
}

func newEmptyProductRate() *ProductRate {
	return &ProductRate{
		ProductRate: types.ProductRate{
			Count:   0,
			Average: 0.0,
			Detail: map[int64]int64{
				1: 0,
				2: 0,
				3: 0,
				4: 0,
				5: 0,
			},
		},
		productID: "",
	}
}

func (r *ProductRate) Response() *types.ProductRate {
	if r == nil {
		return newEmptyProductRate().Response()
	}
	return &r.ProductRate
}

func NewProductRates(reviews entity.AggregatedProductReviews) ProductRates {
	res := make(ProductRates, len(reviews))
	for i, review := range reviews {
		res[i] = newProductRate(review)
	}
	return res
}

func (rs ProductRates) MapByProductID() map[string]*ProductRate {
	res := make(map[string]*ProductRate, len(rs))
	for _, r := range rs {
		res[r.productID] = r
	}
	return res
}

func (rs ProductRates) Response() []*types.ProductRate {
	res := make([]*types.ProductRate, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}

type MerchantCenterItem struct {
	types.MerchantCenterItem
}

type MerchantCenterItems []*MerchantCenterItem

type NewMerchantCenterItemParams struct {
	Product     *Product
	Coordinator *Coordinator
	ProductType *ProductType
	Category    *Category
	Now         time.Time
	WebURL      func() *url.URL
}

type NewMerchantCenterItemsParams struct {
	Products     Products
	Coordinators map[string]*Coordinator
	Details      *ProductDetailsParams
	Now          time.Time
	WebURL       func() *url.URL
}

func NewMerchantCenterItem(params *NewMerchantCenterItemParams) *MerchantCenterItem {
	const (
		dateFormat   = time.RFC3339 // 日付形式（ISO 8601）
		priceFormat  = "%.0f JPY"   // 金額形式（ISO 4217）
		weightFormat = "%.0f kg"    // 重量形式（キログラム単位）
		pathFormat   = "/items/%s"  // 商品リンクのパス形式
	)

	var (
		description     string
		coordinatorName string
		productTypeName string
		expirationDate  string
	)

	if params.Product.Description == "" {
		description = params.Product.Name
	} else {
		description = params.Product.Description
	}
	if len(description) > 5000 {
		description = description[:4997] + "..."
	}

	if params.Coordinator != nil {
		coordinatorName = params.Coordinator.Username
	}

	switch {
	case params.ProductType != nil && params.Category != nil:
		productTypeName = strings.Join([]string{params.Category.Name, params.ProductType.Name}, " > ")
	case params.Category != nil:
		productTypeName = params.Category.Name
	case params.ProductType != nil:
		productTypeName = params.ProductType.Name
	}

	expiration := jst.ParseFromUnix(params.Product.EndAt).Sub(params.Now)
	if params.Product.EndAt != 0 && expiration > 0 && expiration < 30*24*time.Hour { // 30日以内の有効期限
		expirationDate = jst.ParseFromUnix(params.Product.EndAt).Format(dateFormat)
	}

	link := params.WebURL()
	link.Path = fmt.Sprintf(pathFormat, params.Product.ID)

	return &MerchantCenterItem{
		MerchantCenterItem: types.MerchantCenterItem{
			// 商品基本情報
			ID:                   params.Product.ID,
			Title:                params.Product.Name,
			Description:          html.EscapeString(description),
			Link:                 link.String(),
			ImageLink:            params.Product.ThumbnailURL,
			AdditionalImageLinks: params.Product.MediaURLs(),
			// 価格と在庫状況
			Availability:       params.Product.MerchantCenterItemCondition(),
			AvailabilityDate:   jst.ParseFromUnix(params.Product.StartAt).Format(dateFormat),
			CostOfGoodsSold:    fmt.Sprintf(priceFormat, float64(params.Product.cost)),
			ExpirationDate:     expirationDate,
			Price:              fmt.Sprintf(priceFormat, float64(params.Product.Price)),
			UnitPricingMeasure: "ct", // 個数単位
			// 商品カテゴリ
			ProductType: productTypeName,
			// 商品 ID
			Brand: coordinatorName,
			// 詳細な商品説明
			Condition:        "new", // 新品
			ProductWeight:    fmt.Sprintf(weightFormat, params.Product.Weight),
			ProductHighlight: params.Product.RecommendedPoint1,
			// 送料
			ShippingWeight:   fmt.Sprintf(weightFormat, params.Product.Weight),
			ShipsFromCountry: "JP", // 日本からの発送
		},
	}
}

func (i *MerchantCenterItem) Response() *types.MerchantCenterItem {
	return &i.MerchantCenterItem
}

func NewMerchantCenterItems(params *NewMerchantCenterItemsParams) MerchantCenterItems {
	res := make(MerchantCenterItems, len(params.Products))
	for i, product := range params.Products {
		p := &NewMerchantCenterItemParams{
			Product:     product,
			Coordinator: params.Coordinators[product.CoordinatorID],
			ProductType: params.Details.ProductTypes[product.ProductTypeID],
			Category:    params.Details.Categories[product.CategoryID],
			WebURL:      params.WebURL,
		}
		res[i] = NewMerchantCenterItem(p)
	}
	return res
}

func (is MerchantCenterItems) Response() []*types.MerchantCenterItem {
	res := make([]*types.MerchantCenterItem, len(is))
	for i := range is {
		res[i] = is[i].Response()
	}
	return res
}
