package handler

import (
	"net/http"
	"strconv"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/and-period/marche/api/internal/gateway/admin/v1/service"
	"github.com/and-period/marche/api/internal/gateway/util"
	sentity "github.com/and-period/marche/api/internal/store/entity"
	store "github.com/and-period/marche/api/internal/store/service"
	uentity "github.com/and-period/marche/api/internal/user/entity"
	user "github.com/and-period/marche/api/internal/user/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *apiV1Handler) storeRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListStores)
	arg.GET("/:storeId", h.GetStore)
}

func (h *apiV1Handler) ListStores(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	const (
		defaultLimit  = "20"
		defaultOffset = "0"
	)

	limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", defaultLimit), 10, 64)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := strconv.ParseInt(ctx.DefaultQuery("offset", defaultOffset), 10, 64)
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

	storeID, err := strconv.ParseInt(ctx.Param("storeId"), 10, 64)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	var (
		sstore  *sentity.Store
		sstaffs sentity.Staffs
		ushops  uentity.Shops
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
		shopsIn := &user.MultiGetShopsInput{ShopIDs: sstaffs.UserIDs()}
		ushops, err = h.user.MultiGetShops(ectx, shopsIn)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	staffs := service.NewStaffs(sstaffs, ushops.Map())

	res := &response.StoreResponse{
		Store: service.NewStore(sstore, staffs).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
