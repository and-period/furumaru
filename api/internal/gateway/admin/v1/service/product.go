package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/shopspring/decimal"
)

// ProductStatus - 商品販売状況
type ProductStatus types.ProductStatus

// ProductScope - 商品公開範囲
type ProductScope types.ProductScope

// StorageMethodType - 保存方法
type StorageMethodType types.StorageMethodType

// DeliveryType - 配送方法
type DeliveryType types.DeliveryType

type Product struct {
	types.Product
	revisionID int64
}

type Products []*Product

type ProductMedia struct {
	types.ProductMedia
}

type MultiProductMedia []*ProductMedia

func NewProductStatus(status entity.ProductStatus) ProductStatus {
	switch status {
	case entity.ProductStatusPrivate:
		return ProductStatus(types.ProductStatusPrivate)
	case entity.ProductStatusPresale:
		return ProductStatus(types.ProductStatusPresale)
	case entity.ProductStatusForSale:
		return ProductStatus(types.ProductStatusForSale)
	case entity.ProductStatusOutOfSale:
		return ProductStatus(types.ProductStatusOutOfSale)
	case entity.ProductStatusArchived:
		return ProductStatus(types.ProductStatusArchived)
	default:
		return ProductStatus(types.ProductStatusUnknown)
	}
}

func (s ProductStatus) Response() types.ProductStatus {
	return types.ProductStatus(s)
}

func NewProductScope(scope entity.ProductScope) ProductScope {
	switch scope {
	case entity.ProductScopePublic:
		return ProductScope(types.ProductScopePublic)
	case entity.ProductScopeLimited:
		return ProductScope(types.ProductScopeLimited)
	case entity.ProductScopePrivate:
		return ProductScope(types.ProductScopePrivate)
	default:
		return ProductScope(types.ProductScopeUnknown)
	}
}

func (s ProductScope) StoreEntity() entity.ProductScope {
	switch types.ProductScope(s) {
	case types.ProductScopePublic:
		return entity.ProductScopePublic
	case types.ProductScopePrivate:
		return entity.ProductScopePrivate
	default:
		return entity.ProductScopeUnknown
	}
}

func (s ProductScope) Response() types.ProductScope {
	return types.ProductScope(s)
}

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

func (t StorageMethodType) StoreEntity() entity.StorageMethodType {
	switch types.StorageMethodType(t) {
	case types.StorageMethodTypeNormal:
		return entity.StorageMethodTypeNormal
	case types.StorageMethodTypeCoolDark:
		return entity.StorageMethodTypeCoolDark
	case types.StorageMethodTypeRefrigerated:
		return entity.StorageMethodTypeRefrigerated
	case types.StorageMethodTypeFrozen:
		return entity.StorageMethodTypeFrozen
	default:
		return entity.StorageMethodTypeUnknown
	}
}

func (t StorageMethodType) Response() types.StorageMethodType {
	return types.StorageMethodType(t)
}

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

func (t DeliveryType) StoreEntity() entity.DeliveryType {
	switch types.DeliveryType(t) {
	case types.DeliveryTypeNormal:
		return entity.DeliveryTypeNormal
	case types.DeliveryTypeFrozen:
		return entity.DeliveryTypeFrozen
	case types.DeliveryTypeRefrigerated:
		return entity.DeliveryTypeRefrigerated
	default:
		return entity.DeliveryTypeUnknown
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
		Product: types.Product{
			ID:                   product.ID,
			CoordinatorID:        product.CoordinatorID,
			ProducerID:           product.ProducerID,
			CategoryID:           "",
			ProductTypeID:        product.TypeID,
			ProductTagIDs:        product.TagIDs,
			Name:                 product.Name,
			Description:          product.Description,
			Scope:                NewProductScope(product.Scope).Response(),
			Status:               NewProductStatus(product.Status).Response(),
			Inventory:            product.Inventory,
			Weight:               NewProductWeight(product.Weight, product.WeightUnit),
			ItemUnit:             product.ItemUnit,
			ItemDescription:      product.ItemDescription,
			Media:                NewMultiProductMedia(product.Media).Response(),
			Price:                product.Price,
			Cost:                 product.Cost,
			ExpirationDate:       product.ExpirationDate,
			RecommendedPoint1:    point1,
			RecommendedPoint2:    point2,
			RecommendedPoint3:    point3,
			StorageMethodType:    NewStorageMethodType(product.StorageMethodType).Response(),
			DeliveryType:         NewDeliveryType(product.DeliveryType).Response(),
			Box60Rate:            product.Box60Rate,
			Box80Rate:            product.Box80Rate,
			Box100Rate:           product.Box100Rate,
			OriginPrefectureCode: product.OriginPrefectureCode,
			OriginCity:           product.OriginCity,
			StartAt:              product.StartAt.Unix(),
			EndAt:                product.EndAt.Unix(),
			CreatedAt:            product.CreatedAt.Unix(),
			UpdatedAt:            product.CreatedAt.Unix(),
		},
		revisionID: product.ProductRevision.ID,
	}
}

func (p *Product) Fill(category *Category) {
	if category != nil {
		p.CategoryID = category.ID
	}
}

func (p *Product) Response() *types.Product {
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
	return set.UniqBy(ps, func(p *Product) string {
		return p.ProducerID
	})
}

func (ps Products) CategoryIDs() []string {
	return set.UniqBy(ps, func(p *Product) string {
		return p.CategoryID
	})
}

func (ps Products) ProductTypeIDs() []string {
	return set.UniqBy(ps, func(p *Product) string {
		return p.ProductTypeID
	})
}

func (ps Products) Map() map[string]*Product {
	res := make(map[string]*Product, len(ps))
	for _, p := range ps {
		res[p.ID] = p
	}
	return res
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

func (ps Products) Response() []*types.Product {
	res := make([]*types.Product, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}

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

func (m MultiProductMedia) Response() []*types.ProductMedia {
	res := make([]*types.ProductMedia, len(m))
	for i := range m {
		res[i] = m[i].Response()
	}
	return res
}
