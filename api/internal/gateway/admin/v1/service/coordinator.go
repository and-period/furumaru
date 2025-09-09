package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Coordinator struct {
	types.Coordinator
}

type Coordinators []*Coordinator

func NewCoordinator(coordinator *entity.Coordinator, shop *Shop) *Coordinator {
	return &Coordinator{
		Coordinator: types.Coordinator{
			ID:                coordinator.ID,
			ShopID:            shop.GetID(),
			Status:            NewAdminStatus(coordinator.Status).Response(),
			Lastname:          coordinator.Lastname,
			Firstname:         coordinator.Firstname,
			LastnameKana:      coordinator.LastnameKana,
			FirstnameKana:     coordinator.FirstnameKana,
			Username:          coordinator.Username,
			Profile:           coordinator.Profile,
			ThumbnailURL:      coordinator.ThumbnailURL,
			HeaderURL:         coordinator.HeaderURL,
			PromotionVideoURL: coordinator.PromotionVideoURL,
			BonusVideoURL:     coordinator.BonusVideoURL,
			InstagramID:       coordinator.InstagramID,
			FacebookID:        coordinator.FacebookID,
			Email:             coordinator.Email,
			PhoneNumber:       coordinator.PhoneNumber,
			PostalCode:        coordinator.PostalCode,
			PrefectureCode:    coordinator.PrefectureCode,
			City:              coordinator.City,
			AddressLine1:      coordinator.AddressLine1,
			AddressLine2:      coordinator.AddressLine2,
			ProducerTotal:     0,
			CreatedAt:         coordinator.CreatedAt.Unix(),
			UpdatedAt:         coordinator.CreatedAt.Unix(),
		},
	}
}

func (c *Coordinator) AuthUser() *AuthUser {
	return &AuthUser{
		AuthUser: types.AuthUser{
			AdminID:      c.ID,
			ShopIDs:      []string{c.ShopID},
			Type:         AdminTypeCoordinator.Response(),
			Username:     c.Username,
			Email:        c.Email,
			ThumbnailURL: c.ThumbnailURL,
		},
	}
}

func (c *Coordinator) Response() *types.Coordinator {
	return &c.Coordinator
}

func NewCoordinators(coordinators entity.Coordinators, shops map[string]*Shop) Coordinators {
	res := make(Coordinators, len(coordinators))
	for i := range coordinators {
		res[i] = NewCoordinator(coordinators[i], shops[coordinators[i].AdminID])
	}
	return res
}

func (cs Coordinators) SetProducerTotal(totalMap map[string]int64) {
	for _, c := range cs {
		total, ok := totalMap[c.ID]
		if !ok {
			continue
		}
		c.ProducerTotal = total
	}
}

func (cs Coordinators) Map() map[string]*Coordinator {
	res := make(map[string]*Coordinator, len(cs))
	for _, c := range cs {
		res[c.ID] = c
	}
	return res
}

func (cs Coordinators) Response() []*types.Coordinator {
	res := make([]*types.Coordinator, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
