package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Address struct {
	types.Address
	revisionID int64
}

type Addresses []*Address

func NewAddress(address *entity.Address) *Address {
	return &Address{
		Address: types.Address{
			ID:             address.ID,
			IsDefault:      address.IsDefault,
			Lastname:       address.Lastname,
			Firstname:      address.Firstname,
			LastnameKana:   address.LastnameKana,
			FirstnameKana:  address.FirstnameKana,
			PostalCode:     address.PostalCode,
			Prefecture:     address.Prefecture,
			PrefectureCode: address.PrefectureCode,
			City:           address.City,
			AddressLine1:   address.AddressLine1,
			AddressLine2:   address.AddressLine2,
			PhoneNumber:    address.PhoneNumber,
		},
		revisionID: address.AddressRevision.ID,
	}
}

func (a *Address) Response() *types.Address {
	if a == nil {
		return nil
	}
	return &a.Address
}

func NewAddresses(addresses entity.Addresses) Addresses {
	res := make(Addresses, len(addresses))
	for i := range addresses {
		res[i] = NewAddress(addresses[i])
	}
	return res
}

func (as Addresses) MapByRevision() map[int64]*Address {
	res := make(map[int64]*Address, len(as))
	for _, a := range as {
		res[a.revisionID] = a
	}
	return res
}

func (as Addresses) Response() []*types.Address {
	res := make([]*types.Address, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}
