package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/shopspring/decimal"
)

// DeliveryType - 配送方法
type DeliveryType int32

const (
	DeliveryTypeUnknown      DeliveryType = 0
	DeliveryTypeNormal       DeliveryType = 1 // 通常便
	DeliveryTypeRefrigerated DeliveryType = 2 // 冷蔵便
	DeliveryTypeFrozen       DeliveryType = 3 // 冷凍便
)

type Product struct {
	response.Product
}

type Products []*Product

type ProductMedia struct {
	response.ProductMedia
}

type MultiProductMedia []*ProductMedia

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

func (t DeliveryType) StoreEntity() entity.DeliveryType {
	switch t {
	case DeliveryTypeNormal:
		return entity.DeliveryTypeNormal
	case DeliveryTypeFrozen:
		return entity.DeliveryTypeFrozen
	case DeliveryTypeRefrigerated:
		return entity.DeliveryTypeRefrigerated
	default:
		return entity.DeliveryTypeUnknown
	}
}

func (t DeliveryType) Response() int32 {
	return int32(t)
}

func NewProduct(product *entity.Product) *Product {
	return &Product{
		Product: response.Product{
			ID:               product.ID,
			ProducerID:       product.ProducerID,
			StoreName:        "dummy", // TODO: 詳細の実装
			CategoryID:       product.CategoryID,
			CategoryName:     "dummy", // TODO: 詳細の実装
			TypeID:           product.TypeID,
			TypeName:         "dummy", // TODO: 詳細の実装
			Name:             product.Name,
			Description:      product.Description,
			Public:           product.Public,
			Inventory:        product.Inventory,
			Weight:           NewProductWeight(product.Weight, product.WeightUnit),
			ItemUnit:         product.ItemUnit,
			ItemDescription:  product.ItemDescription,
			Media:            NewMultiProductMedia(product.Media).Response(),
			Price:            product.Price,
			DeliveryType:     NewDeliveryType(product.DeliveryType).Response(),
			Box60Rate:        product.Box60Rate,
			Box80Rate:        product.Box80Rate,
			Box100Rate:       product.Box100Rate,
			OriginPrefecture: product.OriginPrefecture,
			OriginCity:       product.OriginCity,
			CreatedBy:        product.CreatedBy,
			UpdatedBy:        product.UpdatedBy,
			CreatedAt:        product.CreatedAt.Unix(),
			UpdatedAt:        product.CreatedAt.Unix(),
		},
	}
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

func NewProductWeightFromRequest(weight float64) (int64, entity.WeightUnit) {
	dweight := decimal.NewFromFloat(weight).Truncate(1) // 少数第一位までを取得
	if dweight.IsInteger() {
		// kg単位のままで表すことが可能なため(request値はkg前提)
		return dweight.IntPart(), entity.WeightUnitKilogram
	}
	// 少数点が含まれている場合、そのままintに変換できないためg単位に変換
	dweight = dweight.Mul(decimal.NewFromInt(1000))
	return dweight.IntPart(), entity.WeightUnitGram
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

func (ps Products) ProducerIDs() []string {
	set := set.New(len(ps))
	for i := range ps {
		set.AddStrings(ps[i].ProducerID)
	}
	return set.Strings()
}

func (ps Products) CategoryIDs() []string {
	set := set.New(len(ps))
	for i := range ps {
		set.AddStrings(ps[i].CategoryID)
	}
	return set.Strings()
}

func (ps Products) ProductTypeIDs() []string {
	set := set.New(len(ps))
	for i := range ps {
		set.AddStrings(ps[i].TypeID)
	}
	return set.Strings()
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
