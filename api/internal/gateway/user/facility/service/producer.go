package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Producer struct {
	types.Producer
}

type Producers []*Producer

func NewProducer(producer *entity.Producer, shops entity.Shops) *Producer {
	var coordinatorID string
	if len(shops) > 0 {
		// クライアント側の実装互換のため、はじめのコーディネータのIDを設定
		coordinatorID = shops[0].CoordinatorID
	}
	return &Producer{
		Producer: types.Producer{
			ID:                producer.ID,
			Username:          producer.Username,
			Profile:           producer.Profile,
			ThumbnailURL:      producer.ThumbnailURL,
			HeaderURL:         producer.HeaderURL,
			PromotionVideoURL: producer.PromotionVideoURL,
			InstagramID:       producer.InstagramID,
			FacebookID:        producer.FacebookID,
			Prefecture:        producer.Prefecture,
			City:              producer.City,
			CoordinatorID:     coordinatorID,
		},
	}
}

func (p *Producer) Response() *types.Producer {
	return &p.Producer
}

func NewProducers(producers entity.Producers, shops map[string]entity.Shops) Producers {
	res := make(Producers, len(producers))
	for i := range producers {
		res[i] = NewProducer(producers[i], shops[producers[i].ID])
	}
	return res
}

func (ps Producers) Map() map[string]*Producer {
	res := make(map[string]*Producer, len(ps))
	for i := range ps {
		res[ps[i].ID] = ps[i]
	}
	return res
}

func (ps Producers) Response() []*types.Producer {
	res := make([]*types.Producer, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
