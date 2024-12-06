package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// SpotUserType - 投稿者の種別
type SpotUserType int32

const (
	SpotUserTypeUnknown     SpotUserType = 0
	SpotUserTypeUser        SpotUserType = 1 // ユーザー
	SpotUserTypeCoordinator SpotUserType = 2 // コーディネータ
	SpotUserTypeProducer    SpotUserType = 3 // 生産者
)

func NewSpotUserType(t sentity.SpotUserType) SpotUserType {
	switch t {
	case sentity.SpotUserTypeUser:
		return SpotUserTypeUser
	case sentity.SpotUserTypeCoordinator:
		return SpotUserTypeCoordinator
	case sentity.SpotUserTypeProducer:
		return SpotUserTypeProducer
	default:
		return SpotUserTypeUnknown
	}
}

func (t SpotUserType) Response() int32 {
	return int32(t)
}

type Spot struct {
	response.Spot
}

type Spots []*Spot

func NewSpotByUser(spot *sentity.Spot, user *uentity.User) *Spot {
	return &Spot{
		Spot: response.Spot{
			ID:               spot.ID,
			TypeID:           spot.TypeID,
			Name:             spot.Name,
			Description:      spot.Description,
			ThumbnailURL:     spot.ThumbnailURL,
			Longitude:        spot.Longitude,
			Latitude:         spot.Latitude,
			UserType:         SpotUserTypeUser.Response(),
			UserID:           spot.UserID,
			Username:         user.Username,
			UserThumbnailURL: user.ThumbnailURL,
			CreatedAt:        jst.Unix(spot.CreatedAt),
			UpdatedAt:        jst.Unix(spot.UpdatedAt),
		},
	}
}

func NewSpotByCoordinator(spot *sentity.Spot, coordinator *Coordinator) *Spot {
	return &Spot{
		Spot: response.Spot{
			ID:               spot.ID,
			TypeID:           spot.TypeID,
			Name:             spot.Name,
			Description:      spot.Description,
			ThumbnailURL:     spot.ThumbnailURL,
			Longitude:        spot.Longitude,
			Latitude:         spot.Latitude,
			UserType:         SpotUserTypeCoordinator.Response(),
			UserID:           spot.UserID,
			Username:         coordinator.Username,
			UserThumbnailURL: coordinator.ThumbnailURL,
			CreatedAt:        jst.Unix(spot.CreatedAt),
			UpdatedAt:        jst.Unix(spot.UpdatedAt),
		},
	}
}

func NewSpotByProducer(spot *sentity.Spot, producer *Producer) *Spot {
	return &Spot{
		Spot: response.Spot{
			ID:               spot.ID,
			TypeID:           spot.TypeID,
			Name:             spot.Name,
			Description:      spot.Description,
			ThumbnailURL:     spot.ThumbnailURL,
			Longitude:        spot.Longitude,
			Latitude:         spot.Latitude,
			UserType:         SpotUserTypeProducer.Response(),
			UserID:           spot.UserID,
			Username:         producer.Username,
			UserThumbnailURL: producer.ThumbnailURL,
			CreatedAt:        jst.Unix(spot.CreatedAt),
			UpdatedAt:        jst.Unix(spot.UpdatedAt),
		},
	}
}

func (s *Spot) Response() *response.Spot {
	return &s.Spot
}

func NewSpotsByUser(spots []*sentity.Spot, users map[string]*uentity.User) Spots {
	res := make(Spots, 0, len(spots))
	for _, spot := range spots {
		if !spot.Approved {
			continue
		}
		user, ok := users[spot.UserID]
		if !ok {
			continue
		}
		res = append(res, NewSpotByUser(spot, user))
	}
	return res
}

func NewSpotsByCoordinator(spots []*sentity.Spot, coordinators map[string]*Coordinator) Spots {
	res := make(Spots, 0, len(spots))
	for _, spot := range spots {
		if !spot.Approved {
			continue
		}
		coordinator, ok := coordinators[spot.UserID]
		if !ok {
			continue
		}
		res = append(res, NewSpotByCoordinator(spot, coordinator))
	}
	return res
}

func NewSpotsByProducer(spots []*sentity.Spot, producers map[string]*Producer) Spots {
	res := make(Spots, 0, len(spots))
	for _, spot := range spots {
		if !spot.Approved {
			continue
		}
		producer, ok := producers[spot.UserID]
		if !ok {
			continue
		}
		res = append(res, NewSpotByProducer(spot, producer))
	}
	return res
}

func (ss Spots) Response() []*response.Spot {
	res := make([]*response.Spot, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
