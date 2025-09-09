package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ProductTag struct {
	types.ProductTag
}

type ProductTags []*ProductTag

func NewProductTag(tag *entity.ProductTag) *ProductTag {
	return &ProductTag{
		ProductTag: types.ProductTag{
			ID:   tag.ID,
			Name: tag.Name,
		},
	}
}

func (t *ProductTag) Response() *types.ProductTag {
	return &t.ProductTag
}

func NewProductTags(tags entity.ProductTags) ProductTags {
	res := make(ProductTags, len(tags))
	for i := range tags {
		res[i] = NewProductTag(tags[i])
	}
	return res
}

func (ts ProductTags) Response() []*types.ProductTag {
	res := make([]*types.ProductTag, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
