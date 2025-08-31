package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

type Coordinator struct {
	response.Coordinator
	ShopID string
}

type Coordinators []*Coordinator

func NewCoordinator(coordinator *uentity.Coordinator, shop *sentity.Shop) *Coordinator {
	return &Coordinator{
		Coordinator: response.Coordinator{
			ID:                coordinator.ID,
			Username:          coordinator.Username,
			Profile:           coordinator.Profile,
			ThumbnailURL:      coordinator.ThumbnailURL,
			HeaderURL:         coordinator.HeaderURL,
			PromotionVideoURL: coordinator.PromotionVideoURL,
			InstagramID:       coordinator.InstagramID,
			FacebookID:        coordinator.FacebookID,
			Prefecture:        coordinator.Prefecture,
			City:              coordinator.City,
			MarcheName:        shop.Name,
			ProductTypeIDs:    shop.ProductTypeIDs,
			BusinessDays:      shop.BusinessDays,
		},
		ShopID: shop.ID,
	}
}

func (c *Coordinator) Response() *response.Coordinator {
	return &c.Coordinator
}

func NewCoordinators(coordinators uentity.Coordinators, shops map[string]*sentity.Shop) Coordinators {
	res := make(Coordinators, len(coordinators))
	for i := range coordinators {
		res[i] = NewCoordinator(coordinators[i], shops[coordinators[i].ID])
	}
	return res
}

func (cs Coordinators) Map() map[string]*Coordinator {
	res := make(map[string]*Coordinator, len(cs))
	for i := range cs {
		res[cs[i].ID] = cs[i]
	}
	return res
}

func (cs Coordinators) Response() []*response.Coordinator {
	res := make([]*response.Coordinator, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
