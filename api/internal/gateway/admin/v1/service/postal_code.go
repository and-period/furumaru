package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type PostalCode struct {
	response.PostalCode
}

func NewPostalCode(c *entity.PostalCode) *PostalCode {
	return &PostalCode{
		PostalCode: response.PostalCode{
			PostalCode:     c.PostalCode,
			PrefectureCode: c.PrefectureCode,
			Prefecture:     c.Prefecture,
			City:           c.City,
			Town:           c.Town,
		},
	}
}

func (c *PostalCode) Response() *response.PostalCode {
	return &c.PostalCode
}
