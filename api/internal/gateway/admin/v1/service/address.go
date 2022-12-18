package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Address struct {
	response.Address
	id string
}

type Addresses []*Address

func NewAddress(address *entity.Address) *Address {
	return &Address{
		Address: response.Address{
			Lastname:     address.Lastname,
			Firstname:    address.Firstname,
			PostalCode:   address.PostalCode,
			Prefecture:   address.Prefecture,
			City:         address.City,
			AddressLine1: address.AddressLine1,
			AddressLine2: address.AddressLine2,
			PhoneNumber:  address.PhoneNumber,
		},
		id: address.ID,
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

func (as Addresses) Map() map[string]*Address {
	res := make(map[string]*Address, len(as))
	for _, a := range as {
		res[a.id] = a
	}
	return res
}
