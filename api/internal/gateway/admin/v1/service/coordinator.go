package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Coordinator struct {
	response.Coordinator
}

type Coordinators []*Coordinator

func NewCoordinator(coordinator *entity.Coordinator) *Coordinator {
	return &Coordinator{
		Coordinator: response.Coordinator{
			ID:               coordinator.ID,
			Lastname:         coordinator.Lastname,
			Firstname:        coordinator.Firstname,
			LastnameKana:     coordinator.LastnameKana,
			FirstnameKana:    coordinator.FirstnameKana,
			CompanyName:      coordinator.CompanyName,
			StoreName:        coordinator.StoreName,
			ThumbnailURL:     coordinator.ThumbnailURL,
			Thumbnails:       NewImages(coordinator.Thumbnails).Response(),
			HeaderURL:        coordinator.HeaderURL,
			Headers:          NewImages(coordinator.Headers).Response(),
			TwitterAccount:   coordinator.TwitterAccount,
			InstagramAccount: coordinator.InstagramAccount,
			FacebookAccount:  coordinator.FacebookAccount,
			Email:            coordinator.Email,
			PhoneNumber:      coordinator.PhoneNumber,
			PostalCode:       coordinator.PostalCode,
			Prefecture:       coordinator.Prefecture,
			City:             coordinator.City,
			AddressLine1:     coordinator.AddressLine1,
			AddressLine2:     coordinator.AddressLine2,
			CreatedAt:        coordinator.CreatedAt.Unix(),
			UpdatedAt:        coordinator.CreatedAt.Unix(),
		},
	}
}

func (p *Coordinator) Response() *response.Coordinator {
	return &p.Coordinator
}

func NewCoordinators(coordinators entity.Coordinators) Coordinators {
	res := make(Coordinators, len(coordinators))
	for i := range coordinators {
		res[i] = NewCoordinator(coordinators[i])
	}
	return res
}

func (ps Coordinators) Response() []*response.Coordinator {
	res := make([]*response.Coordinator, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
