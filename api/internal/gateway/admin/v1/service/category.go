package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Category struct {
	*response.Category
}

type Categories []*Category

func NewCategory(category *entity.Category) *Category {
	return &Category{
		Category: &response.Category{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Unix(),
			UpdatedAt: category.UpdatedAt.Unix(),
		},
	}
}

func (c *Category) Response() *response.Category {
	return c.Category
}

func NewCategories(categories entity.Categories) Categories {
	res := make(Categories, len(categories))
	for i := range categories {
		res[i] = NewCategory(categories[i])
	}
	return res
}

func (cs Categories) Response() []*response.Category {
	res := make([]*response.Category, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
