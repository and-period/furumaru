package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *apiV1Handler) storeRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListStores)
	arg.POST("", h.CreateStore)
	arg.GET("/:storeId", h.GetStore)
	arg.PATCH("/:storeId", h.UpdateStore)
}

func (h *apiV1Handler) ListStores(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

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

	in := &store.ListStoresInput{
		Limit:  limit,
		Offset: offset,
	}
	stores, err := h.store.ListStores(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.StoresResponse{
		Stores: service.NewStores(stores).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) GetStore(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	storeID, err := util.GetParamInt64(ctx, "storeId")
	if err != nil {
		badRequest(ctx, err)
		return
	}

	var (
		sstore  *sentity.Store
		sstaffs sentity.Staffs
		uadmins uentity.Admins
	)

	eg, ectx := errgroup.WithContext(c)
	eg.Go(func() (err error) {
		in := &store.GetStoreInput{StoreID: storeID}
		sstore, err = h.store.GetStore(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		staffsIn := &store.ListStaffsByStoreIDInput{StoreID: storeID}
		sstaffs, err = h.store.ListStaffsByStoreID(ectx, staffsIn)
		if err != nil || len(sstaffs) == 0 {
			return
		}
		adminsIn := &user.MultiGetAdminsInput{AdminIDs: sstaffs.UserIDs()}
		uadmins, err = h.user.MultiGetAdmins(ectx, adminsIn)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	staffs := service.NewStaffs(sstaffs, uadmins.Map())

	res := &response.StoreResponse{
		Store: service.NewStore(sstore, staffs).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) CreateStore(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.CreateStoreRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.CreateStoreInput{
		Name: req.Name,
	}
	sstore, err := h.store.CreateStore(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.StoreResponse{
		Store: service.NewStore(sstore, nil).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) UpdateStore(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	storeID, err := util.GetParamInt64(ctx, "storeId")
	if err != nil {
		badRequest(ctx, err)
		return
	}

	req := &request.UpdateStoreRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.UpdateStoreInput{
		StoreID:      storeID,
		Name:         req.Name,
		ThumbnailURL: req.ThumbnailURL,
	}
	if err := h.store.UpdateStore(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
