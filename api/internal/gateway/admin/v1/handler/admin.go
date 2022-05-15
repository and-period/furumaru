package handler

import (
	"net/http"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/request"
	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/and-period/marche/api/internal/gateway/admin/v1/service"
	"github.com/and-period/marche/api/internal/gateway/util"
	"github.com/and-period/marche/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) adminRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.authentication(), h.ListAdmins)
	rg.POST("", h.authentication(), h.CreateAdmin)
	rg.GET("/:adminId", h.authentication(), h.GetAdmin)
	rg.GET("/me", h.authentication(), h.GetAdminMe)
	rg.PATCH("/me/email", h.authentication(), h.UpdateAdminEmail)
	rg.POST("/me/email/verified", h.VerifyAdminEmail)
	rg.PATCH("/me/password", h.authentication(), h.UpdateAdminPassword)
}

func (h *apiV1Handler) ListAdmins(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

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
	roles, err := util.GetQueryInt32s(ctx, "roles")
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.ListAdminsInput{
		Roles:  roles,
		Limit:  limit,
		Offset: offset,
	}
	admins, err := h.user.ListAdmins(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AdminsResponse{
		Admins: service.NewAdmins(admins).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) GetAdmin(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	in := &user.GetAdminInput{
		AdminID: util.GetParam(ctx, "adminId"),
	}
	admin, err := h.user.GetAdmin(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AdminResponse{
		Admin: service.NewAdmin(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) CreateAdmin(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.CreateAdminRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.CreateAdminInput{
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		Email:         req.Email,
		Role:          req.Role,
	}
	admin, err := h.user.CreateAdmin(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AdminResponse{
		Admin: service.NewAdmin(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) GetAdminMe(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	in := &user.GetAdminInput{
		AdminID: getAdminID(ctx),
	}
	admin, err := h.user.GetAdmin(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AdminMeResponse{
		ID:            admin.ID,
		Lastname:      admin.Lastname,
		Firstname:     admin.Firstname,
		LastnameKana:  admin.LastnameKana,
		FirstnameKana: admin.FirstnameKana,
		Email:         admin.Email,
		Role:          service.NewAdminRole(admin.Role).Response(),
		ThumbnailURL:  admin.ThumbnailURL,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) UpdateAdminEmail(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAdminEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdminEmailInput{
		AccessToken: token,
		Email:       req.Email,
	}
	if err := h.user.UpdateAdminEmail(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) VerifyAdminEmail(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.VerifyAdminEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.VerifyAdminEmailInput{
		AccessToken: token,
		VerifyCode:  req.VerifyCode,
	}
	if err := h.user.VerifyAdminEmail(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) UpdateAdminPassword(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAdminPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdminPasswordInput{
		AccessToken:          token,
		OldPassword:          req.OldPassword,
		NewPassword:          req.NewPassword,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.UpdateAdminPassword(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
