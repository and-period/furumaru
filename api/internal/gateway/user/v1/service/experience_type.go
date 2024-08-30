package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ExperienceType struct {
	response.ExperienceType
}

type ExperienceTypes []*ExperienceType

func NewExperienceType(experienceType *entity.ExperienceType) *ExperienceType {
	return &ExperienceType{
		ExperienceType: response.ExperienceType{
			ID:   experienceType.ID,
			Name: experienceType.Name,
		},
	}
}

func (t *ExperienceType) Response() *response.ExperienceType {
	return &t.ExperienceType
}

func NewExperienceTypes(experienceTypes entity.ExperienceTypes) ExperienceTypes {
	res := make(ExperienceTypes, len(experienceTypes))
	for i := range experienceTypes {
		res[i] = NewExperienceType(experienceTypes[i])
	}
	return res
}

func (ts ExperienceTypes) Response() []*response.ExperienceType {
	res := make([]*response.ExperienceType, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
