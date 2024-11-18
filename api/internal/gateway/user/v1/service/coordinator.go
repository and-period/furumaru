package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Coordinator struct {
	response.Coordinator
}

type Coordinators []*Coordinator

func NewCoordinator(coordinator *entity.Coordinator) *Coordinator {
	return &Coordinator{
		Coordinator: response.Coordinator{
			ID:                coordinator.ID,
			MarcheName:        coordinator.MarcheName,
			Username:          coordinator.Username,
			Profile:           coordinator.Profile,
			ProductTypeIDs:    coordinator.ProductTypeIDs,
			BusinessDays:      coordinator.BusinessDays,
			ThumbnailURL:      coordinator.ThumbnailURL,
			HeaderURL:         coordinator.HeaderURL,
			PromotionVideoURL: coordinator.PromotionVideoURL,
			InstagramID:       coordinator.InstagramID,
			FacebookID:        coordinator.FacebookID,
			Prefecture:        coordinator.Prefecture,
			City:              coordinator.City,
		},
	}
}

func (c *Coordinator) Response() *response.Coordinator {
	return &c.Coordinator
}

func NewCoordinators(coordinators entity.Coordinators) Coordinators {
	res := make(Coordinators, len(coordinators))
	for i := range coordinators {
		res[i] = NewCoordinator(coordinators[i])
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
