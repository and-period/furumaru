package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type SpotType struct {
	response.SpotType
}

type SpotTypes []*SpotType

func NewSpotType(spotType *entity.SpotType) *SpotType {
	return &SpotType{
		SpotType: response.SpotType{
			ID:        spotType.ID,
			Name:      spotType.Name,
			CreatedAt: spotType.CreatedAt.Unix(),
			UpdatedAt: spotType.CreatedAt.Unix(),
		},
	}
}

func (t *SpotType) Response() *response.SpotType {
	return &t.SpotType
}

func NewSpotTypes(spotTypes entity.SpotTypes) SpotTypes {
	res := make(SpotTypes, len(spotTypes))
	for i := range spotTypes {
		res[i] = NewSpotType(spotTypes[i])
	}
	return res
}

func (ts SpotTypes) Response() []*response.SpotType {
	res := make([]*response.SpotType, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
