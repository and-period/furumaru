package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ExperienceType struct {
	types.ExperienceType
}

type ExperienceTypes []*ExperienceType

func NewExperienceType(experienceType *entity.ExperienceType) *ExperienceType {
	return &ExperienceType{
		ExperienceType: types.ExperienceType{
			ID:        experienceType.ID,
			Name:      experienceType.Name,
			CreatedAt: experienceType.CreatedAt.Unix(),
			UpdatedAt: experienceType.CreatedAt.Unix(),
		},
	}
}

func (t *ExperienceType) Response() *types.ExperienceType {
	return &t.ExperienceType
}

func NewExperienceTypes(experienceTypes entity.ExperienceTypes) ExperienceTypes {
	res := make(ExperienceTypes, len(experienceTypes))
	for i := range experienceTypes {
		res[i] = NewExperienceType(experienceTypes[i])
	}
	return res
}

func (ts ExperienceTypes) Response() []*types.ExperienceType {
	res := make([]*types.ExperienceType, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
