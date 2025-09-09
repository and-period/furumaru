package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type ContactCategory struct {
	types.ContactCategory
}

type ContactCategories []*ContactCategory

func NewContactCategory(contactCategory *entity.ContactCategory) *ContactCategory {
	return &ContactCategory{
		ContactCategory: types.ContactCategory{
			ID:        contactCategory.ID,
			Title:     contactCategory.Title,
			CreatedAt: contactCategory.CreatedAt.Unix(),
			UpdatedAt: contactCategory.UpdatedAt.Unix(),
		},
	}
}

func (c *ContactCategory) Response() *types.ContactCategory {
	return &c.ContactCategory
}

func NewContactCategories(contactCategories entity.ContactCategories) ContactCategories {
	res := make(ContactCategories, len(contactCategories))
	for i := range contactCategories {
		res[i] = NewContactCategory(contactCategories[i])
	}
	return res
}

func (cs ContactCategories) Response() []*types.ContactCategory {
	res := make([]*types.ContactCategory, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
