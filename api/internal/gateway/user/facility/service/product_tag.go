package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ProductTag struct {
	response.ProductTag
}

type ProductTags []*ProductTag

func NewProductTag(tag *entity.ProductTag) *ProductTag {
	return &ProductTag{
		ProductTag: response.ProductTag{
			ID:   tag.ID,
			Name: tag.Name,
		},
	}
}

func (t *ProductTag) Response() *response.ProductTag {
	return &t.ProductTag
}

func NewProductTags(tags entity.ProductTags) ProductTags {
	res := make(ProductTags, len(tags))
	for i := range tags {
		res[i] = NewProductTag(tags[i])
	}
	return res
}

func (ts ProductTags) Response() []*response.ProductTag {
	res := make([]*response.ProductTag, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
