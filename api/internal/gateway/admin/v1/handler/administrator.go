package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) administratorRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/administrators", h.authentication)

	r.GET("", h.ListAdministrators)
	r.POST("", h.CreateAdministrator)
	r.GET("/:adminId", h.GetAdministrator)
	r.PATCH("/:adminId", h.UpdateAdministrator)
	r.PATCH("/:adminId/email", h.UpdateAdministratorEmail)
	r.PATCH("/:adminId/password", h.ResetAdministratorPassword)
	r.DELETE("/:adminId", h.DeleteAdministrator)
}

func (h *handler) ListAdministrators(ctx *gin.Context) {
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

	in := &user.ListAdministratorsInput{
		Limit:  limit,
		Offset: offset,
	}
	admins, total, err := h.user.ListAdministrators(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AdministratorsResponse{
		Administrators: service.NewAdministrators(admins).Response(),
		Total:          total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetAdministrator(ctx *gin.Context) {
	in := &user.GetAdministratorInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
	}
	admin, err := h.user.GetAdministrator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AdministratorResponse{
		Administrator: service.NewAdministrator(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateAdministrator(ctx *gin.Context) {
	req := &request.CreateAdministratorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.CreateAdministratorInput{
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		Email:         req.Email,
		PhoneNumber:   req.PhoneNumber,
	}
	admin, err := h.user.CreateAdministrator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AdministratorResponse{
		Administrator: service.NewAdministrator(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateAdministrator(ctx *gin.Context) {
	req := &request.UpdateAdministratorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdministratorInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
		Lastname:        req.Lastname,
		Firstname:       req.Firstname,
		LastnameKana:    req.LastnameKana,
		FirstnameKana:   req.FirstnameKana,
		PhoneNumber:     req.PhoneNumber,
	}
	if err := h.user.UpdateAdministrator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateAdministratorEmail(ctx *gin.Context) {
	req := &request.UpdateAdministratorEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdministratorEmailInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
		Email:           req.Email,
	}
	if err := h.user.UpdateAdministratorEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) ResetAdministratorPassword(ctx *gin.Context) {
	in := &user.ResetAdministratorPasswordInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
	}
	if err := h.user.ResetAdministratorPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteAdministrator(ctx *gin.Context) {
	in := &user.DeleteAdministratorInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
	}
	if err := h.user.DeleteAdministrator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) getAdministrator(
	ctx context.Context,
	administratorID string,
) (*service.Administrator, error) {
	in := &user.GetAdministratorInput{
		AdministratorID: administratorID,
	}
	administrator, err := h.user.GetAdministrator(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAdministrator(administrator), nil
}
