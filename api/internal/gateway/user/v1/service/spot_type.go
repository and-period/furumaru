package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type SpotType struct {
	types.SpotType
}

type SpotTypes []*SpotType

func NewSpotType(spotType *entity.SpotType) *SpotType {
	return &SpotType{
		SpotType: types.SpotType{
			ID:   spotType.ID,
			Name: spotType.Name,
		},
	}
}

func (t *SpotType) Response() *types.SpotType {
	if t == nil {
		return nil
	}
	return &t.SpotType
}

func NewSpotTypes(spotTypes entity.SpotTypes) SpotTypes {
	res := make(SpotTypes, len(spotTypes))
	for i := range spotTypes {
		res[i] = NewSpotType(spotTypes[i])
	}
	return res
}

func (ts SpotTypes) Response() []*types.SpotType {
	res := make([]*types.SpotType, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
