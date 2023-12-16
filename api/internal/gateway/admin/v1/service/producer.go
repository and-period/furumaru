package service

import (
	"strings"

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
			ID:                producer.ID,
			Status:            NewAdminStatus(producer.Status).Response(),
			CoordinatorID:     producer.CoordinatorID,
			Lastname:          producer.Lastname,
			Firstname:         producer.Firstname,
			LastnameKana:      producer.LastnameKana,
			FirstnameKana:     producer.FirstnameKana,
			Username:          producer.Username,
			Profile:           producer.Profile,
			ThumbnailURL:      producer.ThumbnailURL,
			Thumbnails:        NewImages(producer.Thumbnails).Response(),
			HeaderURL:         producer.HeaderURL,
			Headers:           NewImages(producer.Headers).Response(),
			PromotionVideoURL: producer.PromotionVideoURL,
			BonusVideoURL:     producer.BonusVideoURL,
			InstagramID:       producer.InstagramID,
			FacebookID:        producer.FacebookID,
			Email:             producer.Email,
			PhoneNumber:       producer.PhoneNumber,
			PostalCode:        producer.PostalCode,
			PrefectureCode:    producer.PrefectureCode,
			City:              producer.City,
			AddressLine1:      producer.AddressLine1,
			AddressLine2:      producer.AddressLine2,
			CreatedAt:         producer.CreatedAt.Unix(),
			UpdatedAt:         producer.CreatedAt.Unix(),
		},
	}
}

func (p *Producer) AuthUser() *AuthUser {
	return &AuthUser{
		AuthUser: response.AuthUser{
			AdminID:      p.ID,
			Role:         AdminRoleProducer.Response(),
			Username:     p.Username,
			Email:        p.Email,
			ThumbnailURL: p.ThumbnailURL,
			Thumbnails:   p.Thumbnails,
		},
	}
}

func (p *Producer) Response() *response.Producer {
	return &p.Producer
}

func (p *Producer) Name() string {
	return strings.TrimSpace(strings.Join([]string{p.Lastname, p.Firstname}, " "))
}

func NewProducers(producers entity.Producers) Producers {
	res := make(Producers, len(producers))
	for i := range producers {
		res[i] = NewProducer(producers[i])
	}
	return res
}

func (ps Producers) IDs() []string {
	res := make([]string, len(ps))
	for i := range ps {
		res[i] = ps[i].ID
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

func (ps Producers) Contains(producerIDs ...string) bool {
	pm := ps.Map()
	for _, producerID := range producerIDs {
		if _, ok := pm[producerID]; ok {
			continue
		}
		return false
	}
	return true
}

func (ps Producers) Response() []*response.Producer {
	res := make([]*response.Producer, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
