package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
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
		return nil, 0, internalError(err)
	}
	params := &database.ListCoordinatorsParams{
		Name:   in.Name,
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
		return nil, 0, internalError(err)
	}
	return coordinators, total, nil
}

func (s *service) MultiGetCoordinators(
	ctx context.Context, in *user.MultiGetCoordinatorsInput,
) (entity.Coordinators, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	var (
		coordinators entity.Coordinators
		err          error
	)
	if in.WithDeleted {
		coordinators, err = s.db.Coordinator.MultiGetWithDeleted(ctx, in.CoordinatorIDs)
	} else {
		coordinators, err = s.db.Coordinator.MultiGet(ctx, in.CoordinatorIDs)
	}
	return coordinators, internalError(err)
}

func (s *service) GetCoordinator(
	ctx context.Context, in *user.GetCoordinatorInput,
) (*entity.Coordinator, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	var (
		coordinator *entity.Coordinator
		err         error
	)
	if in.WithDeleted {
		coordinator, err = s.db.Coordinator.GetWithDeleted(ctx, in.CoordinatorID)
	} else {
		coordinator, err = s.db.Coordinator.Get(ctx, in.CoordinatorID)
	}
	return coordinator, internalError(err)
}

func (s *service) CreateCoordinator(
	ctx context.Context, in *user.CreateCoordinatorInput,
) (*entity.Coordinator, string, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, "", internalError(err)
	}
	productTypes, err := s.multiGetProductTypes(ctx, in.ProductTypeIDs)
	if err != nil {
		return nil, "", internalError(err)
	}
	if len(productTypes) != len(in.ProductTypeIDs) {
		return nil, "", fmt.Errorf("api: invalid product type ids: %w", exception.ErrInvalidArgument)
	}
	cognitoID := uuid.Base58Encode(uuid.New())
	password := random.NewStrings(size)
	adminParams := &entity.NewAdminParams{
		CognitoID:     cognitoID,
		Type:          entity.AdminTypeCoordinator,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		Email:         in.Email,
	}
	params := &entity.NewCoordinatorParams{
		Admin:             entity.NewAdmin(adminParams),
		MarcheName:        in.MarcheName,
		Username:          in.Username,
		Profile:           in.Profile,
		ProductTypeIDs:    in.ProductTypeIDs,
		ThumbnailURL:      in.ThumbnailURL,
		HeaderURL:         in.HeaderURL,
		PromotionVideoURL: in.PromotionVideoURL,
		BonusVideoURL:     in.BonusVideoURL,
		InstagramID:       in.InstagramID,
		FacebookID:        in.FacebookID,
		PhoneNumber:       in.PhoneNumber,
		PostalCode:        in.PostalCode,
		PrefectureCode:    in.PrefectureCode,
		City:              in.City,
		AddressLine1:      in.AddressLine1,
		AddressLine2:      in.AddressLine2,
		BusinessDays:      in.BusinessDays,
	}
	coordinator, err := entity.NewCoordinator(params)
	if err != nil {
		return nil, "", fmt.Errorf("service: failed to new coordinator: %w: %s", exception.ErrInvalidArgument, err.Error())
	}
	auth := s.createCognitoAdmin(cognitoID, in.Email, password)
	if err := s.db.Coordinator.Create(ctx, coordinator, auth); err != nil {
		return nil, "", internalError(err)
	}
	s.logger.Debug("Create coordinator", zap.String("coordinatorId", coordinator.ID), zap.String("password", password))
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyRegisterAdmin(context.Background(), coordinator.ID, password)
		if err != nil {
			s.logger.Warn("Failed to notify register admin", zap.String("coordinatorId", coordinator.ID), zap.Error(err))
		}
	}()
	return coordinator, password, nil
}

func (s *service) UpdateCoordinator(ctx context.Context, in *user.UpdateCoordinatorInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if _, err := codes.ToPrefectureJapanese(in.PrefectureCode); err != nil {
		return fmt.Errorf("service: invalid prefecture code: %w: %s", exception.ErrInvalidArgument, err.Error())
	}
	productTypes, err := s.multiGetProductTypes(ctx, in.ProductTypeIDs)
	if err != nil {
		return internalError(err)
	}
	if len(productTypes) != len(in.ProductTypeIDs) {
		return fmt.Errorf("api: invalid product type ids: %w", exception.ErrInvalidArgument)
	}
	params := &database.UpdateCoordinatorParams{
		Lastname:          in.Lastname,
		Firstname:         in.Firstname,
		LastnameKana:      in.LastnameKana,
		FirstnameKana:     in.FirstnameKana,
		MarcheName:        in.MarcheName,
		Username:          in.Username,
		Profile:           in.Profile,
		ProductTypeIDs:    in.ProductTypeIDs,
		ThumbnailURL:      in.ThumbnailURL,
		HeaderURL:         in.HeaderURL,
		PromotionVideoURL: in.PromotionVideoURL,
		BonusVideoURL:     in.BonusVideoURL,
		InstagramID:       in.InstagramID,
		FacebookID:        in.FacebookID,
		PhoneNumber:       in.PhoneNumber,
		PostalCode:        in.PostalCode,
		PrefectureCode:    in.PrefectureCode,
		City:              in.City,
		AddressLine1:      in.AddressLine1,
		AddressLine2:      in.AddressLine2,
		BusinessDays:      in.BusinessDays,
	}
	if err := s.db.Coordinator.Update(ctx, in.CoordinatorID, params); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) UpdateCoordinatorEmail(ctx context.Context, in *user.UpdateCoordinatorEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	coordinator, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if err != nil {
		return internalError(err)
	}
	params := &cognito.AdminChangeEmailParams{
		Username: coordinator.CognitoID,
		Email:    in.Email,
	}
	if err := s.adminAuth.AdminChangeEmail(ctx, params); err != nil {
		return internalError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, in.CoordinatorID, in.Email)
	return internalError(err)
}

func (s *service) ResetCoordinatorPassword(ctx context.Context, in *user.ResetCoordinatorPasswordInput) error {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	coordinator, err := s.db.Coordinator.Get(ctx, in.CoordinatorID)
	if err != nil {
		return internalError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminChangePasswordParams{
		Username:  coordinator.CognitoID,
		Password:  password,
		Permanent: true,
	}
	if err := s.adminAuth.AdminChangePassword(ctx, params); err != nil {
		return internalError(err)
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

func (s *service) RemoveCoordinatorProductType(ctx context.Context, in *user.RemoveCoordinatorProductTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Coordinator.RemoveProductTypeID(ctx, in.ProductTypeID)
	return internalError(err)
}

func (s *service) DeleteCoordinator(ctx context.Context, in *user.DeleteCoordinatorInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Coordinator.Delete(ctx, in.CoordinatorID, s.deleteCognitoAdmin(in.CoordinatorID))
	return internalError(err)
}

func (s *service) AggregateRealatedProducers(
	ctx context.Context, in *user.AggregateRealatedProducersInput,
) (map[string]int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	res, err := s.db.Producer.AggregateByCoordinatorID(ctx, in.CoordinatorIDs)
	return res, internalError(err)
}

func (s *service) multiGetProductTypes(ctx context.Context, productTypeIDs []string) (sentity.ProductTypes, error) {
	if len(productTypeIDs) == 0 {
		return sentity.ProductTypes{}, nil
	}
	in := &store.MultiGetProductTypesInput{
		ProductTypeIDs: productTypeIDs,
	}
	return s.store.MultiGetProductTypes(ctx, in)
}
