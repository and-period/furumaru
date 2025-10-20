package entity

import "time"

type ShopProducer struct {
	ShopID     string    `gorm:"primaryKey;<-:create"` // 店舗ID
	ProducerID string    `gorm:"primaryKey;<-:create"` // 生産者ID
	CreatedAt  time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time `gorm:""`                     // 更新日時
}

type ShopProducers []*ShopProducer

func (ps ShopProducers) ProducerIDs() []string {
	res := make([]string, len(ps))
	for i, p := range ps {
		res[i] = p.ProducerID
	}
	return res
}

func (ps ShopProducers) GroupByShopID() map[string]ShopProducers {
	res := make(map[string]ShopProducers, len(ps))
	for _, p := range ps {
		if _, ok := res[p.ShopID]; !ok {
			res[p.ShopID] = make(ShopProducers, 0)
		}
		res[p.ShopID] = append(res[p.ShopID], p)
	}
	return res
}
