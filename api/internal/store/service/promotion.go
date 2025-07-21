package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListPromotions(
	ctx context.Context, in *store.ListPromotionsInput,
) (entity.Promotions, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders, err := s.newListPromotionsOrders(in.Orders)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"service: invalid list promotions orders: err=%s: %w",
			err.Error(),
			exception.ErrInvalidArgument,
		)
	}
	params := &database.ListPromotionsParams{
		ShopID:        in.ShopID,
		Title:         in.Title,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
		Orders:        orders,
		WithAllTarget: in.WithAllTarget,
	}
	var (
		promotions entity.Promotions
		total      int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		promotions, err = s.db.Promotion.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Promotion.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return promotions, total, nil
}

func (s *service) newListPromotionsOrders(
	in []*store.ListPromotionsOrder,
) ([]*database.ListPromotionsOrder, error) {
	res := make([]*database.ListPromotionsOrder, len(in))
	for i := range in {
		var key database.ListPromotionsOrderKey
		switch in[i].Key {
		case store.ListPromotionsOrderByTitle:
			key = database.ListPromotionsOrderByTitle
		case store.ListPromotionsOrderByPublic:
			key = database.ListPromotionsOrderByPublic
		case store.ListPromotionsOrderByStartAt:
			key = database.ListPromotionsOrderByStartAt
		case store.ListPromotionsOrderByEndAt:
			key = database.ListPromotionsOrderByEndAt
		case store.ListPromotionsOrderByCreatedAt:
			key = database.ListPromotionsOrderByCreatedAt
		case store.ListPromotionsOrderByUpdatedAt:
			key = database.ListPromotionsOrderByUpdatedAt
		default:
			return nil, errors.New("service: invalid order key")
		}
		res[i] = &database.ListPromotionsOrder{
			Key:        key,
			OrderByASC: in[i].OrderByASC,
		}
	}
	return res, nil
}

func (s *service) MultiGetPromotions(
	ctx context.Context, in *store.MultiGetPromotionsInput,
) (entity.Promotions, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	promotions, err := s.db.Promotion.MultiGet(ctx, in.PromotionIDs)
	return promotions, internalError(err)
}

func (s *service) GetPromotion(
	ctx context.Context,
	in *store.GetPromotionInput,
) (*entity.Promotion, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	promotion, err := s.db.Promotion.Get(ctx, in.PromotionID)
	if err != nil {
		return nil, internalError(err)
	}
	if in.OnlyEnabled && !promotion.IsEnabled(in.ShopID) {
		return nil, fmt.Errorf("this promotion is disabled: %w", exception.ErrNotFound)
	}
	return promotion, nil
}

func (s *service) GetPromotionByCode(
	ctx context.Context,
	in *store.GetPromotionByCodeInput,
) (*entity.Promotion, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	promotion, err := s.db.Promotion.GetByCode(ctx, in.PromotionCode)
	if err != nil {
		return nil, internalError(err)
	}
	if in.OnlyEnabled && !promotion.IsEnabled(in.ShopID) {
		return nil, fmt.Errorf("this promotion is disabled: %w", exception.ErrNotFound)
	}
	return promotion, nil
}

func (s *service) CreatePromotion(
	ctx context.Context,
	in *store.CreatePromotionInput,
) (*entity.Promotion, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}

	adminIn := &user.GetAdminInput{
		AdminID: in.AdminID,
	}
	admin, err := s.user.GetAdmin(ctx, adminIn)
	if err != nil {
		return nil, internalError(err)
	}
	var shopID string
	switch admin.Type {
	case uentity.AdminTypeAdministrator:
	case uentity.AdminTypeCoordinator:
		shop, err := s.db.Shop.GetByCoordinatorID(ctx, admin.ID)
		if err != nil {
			return nil, internalError(err)
		}
		shopID = shop.ID
	default:
		return nil, fmt.Errorf("service: invalid admin type: %w", exception.ErrForbidden)
	}

	params := &entity.NewPromotionParams{
		ShopID:       shopID,
		Title:        in.Title,
		Description:  in.Description,
		Public:       in.Public,
		DiscountType: in.DiscountType,
		DiscountRate: in.DiscountRate,
		Code:         in.Code,
		CodeType:     in.CodeType,
		StartAt:      in.StartAt,
		EndAt:        in.EndAt,
	}
	promotion := entity.NewPromotion(params)
	if err := promotion.Validate(); err != nil {
		return nil, fmt.Errorf(
			"api: validation error: %s: %w",
			err.Error(),
			exception.ErrInvalidArgument,
		)
	}
	if err := s.db.Promotion.Create(ctx, promotion); err != nil {
		return nil, internalError(err)
	}
	return promotion, nil
}

func (s *service) UpdatePromotion(ctx context.Context, in *store.UpdatePromotionInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	adminIn := &user.GetAdminInput{
		AdminID: in.AdminID,
	}
	admin, err := s.user.GetAdmin(ctx, adminIn)
	if err != nil {
		return internalError(err)
	}

	promotion, err := s.db.Promotion.Get(ctx, in.PromotionID)
	if err != nil {
		return internalError(err)
	}
	switch promotion.TargetType {
	case entity.PromotionTargetTypeAllShop:
		if admin.Type != uentity.AdminTypeAdministrator {
			return fmt.Errorf(
				"service: cannot update promotion for all shops: %w",
				exception.ErrForbidden,
			)
		}
	case entity.PromotionTargetTypeSpecificShop:
		switch admin.Type {
		case uentity.AdminTypeAdministrator:
		case uentity.AdminTypeCoordinator:
			shop, err := s.db.Shop.GetByCoordinatorID(ctx, admin.ID)
			if err != nil {
				return internalError(err)
			}
			if promotion.ShopID != shop.ID {
				return fmt.Errorf(
					"service: this coordinator does not have permission to update promotion: %w",
					exception.ErrForbidden,
				)
			}
		default:
			return fmt.Errorf(
				"service: cannot update promotion for only shop: %w",
				exception.ErrForbidden,
			)
		}
	}

	params := &database.UpdatePromotionParams{
		Title:        in.Title,
		Description:  in.Description,
		Public:       in.Public,
		DiscountType: in.DiscountType,
		DiscountRate: in.DiscountRate,
		Code:         in.Code,
		CodeType:     in.CodeType,
		StartAt:      in.StartAt,
		EndAt:        in.EndAt,
	}
	err = s.db.Promotion.Update(ctx, in.PromotionID, params)
	return internalError(err)
}

func (s *service) DeletePromotion(ctx context.Context, in *store.DeletePromotionInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Promotion.Delete(ctx, in.PromotionID)
	return internalError(err)
}
