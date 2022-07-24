package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Producer struct {
	response.Producer
}

type Producers []*Producer

func NewProducer(producer *entity.Producer) *Producer {
	return &Producer{
		Producer: response.Producer{
			ID:            producer.ID,
			Lastname:      producer.Lastname,
			Firstname:     producer.Firstname,
			LastnameKana:  producer.LastnameKana,
			FirstnameKana: producer.FirstnameKana,
			StoreName:     producer.StoreName,
			ThumbnailURL:  producer.ThumbnailURL,
			HeaderURL:     producer.HeaderURL,
			Email:         producer.Email,
			PhoneNumber:   producer.PhoneNumber,
			PostalCode:    producer.PostalCode,
			Prefecture:    producer.Prefecture,
			City:          producer.City,
			AddressLine1:  producer.AddressLine1,
			AddressLine2:  producer.AddressLine2,
			CreatedAt:     producer.CreatedAt.Unix(),
			UpdatedAt:     producer.CreatedAt.Unix(),
		},
	}
}

func (p *Producer) Response() *response.Producer {
	return &p.Producer
}

func NewProducers(producers entity.Producers) Producers {
	res := make(Producers, len(producers))
	for i := range producers {
		res[i] = NewProducer(producers[i])
	}
	return res
}

func (ps Producers) Map() map[string]*Producer {
	res := make(map[string]*Producer, len(ps))
	for _, p := range ps {
		res[p.ID] = p
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
