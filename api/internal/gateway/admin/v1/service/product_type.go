package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

type ProductType struct {
	response.ProductType
}

type ProductTypes []*ProductType

func NewProductType(productType *entity.ProductType) *ProductType {
	return &ProductType{
		ProductType: response.ProductType{
			ID:         productType.ID,
			CategoryID: productType.CategoryID,
			Name:       productType.Name,
			CreatedAt:  productType.CreatedAt.Unix(),
			UpdatedAt:  productType.UpdatedAt.Unix(),
		},
	}
}

func (t *ProductType) Response() *response.ProductType {
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
	set := set.New(len(ts))
	for i := range ts {
		set.AddStrings(ts[i].CategoryID)
	}
	return set.Strings()
}

func (ts ProductTypes) Response() []*response.ProductType {
	res := make([]*response.ProductType, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
