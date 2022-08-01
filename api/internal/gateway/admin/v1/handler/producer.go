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

func (h *handler) producerRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListProducers)
	arg.POST("", h.CreateProducer)
	arg.GET("/:producerId", h.GetProducer)
	arg.PATCH("/:producerId", h.UpdateProducer)
	arg.PATCH("/:producerId/email", h.UpdateProducerEmail)
	arg.PATCH("/:producerId/password", h.ResetProducerPassword)
}

func (h *handler) ListProducers(ctx *gin.Context) {
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
	orders, err := h.newProducerOrders(ctx)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.ListProducersInput{
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	producers, total, err := h.user.ListProducers(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProducersResponse{
		Producers: service.NewProducers(producers).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newProducerOrders(ctx *gin.Context) ([]*user.ListProducersOrder, error) {
	params := util.GetOrders(ctx)
	res := make([]*user.ListProducersOrder, len(params))
	for i := range params {
		var key uentity.ProducerOrderBy
		switch params[i].Key {
		case "lastname":
			key = uentity.ProducerOrderByLastname
		case "firstname":
			key = uentity.ProducerOrderByFirstname
		case "storeName":
			key = uentity.ProducerOrderByStoreName
		case "email":
			key = uentity.ProducerOrderByEmail
		case "phoneNumber":
			key = uentity.ProducerOrderByPhoneNumber
		default:
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", params[i].Key, errInvalidOrderkey)
		}
		res[i] = &user.ListProducersOrder{
			Key:        key,
			OrderByASC: params[i].Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) GetProducer(ctx *gin.Context) {
	in := &user.GetProducerInput{
		ProducerID: util.GetParam(ctx, "producerId"),
	}
	producer, err := h.user.GetProducer(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer: service.NewProducer(producer).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProducer(ctx *gin.Context) {
	req := &request.CreateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.CreateProducerInput{
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		StoreName:     req.StoreName,
		ThumbnailURL:  req.ThumbnailURL,
		HeaderURL:     req.HeaderURL,
		Email:         req.Email,
		PhoneNumber:   req.PhoneNumber,
		PostalCode:    req.PostalCode,
		Prefecture:    req.Prefecture,
		City:          req.City,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
	}
	producer, err := h.user.CreateProducer(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer: service.NewProducer(producer).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProducer(ctx *gin.Context) {
	req := &request.UpdateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateProducerInput{
		ProducerID:    util.GetParam(ctx, "producerId"),
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		StoreName:     req.StoreName,
		ThumbnailURL:  req.ThumbnailURL,
		HeaderURL:     req.HeaderURL,
		PhoneNumber:   req.PhoneNumber,
		PostalCode:    req.PostalCode,
		Prefecture:    req.Prefecture,
		City:          req.City,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
	}
	if err := h.user.UpdateProducer(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) UpdateProducerEmail(ctx *gin.Context) {
	req := &request.UpdateProducerEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateProducerEmailInput{
		ProducerID: util.GetParam(ctx, "producerId"),
		Email:      req.Email,
	}
	if err := h.user.UpdateProducerEmail(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ResetProducerPassword(ctx *gin.Context) {
	in := &user.ResetProducerPasswordInput{
		ProducerID: util.GetParam(ctx, "producerId"),
	}
	if err := h.user.ResetProducerPassword(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
