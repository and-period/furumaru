package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) producerRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListProducers)
	arg.POST("", h.CreateProducer)
	arg.GET("/:producerId", h.GetProducer)
}

func (h *apiV1Handler) ListProducers(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ProducersResponse{
		Producers: []*response.Producer{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) GetProducer(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ProducerResponse{}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) CreateProducer(ctx *gin.Context) {
	req := &request.CreateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	// TODO: 詳細の実装
	res := &response.ProducerResponse{}
	ctx.JSON(http.StatusOK, res)
}
