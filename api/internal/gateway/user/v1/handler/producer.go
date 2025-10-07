package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Producer
// @tag.description 生産者関連
func (h *handler) producerRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/producers")

	r.GET("", h.ListProducers)
	r.GET("/:producerId", h.GetProducer)
}

// @Summary     生産者一覧取得
// @Description 生産者の一覧を取得します。
// @Tags        Producer
// @Router      /producers [get]
// @Param       limit query int64 false "取得件数" default(20)
// @Param       offset query int64 false "取得開始位置" default(0)
// @Produce     json
// @Success     200 {object} types.ProducersResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
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

	producersIn := &user.ListProducersInput{
		Limit:  limit,
		Offset: offset,
	}
	producers, total, err := h.user.ListProducers(ctx, producersIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(producers) == 0 {
		res := &types.ProducersResponse{
			Producers: []*types.Producer{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	shops, err := h.listShopsByProducerIDs(ctx, producers.IDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ProducersResponse{
		Producers: service.NewProducers(producers, shops.GroupByProducerID()).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     生産者詳細取得
// @Description 生産者の詳細情報を取得します。
// @Tags        Producer
// @Router      /producers/{producerId} [get]
// @Param       producerId path string true "生産者ID"
// @Produce     json
// @Success     200 {object} types.ProducerResponse
// @Failure     404 {object} util.ErrorResponse "生産者が見つからない"
func (h *handler) GetProducer(ctx *gin.Context) {
	producer, err := h.getProducer(ctx, util.GetParam(ctx, "producerId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		lives       service.LiveSummaries
		archives    service.ArchiveSummaries
		products    service.Products
		experiences service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &listLiveSummariesParams{
			producerID: producer.ID,
			noLimit:    true,
		}
		lives, _, err = h.listLiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &listArchiveSummariesParams{
			producerID: producer.ID,
			noLimit:    true,
		}
		archives, _, err = h.listArchiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		in := &store.ListProductsInput{
			ProducerID:       producer.ID,
			Scopes:           []entity.ProductScope{entity.ProductScopePublic},
			ExcludeOutOfSale: true,
			NoLimit:          true,
		}
		products, err = h.listProducts(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.ListExperiencesInput{
			ProducerID:      producer.ID,
			OnlyPublished:   true,
			ExcludeFinished: true,
			ExcludeDeleted:  true,
			NoLimit:         true,
		}
		experiences, err = h.listExperiences(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ProducerResponse{
		Producer:    producer.Response(),
		Lives:       lives.Response(),
		Archives:    archives.Response(),
		Products:    products.Response(),
		Experiences: experiences.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) listProducersByCoordinatorID(ctx context.Context, coordinatorID string) (service.Producers, error) {
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
	if len(producers) == 0 {
		return service.Producers{}, nil
	}
	shopsIn := &store.ListShopsInput{
		CoordinatorIDs: []string{coordinatorID},
		ProducerIDs:    producers.IDs(),
		NoLimit:        true,
	}
	shops, _, err := h.store.ListShops(ctx, shopsIn)
	if err != nil {
		return nil, err
	}
	return service.NewProducers(producers, shops.GroupByProducerID()), nil
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
	shopsIn := &store.ListShopsInput{
		ProducerIDs: producerIDs,
		NoLimit:     true,
	}
	shops, _, err := h.store.ListShops(ctx, shopsIn)
	if err != nil {
		return nil, err
	}
	return service.NewProducers(producers, shops.GroupByProducerID()), nil
}

func (h *handler) getProducer(ctx context.Context, producerID string) (*service.Producer, error) {
	in := &user.GetProducerInput{
		ProducerID: producerID,
	}
	producer, err := h.user.GetProducer(ctx, in)
	if err != nil {
		return nil, err
	}
	shopsIn := &store.ListShopsInput{
		ProducerIDs: []string{producerID},
		NoLimit:     true,
	}
	shops, _, err := h.store.ListShops(ctx, shopsIn)
	if err != nil {
		return nil, err
	}
	return service.NewProducer(producer, shops), nil
}
