package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/random"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"go.uber.org/zap"
)

func (s *userService) ListCoordinators(
	ctx context.Context, in *user.ListCoordinatorsInput,
) (entity.Coordinators, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListCoordinatorsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	coordinators, err := s.db.Coordinator.List(ctx, params)
	return coordinators, exception.InternalError(err)
}

func (s *userService) GetCoordinator(
	ctx context.Context, in *user.GetCoordinatorInput,
) (*entity.Coordinator, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	coordinator, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	return coordinator, exception.InternalError(err)
}

func (s *userService) CreateCoordinator(
	ctx context.Context, in *user.CreateCoordinatorInput,
) (*entity.Coordinator, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	cognitoID := uuid.Base58Encode(uuid.New())
	password := random.NewStrings(size)
	if err := s.createCognitoAdmin(ctx, cognitoID, in.Email, password); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &entity.NewCoordinatorParams{
		Lastname:         in.Lastname,
		Firstname:        in.Firstname,
		LastnameKana:     in.LastnameKana,
		FirstnameKana:    in.FirstnameKana,
		CompanyName:      in.CompanyName,
		StoreName:        in.StoreName,
		ThumbnailURL:     in.ThumbnailURL,
		HeaderURL:        in.HeaderURL,
		TwitterAccount:   in.TwitterAccount,
		InstagramAccount: in.InstagramAccount,
		FacebookAccount:  in.FacebookAccount,
		Email:            in.Email,
		PhoneNumber:      in.PhoneNumber,
		PostalCode:       in.PostalCode,
		Prefecture:       in.Prefecture,
		City:             in.City,
		AddressLine1:     in.AddressLine1,
		AddressLine2:     in.AddressLine2,
	}
	coordinator := entity.NewCoordinator(params)
	auth := entity.NewAdminAuth(coordinator.ID, cognitoID, entity.AdminRoleCoordinator)
	if err := s.db.Coordinator.Create(ctx, auth, coordinator); err != nil {
		return nil, exception.InternalError(err)
	}
	s.logger.Debug("Create coordinator",
		zap.String("coordinatorId", coordinator.ID), zap.String("password", password))
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyRegisterAdmin(context.Background(), coordinator.Name(), coordinator.Email, password)
		if err != nil {
			s.logger.Warn("Failed to notify register admin", zap.String("coordinatorId", coordinator.ID), zap.Error(err))
		}
	}()
	return coordinator, nil
}
