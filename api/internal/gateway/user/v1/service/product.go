package service

import (
	"fmt"
	"html"
	"net/url"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/format"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/shopspring/decimal"
)

// ProductStatus - 商品販売状況
type ProductStatus int32

const (
	ProductStatusUnknown   ProductStatus = 0
	ProductStatusPresale   ProductStatus = 1 // 予約受付中
	ProductStatusForSale   ProductStatus = 2 // 販売中
	ProductStatusOutOfSale ProductStatus = 3 // 販売期間外
)

func NewProductStatus(status entity.ProductStatus) ProductStatus {
	switch status {
	case entity.ProductStatusPrivate:
		return ProductStatusUnknown
	case entity.ProductStatusPresale:
		return ProductStatusPresale
	case entity.ProductStatusForSale:
		return ProductStatusForSale
	case entity.ProductStatusOutOfSale:
		return ProductStatusOutOfSale
	default:
		return ProductStatusUnknown
	}
}

func (s ProductStatus) Response() int32 {
	return int32(s)
}

// StorageMethodType - 保存方法
type StorageMethodType int32

const (
	StorageMethodTypeUnknown      StorageMethodType = 0
	StorageMethodTypeNormal       StorageMethodType = 1 // 常温保存
	StorageMethodTypeCoolDark     StorageMethodType = 2 // 冷暗所保存
	StorageMethodTypeRefrigerated StorageMethodType = 3 // 冷蔵保存
	StorageMethodTypeFrozen       StorageMethodType = 4 // 冷凍保存
)

func NewStorageMethodType(typ entity.StorageMethodType) StorageMethodType {
	switch typ {
	case entity.StorageMethodTypeNormal:
		return StorageMethodTypeNormal
	case entity.StorageMethodTypeCoolDark:
		return StorageMethodTypeCoolDark
	case entity.StorageMethodTypeRefrigerated:
		return StorageMethodTypeRefrigerated
	case entity.StorageMethodTypeFrozen:
		return StorageMethodTypeFrozen
	default:
		return StorageMethodTypeUnknown
	}
}

func (t StorageMethodType) Response() int32 {
	return int32(t)
}

// DeliveryType - 配送方法
type DeliveryType int32

const (
	DeliveryTypeUnknown      DeliveryType = 0
	DeliveryTypeNormal       DeliveryType = 1 // 常温便
	DeliveryTypeRefrigerated DeliveryType = 2 // 冷蔵便
	DeliveryTypeFrozen       DeliveryType = 3 // 冷凍便
)

func NewDeliveryType(typ entity.DeliveryType) DeliveryType {
	switch typ {
	case entity.DeliveryTypeNormal:
		return DeliveryTypeNormal
	case entity.DeliveryTypeFrozen:
		return DeliveryTypeFrozen
	case entity.DeliveryTypeRefrigerated:
		return DeliveryTypeRefrigerated
	default:
		return DeliveryTypeUnknown
	}
}

func (t DeliveryType) Response() int32 {
	return int32(t)
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
	response.Product
	revisionID int64
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
	return &Product{
		Product: response.Product{
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
			Media:             NewMultiProductMedia(product.Media).Response(),
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
	}
}

func (p *Product) Response() *response.Product {
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

func (ps Products) Response() []*response.Product {
	res := make([]*response.Product, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}

type ProductMedia struct {
	response.ProductMedia
}

type MultiProductMedia []*ProductMedia

func NewProductMedia(media *entity.ProductMedia) *ProductMedia {
	return &ProductMedia{
		ProductMedia: response.ProductMedia{
			URL:         media.URL,
			IsThumbnail: media.IsThumbnail,
		},
	}
}

func (m *ProductMedia) Response() *response.ProductMedia {
	return &m.ProductMedia
}

func NewMultiProductMedia(media entity.MultiProductMedia) MultiProductMedia {
	res := make(MultiProductMedia, len(media))
	for i := range media {
		res[i] = NewProductMedia(media[i])
	}
	return res
}

func (m MultiProductMedia) Response() []*response.ProductMedia {
	res := make([]*response.ProductMedia, len(m))
	for i := range m {
		res[i] = m[i].Response()
	}
	return res
}

type ProductRate struct {
	response.ProductRate
	productID string
}

type ProductRates []*ProductRate

func newProductRate(review *entity.AggregatedProductReview) *ProductRate {
	return &ProductRate{
		ProductRate: response.ProductRate{
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
		ProductRate: response.ProductRate{
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

func (r *ProductRate) Response() *response.ProductRate {
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

func (rs ProductRates) Response() []*response.ProductRate {
	res := make([]*response.ProductRate, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}

type MerchantCenterItem struct {
	response.MerchantCenterItem
}

type MerchantCenterItems []*MerchantCenterItem

type NewMerchantCenterItemParams struct {
	Product     *Product
	Coordinator *Coordinator
	ProductType *ProductType
	Category    *Category
	WebURL      func() *url.URL
}

type NewMerchantCenterItemsParams struct {
	Products     Products
	Coordinators map[string]*Coordinator
	Details      *ProductDetailsParams
	WebURL       func() *url.URL
}

func NewMerchantCenterItem(params *NewMerchantCenterItemParams) *MerchantCenterItem {
	const condition = "new"

	var (
		description     string
		coordinatorName string
		productTypeName string
		categoryName    string
		availability    string
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

	if params.ProductType != nil {
		productTypeName = params.ProductType.Name
	}

	if params.Category != nil {
		categoryName = params.Category.Name
	}

	if params.Product.Inventory > 0 {
		availability = "in_stock"
	} else {
		availability = "out_of_stock"
	}

	link := params.WebURL()
	link.Path = fmt.Sprintf("/products/%s", params.Product.ID)

	return &MerchantCenterItem{
		MerchantCenterItem: response.MerchantCenterItem{
			ID:                    params.Product.ID,
			Title:                 params.Product.Name,
			Description:           html.EscapeString(description),
			Link:                  link.String(),
			ImageLink:             params.Product.ThumbnailURL,
			Condition:             condition,
			Availability:          availability,
			Price:                 fmt.Sprintf("%.0f JPY", float64(params.Product.Price)),
			Brand:                 coordinatorName,
			GoogleProductCategory: categoryName,
			ProductType:           productTypeName,
			ShippingWeight:        fmt.Sprintf("%.0f g", params.Product.Weight),
		},
	}
}

func (i *MerchantCenterItem) Response() *response.MerchantCenterItem {
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

func (is MerchantCenterItems) Response() []*response.MerchantCenterItem {
	res := make([]*response.MerchantCenterItem, len(is))
	for i := range is {
		res[i] = is[i].Response()
	}
	return res
}
