package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) CreateFacilityUser(ctx context.Context, in *user.CreateFacilityUserInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	if s.now().Before(in.LastCheckInAt) {
		return nil, fmt.Errorf("user: invalid last checkin at: %w", exception.ErrInvalidArgument)
	}
	params := &entity.NewUserParams{
		UserType:      entity.UserTypeFacilityUser,
		Registered:    false, // 施設利用者はゲストと同じ扱いに
		ProducerID:    in.ProducerID,
		ProviderType:  in.ProviderType,
		ExternalID:    in.ProviderID,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		Email:         in.Email,
		PhoneNumber:   in.PhoneNumber,
		LastCheckInAt: in.LastCheckInAt,
	}
	user := entity.NewUser(params)
	if err := s.db.FacilityUser.Create(ctx, user); err != nil {
		return nil, internalError(err)
	}
	return user, nil
}
