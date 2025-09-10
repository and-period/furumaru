package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Live
// @tag.description ライブ配信関連
func (h *handler) liveRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules/:scheduleId/lives", h.authentication, h.filterAccessSchedule)

	r.GET("", h.ListLives)
	r.POST("", h.CreateLive)
	r.GET("/:liveId", h.filterAccessLive, h.GetLive)
	r.PATCH("/:liveId", h.filterAccessLive, h.UpdateLive)
	r.DELETE("/:liveId", h.filterAccessLive, h.DeleteLive)
}

func (h *handler) filterAccessLive(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			live, err := h.getLive(ctx, util.GetParam(ctx, "liveId"))
			if err != nil {
				return false, err
			}
			shop, err := h.getShop(ctx, getShopID(ctx))
			if err != nil {
				return false, err
			}
			return slices.Contains(shop.ProducerIDs, live.ProducerID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     ライブ配信一覧取得
// @Description 指定されたスケジュールのライブ配信一覧を取得します。
// @Tags        Live
// @Router      /v1/schedules/{scheduleId}/lives [get]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.LivesResponse
func (h *handler) ListLives(ctx *gin.Context) {
	scheduleID := util.GetParam(ctx, "scheduleId")
	in := &store.ListLivesInput{
		ScheduleIDs: []string{scheduleID},
		NoLimit:     true,
	}
	lives, total, err := h.store.ListLives(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	var (
		producers service.Producers
		products  service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, lives.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, lives.ProductIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.LivesResponse{
		Lives:     service.NewLives(lives).Response(),
		Producers: producers.Response(),
		Products:  products.Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ライブ配信取得
// @Description 指定されたライブ配信の詳細情報を取得します。
// @Tags        Live
// @Router      /v1/schedules/{scheduleId}/lives/{liveId} [get]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       liveId path string true "ライブ配信ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.LiveResponse
// @Failure     403 {object} util.ErrorResponse "ライブ配信の参照権限がない"
// @Failure     404 {object} util.ErrorResponse "ライブ配信が存在しない"
func (h *handler) GetLive(ctx *gin.Context) {
	liveID := util.GetParam(ctx, "liveId")
	live, err := h.getLive(ctx, liveID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	var (
		producer *service.Producer
		products service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, live.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, live.ProductIDs)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.LiveResponse{
		Live:     live.Response(),
		Producer: producer.Response(),
		Products: products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ライブ配信登録
// @Description 新しいライブ配信を登録します。
// @Tags        Live
// @Router      /v1/schedules/{scheduleId}/lives [post]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.CreateLiveRequest true "ライブ配信情報"
// @Produce     json
// @Success     200 {object} types.LiveResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateLive(ctx *gin.Context) {
	req := &types.CreateLiveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	var (
		producer *service.Producer
		products service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, req.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, req.ProductIDs)
		if len(products) != len(req.ProductIDs) {
			return fmt.Errorf("handler: unmatch products length: %w", exception.ErrInvalidArgument)
		}
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	if getAdminType(ctx).Response() == types.AdminTypeCoordinator {
		shop, err := h.getShop(ctx, getShopID(ctx))
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		if !slices.Contains(shop.ProducerIDs, producer.ID) {
			h.forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	in := &store.CreateLiveInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		ProducerID: req.ProducerID,
		ProductIDs: req.ProductIDs,
		Comment:    req.Comment,
		StartAt:    jst.ParseFromUnix(req.StartAt),
		EndAt:      jst.ParseFromUnix(req.EndAt),
	}
	live, err := h.store.CreateLive(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.LiveResponse{
		Live:     service.NewLive(live).Response(),
		Producer: producer.Response(),
		Products: products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ライブ配信更新
// @Description ライブ配信の情報を更新します。
// @Tags        Live
// @Router      /v1/schedules/{scheduleId}/lives/{liveId} [patch]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       liveId path string true "ライブ配信ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateLiveRequest true "ライブ配信情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "ライブ配信の更新権限がない"
// @Failure     404 {object} util.ErrorResponse "ライブ配信が存在しない"
func (h *handler) UpdateLive(ctx *gin.Context) {
	req := &types.UpdateLiveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	products, err := h.multiGetProducts(ctx, req.ProductIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(products) != len(req.ProductIDs) {
		h.badRequest(ctx, errors.New("handler: unmatch products length"))
		return
	}

	in := &store.UpdateLiveInput{
		LiveID:     util.GetParam(ctx, "liveId"),
		ProductIDs: req.ProductIDs,
		Comment:    req.Comment,
		StartAt:    jst.ParseFromUnix(req.StartAt),
		EndAt:      jst.ParseFromUnix(req.EndAt),
	}
	if err := h.store.UpdateLive(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ライブ配信削除
// @Description ライブ配信を削除します。
// @Tags        Live
// @Router      /v1/schedules/{scheduleId}/lives/{liveId} [delete]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       liveId path string true "ライブ配信ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "ライブ配信の削除権限がない"
// @Failure     404 {object} util.ErrorResponse "ライブ配信が存在しない"
func (h *handler) DeleteLive(ctx *gin.Context) {
	in := &store.DeleteLiveInput{
		LiveID: util.GetParam(ctx, "liveId"),
	}
	if err := h.store.DeleteLive(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) getLive(ctx context.Context, liveID string) (*service.Live, error) {
	in := &store.GetLiveInput{
		LiveID: liveID,
	}
	live, err := h.store.GetLive(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewLive(live), nil
}
