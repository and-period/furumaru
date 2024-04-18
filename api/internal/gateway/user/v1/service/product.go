package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
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

type Product struct {
	response.Product
	revisionID int64
}

type Products []*Product

type ProductMedia struct {
	response.ProductMedia
}

type MultiProductMedia []*ProductMedia

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

func NewProduct(product *entity.Product) *Product {
	var point1, point2, point3 string
	if len(product.RecommendedPoints) > 0 {
		point1 = product.RecommendedPoints[0]
	}
	if len(product.RecommendedPoints) > 1 {
		point2 = product.RecommendedPoints[1]
	}
	if len(product.RecommendedPoints) > 2 {
		point3 = product.RecommendedPoints[2]
	}
	return &Product{
		Product: response.Product{
			ID:                product.ID,
			CoordinatorID:     product.CoordinatorID,
			ProducerID:        product.ProducerID,
			CategoryID:        "",
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
			StartAt:           product.StartAt.Unix(),
			EndAt:             product.EndAt.Unix(),
		},
		revisionID: product.ProductRevision.ID,
	}
}

func (p *Product) Fill(category *Category) {
	if category != nil {
		p.CategoryID = category.ID
	}
}

func (p *Product) Response() *response.Product {
	return &p.Product
}

func NewProducts(products entity.Products) Products {
	res := make(Products, len(products))
	for i := range products {
		res[i] = NewProduct(products[i])
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

func (ps Products) Fill(types map[string]*ProductType, categories map[string]*Category) {
	for _, p := range ps {
		typ, ok := types[p.ProductTypeID]
		if !ok {
			continue
		}
		p.Fill(categories[typ.CategoryID])
	}
}

func (ps Products) Response() []*response.Product {
	res := make([]*response.Product, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}

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
