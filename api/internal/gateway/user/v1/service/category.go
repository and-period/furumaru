package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Category struct {
	types.Category
}

type Categories []*Category

func NewCategory(category *entity.Category) *Category {
	return &Category{
		Category: types.Category{
			ID:   category.ID,
			Name: category.Name,
		},
	}
}

func (c *Category) Response() *types.Category {
	return &c.Category
}

func NewCategories(categories entity.Categories) Categories {
	res := make(Categories, len(categories))
	for i := range categories {
		res[i] = NewCategory(categories[i])
	}
	return res
}

func (cs Categories) Map() map[string]*Category {
	res := make(map[string]*Category, len(cs))
	for _, c := range cs {
		res[c.ID] = c
	}
	return res
}

func (cs Categories) Response() []*types.Category {
	res := make([]*types.Category, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
