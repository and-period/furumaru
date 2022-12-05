package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) producerRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListProducers)
	arg.POST("", h.CreateProducer)
	arg.GET("/:producerId", h.filterAccessProducer, h.GetProducer)
	arg.PATCH("/:producerId", h.filterAccessProducer, h.UpdateProducer)
	arg.PATCH("/:producerId/email", h.filterAccessProducer, h.UpdateProducerEmail)
	arg.PATCH("/:producerId/password", h.filterAccessProducer, h.ResetProducerPassword)
	arg.DELETE("/:producerId", h.filterAccessProducer, h.DeleteProducer)
}

func (h *handler) filterAccessProducer(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			producer, err := h.getProducer(ctx, util.GetParam(ctx, "producerId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, producer.CoordinatorID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Next()
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

	in := &user.ListProducersInput{
		Limit:  limit,
		Offset: offset,
	}
	h.addlistProducerFilters(ctx, in)
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

func (h *handler) addlistProducerFilters(ctx *gin.Context, in *user.ListProducersInput) {
	strs := util.GetQueryStrings(ctx, "filters")
	filters := set.NewEmpty[string](len(strs))
	filters.Add(strs...)
	if filters.Contains("unrelated") { // 未関連状態の生産者のみ取得
		in.OnlyUnrelated = true
		return
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		in.CoordinatorID = getAdminID(ctx)
	}
}

func (h *handler) GetProducer(ctx *gin.Context) {
	producer, err := h.getProducer(ctx, util.GetParam(ctx, "producerId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer: producer.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProducer(ctx *gin.Context) {
	req := &request.CreateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	var thumbnailURL, headerURL string
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if req.ThumbnailURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.ThumbnailURL,
		}
		thumbnailURL, err = h.media.UploadProducerThumbnail(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		if req.HeaderURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.HeaderURL,
		}
		headerURL, err = h.media.UploadProducerHeader(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	in := &user.CreateProducerInput{
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		StoreName:     req.StoreName,
		ThumbnailURL:  thumbnailURL,
		HeaderURL:     headerURL,
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

	var thumbnailURL, headerURL string
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if req.ThumbnailURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.ThumbnailURL,
		}
		thumbnailURL, err = h.media.UploadProducerThumbnail(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		if req.HeaderURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.HeaderURL,
		}
		headerURL, err = h.media.UploadProducerHeader(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	in := &user.UpdateProducerInput{
		ProducerID:    util.GetParam(ctx, "producerId"),
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		StoreName:     req.StoreName,
		ThumbnailURL:  thumbnailURL,
		HeaderURL:     headerURL,
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

func (h *handler) DeleteProducer(ctx *gin.Context) {
	in := &user.DeleteProducerInput{
		ProducerID: util.GetParam(ctx, "producerId"),
	}
	if err := h.user.DeleteProducer(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) multiGetProducers(ctx context.Context, producerIDs []string) (service.Producers, error) {
	in := &user.MultiGetProducersInput{
		ProducerIDs: producerIDs,
	}
	producers, err := h.user.MultiGetProducers(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducers(producers), nil
}

func (h *handler) getProducer(ctx context.Context, producerID string) (*service.Producer, error) {
	in := &user.GetProducerInput{
		ProducerID: producerID,
	}
	producer, err := h.user.GetProducer(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducer(producer), nil
}

func (h *handler) getProducersByCoordinatorID(ctx context.Context, coordinatorID string) (service.Producers, error) {
	in := &user.ListProducersInput{
		CoordinatorID: coordinatorID,
	}
	producers, _, err := h.user.ListProducers(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducers(producers), nil
}
