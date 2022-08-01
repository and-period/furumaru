package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	set "github.com/and-period/furumaru/api/pkg/set/v2"
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

func (t *ProductType) Fill(category *Category) {
	if category != nil {
		t.CategoryName = category.Name
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

func (ts ProductTypes) Fill(categories map[string]*Category) {
	for i := range ts {
		category, ok := categories[ts[i].CategoryID]
		if !ok {
			continue
		}
		ts[i].Fill(category)
	}
}

func (ts ProductTypes) Response() []*response.ProductType {
	res := make([]*response.ProductType, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
