package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Coordinator
// @tag.description コーディネータ関連
func (h *handler) coordinatorRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/coordinators", h.authentication)

	r.GET("", h.ListCoordinators)
	r.POST("", h.CreateCoordinator)
	r.GET("/:coordinatorId", h.GetCoordinator)
	r.PATCH("/:coordinatorId", h.UpdateCoordinator)
	r.PATCH("/:coordinatorId/email", h.UpdateCoordinatorEmail)
	r.PATCH("/:coordinatorId/password", h.ResetCoordinatorPassword)
	r.DELETE("/:coordinatorId", h.DeleteCoordinator)
}

// @Summary     コーディネータ一覧取得
// @Description コーディネータの一覧を取得します。
// @Tags        Coordinator
// @Router      /v1/coordinators [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       username query string false "コーディネータ名(あいまい検索)(64文字以内)" example("&.コーディネータ")
// @Produce     json
// @Success     200 {object} types.CoordinatorsResponse
func (h *handler) ListCoordinators(ctx *gin.Context) {
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

	in := &user.ListCoordinatorsInput{
		Name:   util.GetQuery(ctx, "username", ""),
		Limit:  limit,
		Offset: offset,
	}
	coordinators, total, err := h.user.ListCoordinators(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(coordinators) == 0 {
		res := &types.CoordinatorsResponse{
			Coordinators: []*types.Coordinator{},
			Shops:        []*types.Shop{},
			ProductTypes: []*types.ProductType{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	shops, err := h.listShopsByCoordinatorIDs(ctx, coordinators.IDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		producerTotals map[string]int64
		productTypes   service.ProductTypes
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		aggregateIn := &user.AggregateRealatedProducersInput{
			CoordinatorIDs: coordinators.IDs(),
		}
		producerTotals, err = h.user.AggregateRealatedProducers(ectx, aggregateIn)
		return
	})
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, shops.ProductTypeIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	scoordinator := service.NewCoordinators(coordinators, shops.MapByCoordinatorID())
	scoordinator.SetProducerTotal(producerTotals)

	res := &types.CoordinatorsResponse{
		Coordinators: scoordinator.Response(),
		Shops:        shops.Response(),
		ProductTypes: productTypes.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     コーディネータ詳細取得
// @Description 指定されたIDのコーディネータの詳細情報を取得します。
// @Tags        Coordinator
// @Router      /v1/coordinators/{coordinatorId} [get]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネータID" example("coordinator-id")
// @Produce     json
// @Success     200 {object} types.CoordinatorResponse
// @Failure     404 {object} util.ErrorResponse
func (h *handler) GetCoordinator(ctx *gin.Context) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinator.ID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, shop.ProductTypeIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.CoordinatorResponse{
		Coordinator:  service.NewCoordinator(coordinator, shop).Response(),
		Shop:         shop.Response(),
		ProductTypes: productTypes.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     コーディネータ新規作成
// @Description 新しいコーディネータを作成します。
// @Tags        Coordinator
// @Router      /v1/coordinators [post]
// @Security    bearerauth
// @Accept      json
// @Produce     json
// @Param       request body types.CreateCoordinatorRequest true "コーディネータ作成リクエスト"
// @Success     200 {object} types.CoordinatorResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "コーディネータの登録権限がない"
// @Failure     409 {object} util.ErrorResponse "すでに存在するメールアドレス"
func (h *handler) CreateCoordinator(ctx *gin.Context) {
	req := &types.CreateCoordinatorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, req.ProductTypeIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(productTypes) != len(req.ProductTypeIDs) {
		h.badRequest(ctx, errors.New("handler: unmatch product types length"))
		return
	}

	in := &user.CreateCoordinatorInput{
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		MarcheName:        req.MarcheName,
		Username:          req.Username,
		Profile:           req.Profile,
		ProductTypeIDs:    req.ProductTypeIDs,
		ThumbnailURL:      req.ThumbnailURL,
		HeaderURL:         req.HeaderURL,
		PromotionVideoURL: req.PromotionVideoURL,
		BonusVideoURL:     req.BonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		BusinessDays:      req.BusinessDays,
	}
	coordinator, password, err := h.user.CreateCoordinator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinator.ID)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		// 店舗作成は非同期で行われるため、取得できない場合はレスポンスを返さないだけに留める
		h.httpError(ctx, err)
		return
	}

	res := &types.CoordinatorResponse{
		Coordinator:  service.NewCoordinator(coordinator, shop).Response(),
		Shop:         shop.Response(),
		ProductTypes: productTypes.Response(),
		Password:     password,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     コーディネータ更新
// @Description 指定されたIDのコーディネータ情報を更新します。
// @Tags        Coordinator
// @Router      /v1/coordinators/{coordinatorId} [patch]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネータID" example("coordinator-id")
// @Accept      json
// @Param       request body types.UpdateCoordinatorRequest true "コーディネータ更新リクエスト"
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "コーディネータの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "コーディネータが存在しない"
// @Failure     409 {object} util.ErrorResponse "すでに存在するメールアドレス"
func (h *handler) UpdateCoordinator(ctx *gin.Context) {
	req := &types.UpdateCoordinatorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateCoordinatorInput{
		CoordinatorID:     util.GetParam(ctx, "coordinatorId"),
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		Username:          req.Username,
		Profile:           req.Profile,
		ThumbnailURL:      req.ThumbnailURL,
		HeaderURL:         req.HeaderURL,
		PromotionVideoURL: req.PromotionVideoURL,
		BonusVideoURL:     req.BonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
	}
	if err := h.user.UpdateCoordinator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     コーディネータメールアドレス更新
// @Description 指定されたIDのコーディネータのメールアドレスを更新します。
// @Tags        Coordinator
// @Router      /v1/coordinators/{coordinatorId}/email [patch]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネータID" example("coordinator-id")
// @Accept      json
// @Param       request body types.UpdateCoordinatorEmailRequest true "メールアドレス更新リクエスト"
// @Success     204
// @Failure     403 {object} util.ErrorResponse "コーディネータの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "コーディネータが存在しない"
func (h *handler) UpdateCoordinatorEmail(ctx *gin.Context) {
	req := &types.UpdateCoordinatorEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateCoordinatorEmailInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
		Email:         req.Email,
	}
	if err := h.user.UpdateCoordinatorEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     コーディネータパスワードリセット
// @Description 指定されたIDのコーディネータのパスワードをリセットします。
// @Tags        Coordinator
// @Router      /v1/coordinators/{coordinatorId}/password [patch]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネータID" example("coordinator-id")
// @Success     204
// @Failure     403 {object} util.ErrorResponse "コーディネータの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "コーディネータが存在しない"
func (h *handler) ResetCoordinatorPassword(ctx *gin.Context) {
	in := &user.ResetCoordinatorPasswordInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	if err := h.user.ResetCoordinatorPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     コーディネータ削除
// @Description 指定されたIDのコーディネータを削除します。
// @Tags        Coordinator
// @Router      /v1/coordinators/{coordinatorId} [delete]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネータID" example("coordinator-id")
// @Success     204
// @Failure     403 {object} util.ErrorResponse "コーディネータの削除権限がない"
// @Failure     404 {object} util.ErrorResponse "コーディネータが存在しない"
func (h *handler) DeleteCoordinator(ctx *gin.Context) {
	in := &user.DeleteCoordinatorInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	if err := h.user.DeleteCoordinator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetCoordinators(ctx context.Context, coordinatorIDs []string) (service.Coordinators, error) {
	if len(coordinatorIDs) == 0 {
		return service.Coordinators{}, nil
	}
	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: coordinatorIDs,
	}
	coordinators, err := h.user.MultiGetCoordinators(ctx, in)
	if err != nil {
		return nil, err
	}
	if len(coordinators) == 0 {
		return service.Coordinators{}, nil
	}
	shops, err := h.listShopsByCoordinatorIDs(ctx, coordinatorIDs)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinators(coordinators, shops.MapByCoordinatorID()), nil
}

func (h *handler) multiGetCoordinatorsWithDeleted(ctx context.Context, coordinatorIDs []string) (service.Coordinators, error) {
	if len(coordinatorIDs) == 0 {
		return service.Coordinators{}, nil
	}
	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: coordinatorIDs,
		WithDeleted:    true,
	}
	coordinators, err := h.user.MultiGetCoordinators(ctx, in)
	if err != nil {
		return nil, err
	}
	if len(coordinators) == 0 {
		return service.Coordinators{}, nil
	}
	shops, err := h.listShopsByCoordinatorIDs(ctx, coordinatorIDs)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinators(coordinators, shops.MapByCoordinatorID()), nil
}

func (h *handler) getCoordinator(ctx context.Context, coordinatorID string) (*service.Coordinator, error) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: coordinatorID,
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		return nil, err
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinator.ID)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinator(coordinator, shop), nil
}

func (h *handler) getCoordinatorWithDeleted(ctx context.Context, coordinatorID string) (*service.Coordinator, error) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: coordinatorID,
		WithDeleted:   true,
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		return nil, err
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinator.ID)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinator(coordinator, shop), nil
}
