package handler

import (
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) coordinatorRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListCoordinators)
	arg.POST("", h.CreateCoordinator)
	arg.GET("/:coordinatorId", h.GetCoordinator)
	arg.PATCH("/:coordinatorId", h.UpdateCoordinator)
	arg.PATCH("/:coordinatorId/email", h.UpdateCoordinatorEmail)
	arg.PATCH("/:coordinatorId/password", h.ResetCoordinatorPassword)
}

func (h *handler) ListCoordinators(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	orders, err := h.newCoordinatorOrders(ctx)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.ListCoordinatorsInput{
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	coordinators, total, err := h.user.ListCoordinators(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.CoordinatorsResponse{
		Coordinators: service.NewCoordinators(coordinators).Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newCoordinatorOrders(ctx *gin.Context) ([]*user.ListCoordinatorsOrder, error) {
	params := util.GetOrders(ctx)
	res := make([]*user.ListCoordinatorsOrder, len(params))
	for i := range params {
		var key uentity.CoordinatorOrderBy
		switch params[i].Key {
		case "lastname":
			key = uentity.CoordinatorOrderByLastname
		case "firstname":
			key = uentity.CoordinatorOrderByFirstname
		case "companyName":
			key = uentity.CoordinatorOrderByCompanyName
		case "storeName":
			key = uentity.CoordinatorOrderByStoreName
		case "email":
			key = uentity.CoordinatorOrderByEmail
		case "phoneNumber":
			key = uentity.CoordinatorOrderByPhoneNumber
		default:
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", params[i].Key, errInvalidOrderkey)
		}
		res[i] = &user.ListCoordinatorsOrder{
			Key:        key,
			OrderByASC: params[i].Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) GetCoordinator(ctx *gin.Context) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.CoordinatorResponse{
		Coordinator: service.NewCoordinator(coordinator).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateCoordinator(ctx *gin.Context) {
	req := &request.CreateCoordinatorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.CreateCoordinatorInput{
		Lastname:         req.Lastname,
		Firstname:        req.Firstname,
		LastnameKana:     req.LastnameKana,
		FirstnameKana:    req.FirstnameKana,
		CompanyName:      req.CompanyName,
		StoreName:        req.StoreName,
		ThumbnailURL:     req.ThumbnailURL,
		HeaderURL:        req.HeaderURL,
		TwitterAccount:   req.TwitterAccount,
		InstagramAccount: req.InstagramAccount,
		FacebookAccount:  req.FacebookAccount,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		PostalCode:       req.PostalCode,
		Prefecture:       req.Prefecture,
		City:             req.City,
		AddressLine1:     req.AddressLine1,
		AddressLine2:     req.AddressLine2,
	}
	coordinator, err := h.user.CreateCoordinator(ctx, in)
	if err != nil {
		httpError(ctx, err)
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
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateCoordinatorInput{
		CoordinatorID:    util.GetParam(ctx, "coordinatorId"),
		Lastname:         req.Lastname,
		Firstname:        req.Firstname,
		LastnameKana:     req.LastnameKana,
		FirstnameKana:    req.FirstnameKana,
		CompanyName:      req.CompanyName,
		StoreName:        req.StoreName,
		ThumbnailURL:     req.ThumbnailURL,
		HeaderURL:        req.HeaderURL,
		TwitterAccount:   req.TwitterAccount,
		InstagramAccount: req.InstagramAccount,
		FacebookAccount:  req.FacebookAccount,
		PhoneNumber:      req.PhoneNumber,
		PostalCode:       req.PostalCode,
		Prefecture:       req.Prefecture,
		City:             req.City,
		AddressLine1:     req.AddressLine1,
		AddressLine2:     req.AddressLine2,
	}
	if err := h.user.UpdateCoordinator(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) UpdateCoordinatorEmail(ctx *gin.Context) {
	req := &request.UpdateCoordinatorEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateCoordinatorEmailInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
		Email:         req.Email,
	}
	if err := h.user.UpdateCoordinatorEmail(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ResetCoordinatorPassword(ctx *gin.Context) {
	in := &user.ResetCoordinatorPasswordInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	if err := h.user.ResetCoordinatorPassword(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
