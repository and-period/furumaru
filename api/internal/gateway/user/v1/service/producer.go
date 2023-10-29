package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
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
			CoordinatorID:     producer.CoordinatorID,
			Username:          producer.Username,
			Profile:           producer.Profile,
			ThumbnailURL:      producer.ThumbnailURL,
			Thumbnails:        NewImages(producer.Thumbnails).Response(),
			HeaderURL:         producer.HeaderURL,
			Headers:           NewImages(producer.Headers).Response(),
			PromotionVideoURL: producer.PromotionVideoURL,
			InstagramID:       producer.InstagramID,
			FacebookID:        producer.FacebookID,
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

func (ps Producers) Response() []*response.Producer {
	res := make([]*response.Producer, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
