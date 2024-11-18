package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
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

func NewSpotUserType(userType entity.SpotUserType) SpotUserType {
	switch userType {
	case entity.SpotUserTypeUser:
		return SpotUserTypeUser
	case entity.SpotUserTypeCoordinator:
		return SpotUserTypeCoordinator
	case entity.SpotUserTypeProducer:
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
	UserType SpotUserType
}

type Spots []*Spot

func NewSpot(spot *entity.Spot) *Spot {
	userType := NewSpotUserType(spot.UserType)
	return &Spot{
		Spot: response.Spot{
			ID:           spot.ID,
			UserType:     userType.Response(),
			UserID:       spot.UserID,
			Name:         spot.Name,
			Description:  spot.Description,
			ThumbnailURL: spot.ThumbnailURL,
			Longitude:    spot.Longitude,
			Latitude:     spot.Latitude,
			Approved:     spot.Approved,
			CreatedAt:    jst.Unix(spot.CreatedAt),
			UpdatedAt:    jst.Unix(spot.UpdatedAt),
		},
		UserType: userType,
	}
}

func (s *Spot) Response() *response.Spot {
	return &s.Spot
}

func NewSpots(spots []*entity.Spot) Spots {
	res := make(Spots, len(spots))
	for i := range spots {
		res[i] = NewSpot(spots[i])
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
