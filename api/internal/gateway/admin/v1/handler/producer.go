package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) producerRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/producers", h.authentication)

	r.GET("", h.ListProducers)
	r.POST("", h.CreateProducer)
	r.GET("/:producerId", h.filterAccessProducer, h.GetProducer)
	r.PATCH("/:producerId", h.filterAccessProducer, h.UpdateProducer)
	r.DELETE("/:producerId", h.filterAccessProducer, h.DeleteProducer)
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
		h.httpError(ctx, err)
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
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.ListProducersInput{
		Username: util.GetQuery(ctx, "username", ""),
		Limit:    limit,
		Offset:   offset,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		in.CoordinatorID = getAdminID(ctx)
	}
	producers, total, err := h.user.ListProducers(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinators, err := h.multiGetCoordinators(ctx, producers.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducersResponse{
		Producers:    service.NewProducers(producers).Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetProducer(ctx *gin.Context) {
	producer, err := h.getProducer(ctx, util.GetParam(ctx, "producerId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, producer.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer:    producer.Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProducer(ctx *gin.Context) {
	req := &request.CreateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		if !currentAdmin(ctx, req.CoordinatorID) {
			h.forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	var thumbnailURL, headerURL, promotionVideoURL, bonusVideoURL string
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
	eg.Go(func() (err error) {
		if req.PromotionVideoURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.PromotionVideoURL,
		}
		promotionVideoURL, err = h.media.UploadProducerPromotionVideo(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		if req.BonusVideoURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.BonusVideoURL,
		}
		bonusVideoURL, err = h.media.UploadProducerBonusVideo(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &user.CreateProducerInput{
		CoordinatorID:     req.CoordinatorID,
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		Username:          req.Username,
		Profile:           req.Profile,
		ThumbnailURL:      thumbnailURL,
		HeaderURL:         headerURL,
		PromotionVideoURL: promotionVideoURL,
		BonusVideoURL:     bonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
	}
	producer, err := h.user.CreateProducer(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, producer.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer:    service.NewProducer(producer).Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProducer(ctx *gin.Context) {
	req := &request.UpdateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	var thumbnailURL, headerURL, promotionVideoURL, bonusVideoURL string
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
	eg.Go(func() (err error) {
		if req.PromotionVideoURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.PromotionVideoURL,
		}
		promotionVideoURL, err = h.media.UploadProducerPromotionVideo(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		if req.BonusVideoURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.BonusVideoURL,
		}
		bonusVideoURL, err = h.media.UploadProducerBonusVideo(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &user.UpdateProducerInput{
		ProducerID:        util.GetParam(ctx, "producerId"),
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		Username:          req.Username,
		Profile:           req.Profile,
		ThumbnailURL:      thumbnailURL,
		HeaderURL:         headerURL,
		PromotionVideoURL: promotionVideoURL,
		BonusVideoURL:     bonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
	}
	if err := h.user.UpdateProducer(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteProducer(ctx *gin.Context) {
	in := &user.DeleteProducerInput{
		ProducerID: util.GetParam(ctx, "producerId"),
	}
	if err := h.user.DeleteProducer(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) multiGetProducers(ctx context.Context, producerIDs []string) (service.Producers, error) {
	if len(producerIDs) == 0 {
		return service.Producers{}, nil
	}
	in := &user.MultiGetProducersInput{
		ProducerIDs: producerIDs,
	}
	producers, err := h.user.MultiGetProducers(ctx, in)
	if err != nil {
		return nil, err
	}
	if len(producers) == 0 {
		return service.Producers{}, nil
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
	if coordinatorID == "" {
		return service.Producers{}, nil
	}
	in := &user.ListProducersInput{
		CoordinatorID: coordinatorID,
	}
	producers, _, err := h.user.ListProducers(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducers(producers), nil
}
