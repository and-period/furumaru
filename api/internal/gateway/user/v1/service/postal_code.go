package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type PostalCode struct {
	types.PostalCode
}

func NewPostalCode(c *entity.PostalCode) *PostalCode {
	return &PostalCode{
		PostalCode: types.PostalCode{
			PostalCode:     c.PostalCode,
			PrefectureCode: c.PrefectureCode,
			Prefecture:     c.Prefecture,
			City:           c.City,
			Town:           c.Town,
		},
	}
}

func (c *PostalCode) Response() *types.PostalCode {
	return &c.PostalCode
}
