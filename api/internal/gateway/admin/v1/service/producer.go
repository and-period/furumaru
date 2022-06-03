package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Producer struct {
	*response.Producer
}

type Producers []*Producer

func NewProducer(admin *entity.Admin) *Producer {
	return &Producer{
		Producer: &response.Producer{
			ID:            admin.ID,
			Lastname:      admin.Lastname,
			Firstname:     admin.Firstname,
			LastnameKana:  admin.LastnameKana,
			FirstnameKana: admin.FirstnameKana,
			StoreName:     admin.StoreName,
			ThumbnailURL:  admin.ThumbnailURL,
			Email:         admin.Email,
			PhoneNumber:   admin.PhoneNumber,
			PostalCode:    admin.PostalCode,
			Prefecture:    admin.Prefecture,
			City:          admin.City,
			AddressLine1:  admin.AddressLine1,
			AddressLine2:  admin.AddressLine2,
			CreatedAt:     admin.CreatedAt.Unix(),
			UpdatedAt:     admin.CreatedAt.Unix(),
		},
	}
}

func (p *Producer) Response() *response.Producer {
	return p.Producer
}

func NewProducers(admins entity.Admins) Producers {
	res := make(Producers, len(admins))
	for i := range admins {
		res[i] = NewProducer(admins[i])
	}
	return res
}

func (ps Producers) Response() []*response.Producer {
	res := make([]*response.Producer, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
