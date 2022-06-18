package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) producerRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListProducers)
	arg.POST("", h.CreateProducer)
	arg.GET("/:producerId", h.GetProducer)
}

func (h *apiV1Handler) ListProducers(ctx *gin.Context) {
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

	in := &user.ListProducersInput{
		Limit:  limit,
		Offset: offset,
	}
	producers, err := h.user.ListProducers(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProducersResponse{
		Producers: service.NewProducers(producers).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) GetProducer(ctx *gin.Context) {
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

func (h *apiV1Handler) CreateProducer(ctx *gin.Context) {
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
