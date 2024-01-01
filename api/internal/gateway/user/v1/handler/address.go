package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) addressRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/addresses", h.authentication)

	r.GET("", h.ListAddresses)
	r.POST("", h.CreateAddress)
	r.GET("/:addressId", h.GetAddress)
	r.PATCH("/:addressId", h.UpdateAddress)
	r.DELETE("/:addressId", h.DeleteAddress)
}

func (h *handler) ListAddresses(ctx *gin.Context) {
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

	in := &user.ListAddressesInput{
		UserID: getUserID(ctx),
		Limit:  limit,
		Offset: offset,
	}
	addresses, total, err := h.user.ListAddresses(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AddressesResponse{
		Addresses: service.NewAddresses(addresses).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetAddress(ctx *gin.Context) {
	in := &user.GetAddressInput{
		AddressID: util.GetParam(ctx, "addressId"),
		UserID:    getUserID(ctx),
	}
	address, err := h.user.GetAddress(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
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
		h.badRequest(ctx, err)
		return
	}
	in := &user.CreateAddressInput{
		UserID:         getUserID(ctx),
		Lastname:       req.Lastname,
		Firstname:      req.Firstname,
		LastnameKana:   req.LastnameKana,
		FirstnameKana:  req.FirstnameKana,
		PostalCode:     req.PostalCode,
		PrefectureCode: req.PrefectureCode,
		City:           req.City,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		PhoneNumber:    req.PhoneNumber,
		IsDefault:      req.IsDefault,
	}
	address, err := h.user.CreateAddress(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
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
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateAddressInput{
		AddressID:      util.GetParam(ctx, "addressId"),
		UserID:         getUserID(ctx),
		Lastname:       req.Lastname,
		Firstname:      req.Firstname,
		LastnameKana:   req.LastnameKana,
		FirstnameKana:  req.FirstnameKana,
		PostalCode:     req.PostalCode,
		PrefectureCode: req.PrefectureCode,
		City:           req.City,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		PhoneNumber:    req.PhoneNumber,
		IsDefault:      req.IsDefault,
	}
	if err := h.user.UpdateAddress(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteAddress(ctx *gin.Context) {
	in := &user.DeleteAddressInput{
		AddressID: util.GetParam(ctx, "addressId"),
		UserID:    getUserID(ctx),
	}
	if err := h.user.DeleteAddress(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

//nolint:unused
func (h *handler) multiGetAddressesByRevision(ctx context.Context, revisionIDs []int64) (service.Addresses, error) {
	if len(revisionIDs) == 0 {
		return service.Addresses{}, nil
	}
	in := &user.MultiGetAddressesByRevisionInput{
		AddressRevisionIDs: revisionIDs,
	}
	addresses, err := h.user.MultiGetAddressesByRevision(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAddresses(addresses), nil
}
