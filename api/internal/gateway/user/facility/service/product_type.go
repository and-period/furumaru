package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

type ProductType struct {
	types.ProductType
}

type ProductTypes []*ProductType

func NewProductType(productType *entity.ProductType) *ProductType {
	return &ProductType{
		ProductType: types.ProductType{
			ID:         productType.ID,
			CategoryID: productType.CategoryID,
			Name:       productType.Name,
			IconURL:    productType.IconURL,
		},
	}
}

func (t *ProductType) Response() *types.ProductType {
	return &t.ProductType
}

func NewProductTypes(productTypes entity.ProductTypes) ProductTypes {
	res := make(ProductTypes, len(productTypes))
	for i := range productTypes {
		res[i] = NewProductType(productTypes[i])
	}
	return res
}

func (ts ProductTypes) CategoryIDs() []string {
	return set.UniqBy(ts, func(t *ProductType) string {
		return t.CategoryID
	})
}

func (ts ProductTypes) Map() map[string]*ProductType {
	res := make(map[string]*ProductType, len(ts))
	for _, t := range ts {
		res[t.ID] = t
	}
	return res
}

func (ts ProductTypes) Response() []*types.ProductType {
	res := make([]*types.ProductType, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
