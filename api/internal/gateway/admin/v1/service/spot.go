package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
)

// SpotUserType - 投稿者の種別
type SpotUserType int32

const (
	SpotUserTypeUnknown SpotUserType = 0
	SpotUserTypeUser    SpotUserType = 1 // ユーザー
	SpotUserTypeAdmin   SpotUserType = 2 // 管理者
)

func NewSpotUserType(userType entity.SpotUserType) SpotUserType {
	switch userType {
	case entity.SpotUserTypeUser:
		return SpotUserTypeUser
	case entity.SpotUserTypeAdmin:
		return SpotUserTypeAdmin
	default:
		return SpotUserTypeUnknown
	}
}

func NewSpotUserTypeFromInt32(userType int32) SpotUserType {
	return SpotUserType(userType)
}

func (t SpotUserType) Response() int32 {
	return int32(t)
}

type Spot struct {
	response.Spot
	userType SpotUserType
}

type Spots []*Spot

func NewSpot(spot *entity.Spot) *Spot {
	return &Spot{
		Spot: response.Spot{
			ID:           spot.ID,
			UserType:     NewSpotUserType(spot.UserType).Response(),
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
		userType: NewSpotUserType(spot.UserType),
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

func (ss Spots) UserIDs() []string {
	return set.UniqBy(ss, func(s *Spot) string {
		return s.UserID
	})
}

func (ss Spots) GroupByUserType() map[SpotUserType]Spots {
	res := make(map[SpotUserType]Spots, 2)
	for _, s := range ss {
		if _, ok := res[s.userType]; !ok {
			res[s.userType] = make(Spots, 0, len(ss))
		}
		res[s.userType] = append(res[s.userType], s)
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
