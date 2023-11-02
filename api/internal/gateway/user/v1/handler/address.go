package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *handler) addressRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListAddresses)
	arg.POST("", h.CreateAddress)
	arg.GET("/:addressId", h.GetAddress)
	arg.PATCH("/:addressId", h.UpdateAddress)
	arg.DELETE("/:addressId", h.DeleteAddress)
}

func (h *handler) ListAddresses(ctx *gin.Context) {
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

	in := &store.ListAddressesInput{
		UserID: getUserID(ctx),
		Limit:  limit,
		Offset: offset,
	}
	addresses, total, err := h.store.ListAddresses(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AddressesResponse{
		Addresses: service.NewAddresses(addresses).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetAddress(ctx *gin.Context) {
	in := &store.GetAddressInput{
		AddressID: util.GetParam(ctx, "addressId"),
		UserID:    getUserID(ctx),
	}
	address, err := h.store.GetAddress(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.AddressResponse{
		Address: service.NewAddress(address).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateAddress(ctx *gin.Context) {
	req := &request.CreateAddressRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	in := &store.CreateAddressInput{
		UserID:       getUserID(ctx),
		Lastname:     req.Lastname,
		Firstname:    req.Firstname,
		PostalCode:   req.PostalCode,
		Prefecture:   req.Prefecture,
		City:         req.City,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine1,
		PhoneNumber:  req.PhoneNumber,
		IsDefault:    req.IsDefault,
	}
	address, err := h.store.CreateAddress(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.AddressResponse{
		Address: service.NewAddress(address).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateAddress(ctx *gin.Context) {
	req := &request.UpdateAddressRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	in := &store.UpdateAddressInput{
		AddressID:    util.GetParam(ctx, "addressId"),
		UserID:       getUserID(ctx),
		Lastname:     req.Lastname,
		Firstname:    req.Firstname,
		PostalCode:   req.PostalCode,
		Prefecture:   req.Prefecture,
		City:         req.City,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine1,
		PhoneNumber:  req.PhoneNumber,
		IsDefault:    req.IsDefault,
	}
	if err := h.store.UpdateAddress(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteAddress(ctx *gin.Context) {
	in := &store.DeleteAddressInput{
		AddressID: util.GetParam(ctx, "addressId"),
		UserID:    getUserID(ctx),
	}
	if err := h.store.DeleteAddress(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}
