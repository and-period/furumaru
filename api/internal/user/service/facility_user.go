package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) GetFacilityUser(ctx context.Context, in *user.GetFacilityUserInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	fuser, err := s.db.FacilityUser.GetByExternalID(ctx, in.ProviderType, in.ProviderID, in.ProducerID)
	if err != nil {
		return nil, internalError(err)
	}
	user, err := s.db.User.Get(ctx, fuser.UserID)
	return user, internalError(err)
}

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

func (s *service) UpdateFacilityUser(ctx context.Context, in *user.UpdateFacilityUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if s.now().Before(in.LastCheckInAt) {
		return fmt.Errorf("user: invalid last checkin at: %w", exception.ErrInvalidArgument)
	}
	params := &database.UpdateFacilityUserParams{
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		PhoneNumber:   in.PhoneNumber,
		LastCheckInAt: in.LastCheckInAt,
	}
	err := s.db.FacilityUser.Update(ctx, in.UserID, params)
	return internalError(err)
}
