package service

import (
	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Address struct {
	response.Address
}

type Addresses []*Address

func NewAddress(address *entity.Address) *Address {
	prefecture, _ := codes.ToPrefectureJapanese(address.Prefecture)
	return &Address{
		Address: response.Address{
			ID:           address.ID,
			IsDefault:    address.IsDefault,
			Lastname:     address.Lastname,
			Firstname:    address.Firstname,
			PostalCode:   address.PostalCode,
			Prefecture:   prefecture,
			City:         address.City,
			AddressLine1: address.AddressLine1,
			AddressLine2: address.AddressLine2,
			PhoneNumber:  address.PhoneNumber,
		},
	}
}

func (a *Address) Response() *response.Address {
	return &a.Address
}

func NewAddresses(addresses entity.Addresses) Addresses {
	res := make(Addresses, len(addresses))
	for i := range addresses {
		res[i] = NewAddress(addresses[i])
	}
	return res
}

func (as Addresses) Response() []*response.Address {
	res := make([]*response.Address, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}
