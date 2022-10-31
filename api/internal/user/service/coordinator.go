package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/random"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListCoordinators(
	ctx context.Context, in *user.ListCoordinatorsInput,
) (entity.Coordinators, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListCoordinatorsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		coordinators entity.Coordinators
		total        int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinators, err = s.db.Coordinator.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Coordinator.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return coordinators, total, nil
}

func (s *service) MultiGetCoordinators(
	ctx context.Context, in *user.MultiGetCoordinatorsInput,
) (entity.Coordinators, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	coordinators, err := s.db.Coordinator.MultiGet(ctx, in.CoordinatorIDs)
	return coordinators, exception.InternalError(err)
}

func (s *service) GetCoordinator(
	ctx context.Context, in *user.GetCoordinatorInput,
) (*entity.Coordinator, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	coordinator, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	return coordinator, exception.InternalError(err)
}

func (s *service) CreateCoordinator(
	ctx context.Context, in *user.CreateCoordinatorInput,
) (*entity.Coordinator, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	cognitoID := uuid.Base58Encode(uuid.New())
	password := random.NewStrings(size)
	adminParams := &entity.NewAdminParams{
		CognitoID:     cognitoID,
		Role:          entity.AdminRoleCoordinator,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		Email:         in.Email,
	}
	params := &entity.NewCoordinatorParams{
		Admin:            entity.NewAdmin(adminParams),
		CompanyName:      in.CompanyName,
		StoreName:        in.StoreName,
		ThumbnailURL:     in.ThumbnailURL,
		HeaderURL:        in.HeaderURL,
		TwitterAccount:   in.TwitterAccount,
		InstagramAccount: in.InstagramAccount,
		FacebookAccount:  in.FacebookAccount,
		PhoneNumber:      in.PhoneNumber,
		PostalCode:       in.PostalCode,
		Prefecture:       in.Prefecture,
		City:             in.City,
		AddressLine1:     in.AddressLine1,
		AddressLine2:     in.AddressLine2,
	}
	coordinator := entity.NewCoordinator(params)
	auth := s.createCognitoAdmin(cognitoID, in.Email, password)
	if err := s.db.Coordinator.Create(ctx, coordinator, auth); err != nil {
		return nil, exception.InternalError(err)
	}
	s.logger.Debug("Create coordinator", zap.String("coordinatorId", coordinator.ID), zap.String("password", password))
	s.waitGroup.Add(2)
	go func() {
		defer s.waitGroup.Done()
		s.resizeCoordinator(context.Background(), coordinator.ID, in.ThumbnailURL, in.HeaderURL)
	}()
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyRegisterAdmin(context.Background(), coordinator.ID, password)
		if err != nil {
			s.logger.Warn("Failed to notify register admin", zap.String("coordinatorId", coordinator.ID), zap.Error(err))
		}
	}()
	return coordinator, nil
}

func (s *service) UpdateCoordinator(ctx context.Context, in *user.UpdateCoordinatorInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	coordinator, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateCoordinatorParams{
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
		PhoneNumber:      in.PhoneNumber,
		PostalCode:       in.PostalCode,
		Prefecture:       in.Prefecture,
		City:             in.City,
		AddressLine1:     in.AddressLine1,
		AddressLine2:     in.AddressLine2,
	}
	if err := s.db.Coordinator.Update(ctx, in.CoordinatorID, params); err != nil {
		return exception.InternalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		var thumbnailURL, headerURL string
		if coordinator.ThumbnailURL != in.ThumbnailURL {
			thumbnailURL = in.ThumbnailURL
		}
		if coordinator.HeaderURL != in.HeaderURL {
			headerURL = in.HeaderURL
		}
		s.resizeCoordinator(context.Background(), coordinator.ID, thumbnailURL, headerURL)
	}()
	return nil
}

func (s *service) UpdateCoordinatorEmail(ctx context.Context, in *user.UpdateCoordinatorEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	coordinator, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.AdminChangeEmailParams{
		Username: coordinator.CognitoID,
		Email:    in.Email,
	}
	if err := s.adminAuth.AdminChangeEmail(ctx, params); err != nil {
		return exception.InternalError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, in.CoordinatorID, in.Email)
	return exception.InternalError(err)
}

func (s *service) ResetCoordinatorPassword(ctx context.Context, in *user.ResetCoordinatorPasswordInput) error {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	coordinator, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if err != nil {
		return exception.InternalError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminChangePasswordParams{
		Username:  coordinator.CognitoID,
		Password:  password,
		Permanent: true,
	}
	if err := s.adminAuth.AdminChangePassword(ctx, params); err != nil {
		return exception.InternalError(err)
	}
	s.logger.Debug("Reset coordinator password",
		zap.String("coordinatorId", in.CoordinatorID), zap.String("password", password),
	)
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyResetAdminPassword(context.Background(), in.CoordinatorID, password)
		if err != nil {
			s.logger.Warn("Failed to notify reset admin password",
				zap.String("coordinatorId", in.CoordinatorID), zap.Error(err),
			)
		}
	}()
	return nil
}

func (s *service) DeleteCoordinator(ctx context.Context, in *user.DeleteCoordinatorInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Coordinator.Delete(ctx, in.CoordinatorID, s.deleteCognitoAdmin(in.CoordinatorID))
	return exception.InternalError(err)
}

func (s *service) resizeCoordinator(ctx context.Context, coordinatorID, thumbnailURL, headerURL string) {
	s.waitGroup.Add(2)
	go func() {
		defer s.waitGroup.Done()
		if thumbnailURL == "" {
			return
		}
		in := &media.ResizeFileInput{
			TargetID: coordinatorID,
			URLs:     []string{thumbnailURL},
		}
		if err := s.media.ResizeCoordinatorThumbnail(ctx, in); err != nil {
			s.logger.Error("Failed to resize coordinator thumbnail",
				zap.String("coordinatorId", coordinatorID), zap.Error(err),
			)
		}
	}()
	go func() {
		defer s.waitGroup.Done()
		if headerURL == "" {
			return
		}
		in := &media.ResizeFileInput{
			TargetID: coordinatorID,
			URLs:     []string{headerURL},
		}
		if err := s.media.ResizeCoordinatorHeader(ctx, in); err != nil {
			s.logger.Error("Failed to resize coordinator header",
				zap.String("coordinatorId", coordinatorID), zap.Error(err),
			)
		}
	}()
}
