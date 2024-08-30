package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) multiGetExperienceTypes(ctx context.Context, experienceTypeIDs []string) (service.ExperienceTypes, error) {
	if len(experienceTypeIDs) == 0 {
		return service.ExperienceTypes{}, nil
	}
	in := &store.MultiGetExperienceTypesInput{
		ExperienceTypeIDs: experienceTypeIDs,
	}
	experienceTypes, err := h.store.MultiGetExperienceTypes(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperienceTypes(experienceTypes), nil
}

func (h *handler) getExperienceType(ctx context.Context, experienceTypeID string) (*service.ExperienceType, error) {
	in := &store.GetExperienceTypeInput{
		ExperienceTypeID: experienceTypeID,
	}
	experienceType, err := h.store.GetExperienceType(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperienceType(experienceType), nil
}
