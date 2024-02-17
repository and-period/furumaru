package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) coordinatorRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/coordinators", h.authentication)

	r.GET("", h.ListCoordinators)
	r.POST("", h.CreateCoordinator)
	r.GET("/:coordinatorId", h.GetCoordinator)
	r.PATCH("/:coordinatorId", h.UpdateCoordinator)
	r.PATCH("/:coordinatorId/email", h.UpdateCoordinatorEmail)
	r.PATCH("/:coordinatorId/password", h.ResetCoordinatorPassword)
	r.DELETE("/:coordinatorId", h.DeleteCoordinator)
}

func (h *handler) ListCoordinators(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.ListCoordinatorsInput{
		Name:   util.GetQuery(ctx, "username", ""),
		Limit:  limit,
		Offset: offset,
	}
	coordinators, total, err := h.user.ListCoordinators(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(coordinators) == 0 {
		res := &response.CoordinatorsResponse{
			Coordinators: []*response.Coordinator{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		producerTotals map[string]int64
		productTypes   service.ProductTypes
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		aggregateIn := &user.AggregateRealatedProducersInput{
			CoordinatorIDs: coordinators.IDs(),
		}
		producerTotals, err = h.user.AggregateRealatedProducers(ctx, aggregateIn)
		return
	})
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, coordinators.ProductTypeIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	scoordinator := service.NewCoordinators(coordinators)
	scoordinator.SetProducerTotal(producerTotals)

	res := &response.CoordinatorsResponse{
		Coordinators: scoordinator.Response(),
		ProductTypes: productTypes.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetCoordinator(ctx *gin.Context) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, coordinator.ProductTypeIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CoordinatorResponse{
		Coordinator:  service.NewCoordinator(coordinator).Response(),
		ProductTypes: productTypes.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateCoordinator(ctx *gin.Context) {
	req := &request.CreateCoordinatorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, req.ProductTypeIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(productTypes) != len(req.ProductTypeIDs) {
		h.badRequest(ctx, errors.New("handler: unmatch product types length"))
		return
	}

	in := &user.CreateCoordinatorInput{
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		MarcheName:        req.MarcheName,
		Username:          req.Username,
		Profile:           req.Profile,
		ProductTypeIDs:    req.ProductTypeIDs,
		ThumbnailURL:      req.ThumbnailURL,
		HeaderURL:         req.HeaderURL,
		PromotionVideoURL: req.PromotionVideoURL,
		BonusVideoURL:     req.BonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		BusinessDays:      req.BusinessDays,
	}
	coordinator, err := h.user.CreateCoordinator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CoordinatorResponse{
		Coordinator: service.NewCoordinator(coordinator).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateCoordinator(ctx *gin.Context) {
	req := &request.UpdateCoordinatorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, req.ProductTypeIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(productTypes) != len(req.ProductTypeIDs) {
		h.badRequest(ctx, errors.New("handler: unmatch product types length"))
		return
	}

	in := &user.UpdateCoordinatorInput{
		CoordinatorID:     util.GetParam(ctx, "coordinatorId"),
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		MarcheName:        req.MarcheName,
		Username:          req.Username,
		Profile:           req.Profile,
		ProductTypeIDs:    req.ProductTypeIDs,
		ThumbnailURL:      req.ThumbnailURL,
		HeaderURL:         req.HeaderURL,
		PromotionVideoURL: req.PromotionVideoURL,
		BonusVideoURL:     req.BonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		BusinessDays:      req.BusinessDays,
	}
	if err := h.user.UpdateCoordinator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateCoordinatorEmail(ctx *gin.Context) {
	req := &request.UpdateCoordinatorEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateCoordinatorEmailInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
		Email:         req.Email,
	}
	if err := h.user.UpdateCoordinatorEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) ResetCoordinatorPassword(ctx *gin.Context) {
	in := &user.ResetCoordinatorPasswordInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	if err := h.user.ResetCoordinatorPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteCoordinator(ctx *gin.Context) {
	in := &user.DeleteCoordinatorInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	if err := h.user.DeleteCoordinator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetCoordinators(ctx context.Context, coordinatorIDs []string) (service.Coordinators, error) {
	if len(coordinatorIDs) == 0 {
		return service.Coordinators{}, nil
	}
	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: coordinatorIDs,
	}
	coordinators, err := h.user.MultiGetCoordinators(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinators(coordinators), nil
}

func (h *handler) multiGetCoordinatorsWithDeleted(ctx context.Context, coordinatorIDs []string) (service.Coordinators, error) {
	if len(coordinatorIDs) == 0 {
		return service.Coordinators{}, nil
	}
	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: coordinatorIDs,
		WithDeleted:    true,
	}
	coordinators, err := h.user.MultiGetCoordinators(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinators(coordinators), nil
}

func (h *handler) getCoordinator(ctx context.Context, coordinatorID string) (*service.Coordinator, error) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: coordinatorID,
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinator(coordinator), nil
}

func (h *handler) getCoordinatorWithDeleted(ctx context.Context, coordinatorID string) (*service.Coordinator, error) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: coordinatorID,
		WithDeleted:   true,
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinator(coordinator), nil
}
