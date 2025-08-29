package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

// @tag.name        Address
// @tag.description 住所関連
func (h *handler) addressRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/addresses", h.authentication)

	r.GET("", h.ListAddresses)
	r.POST("", h.CreateAddress)
	r.GET("/:addressId", h.GetAddress)
	r.PATCH("/:addressId", h.UpdateAddress)
	r.DELETE("/:addressId", h.DeleteAddress)
}

// @Summary     住所一覧取得
// @Description ユーザーの登録済み住所一覧を取得します。
// @Tags        Address
// @Router      /addresses [get]
// @Security    bearerauth
// @Param       limit query int64 false "取得上限数(max:200)" default(20)
// @Param       offset query int64 false "取得開始位置(min:0)" default(0)
// @Produce     json
// @Success     200 {object} response.AddressesResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) ListAddresses(ctx *gin.Context) {
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

	in := &user.ListAddressesInput{
		UserID: h.getUserID(ctx),
		Limit:  limit,
		Offset: offset,
	}
	addresses, total, err := h.user.ListAddresses(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AddressesResponse{
		Addresses: service.NewAddresses(addresses).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     住所詳細取得
// @Description 指定されたIDの住所詳細を取得します。
// @Tags        Address
// @Router      /addresses/{addressId} [get]
// @Security    bearerauth
// @Param       addressId path string true "住所ID"
// @Produce     json
// @Success     200 {object} response.AddressResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "他のユーザーのアドレス情報"
// @Failure     404 {object} util.ErrorResponse "アドレスが存在しない"
func (h *handler) GetAddress(ctx *gin.Context) {
	in := &user.GetAddressInput{
		AddressID: util.GetParam(ctx, "addressId"),
		UserID:    h.getUserID(ctx),
	}
	address, err := h.user.GetAddress(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AddressResponse{
		Address: service.NewAddress(address).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     住所登録
// @Description 新しい住所を登録します。
// @Tags        Address
// @Router      /addresses [post]
// @Security    bearerauth
// @Accept      json
// @Produce     json
// @Param       body body request.CreateAddressRequest true "住所情報"
// @Success     200 {object} response.AddressResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) CreateAddress(ctx *gin.Context) {
	req := &request.CreateAddressRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.CreateAddressInput{
		UserID:         h.getUserID(ctx),
		Lastname:       req.Lastname,
		Firstname:      req.Firstname,
		LastnameKana:   req.LastnameKana,
		FirstnameKana:  req.FirstnameKana,
		PostalCode:     req.PostalCode,
		PrefectureCode: req.PrefectureCode,
		City:           req.City,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		PhoneNumber:    req.PhoneNumber,
		IsDefault:      req.IsDefault,
	}
	address, err := h.user.CreateAddress(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AddressResponse{
		Address: service.NewAddress(address).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     住所更新
// @Description 指定されたIDの住所情報を更新します。
// @Tags        Address
// @Router      /addresses/{addressId} [patch]
// @Security    bearerauth
// @Accept      json
// @Param       addressId path string true "住所ID"
// @Param       body body request.UpdateAddressRequest true "更新する住所情報"
// @Success     204 "更新成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "他のユーザーのアドレス情報"
// @Failure     404 {object} util.ErrorResponse "アドレスが存在しない"
func (h *handler) UpdateAddress(ctx *gin.Context) {
	req := &request.UpdateAddressRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateAddressInput{
		AddressID:      util.GetParam(ctx, "addressId"),
		UserID:         h.getUserID(ctx),
		Lastname:       req.Lastname,
		Firstname:      req.Firstname,
		LastnameKana:   req.LastnameKana,
		FirstnameKana:  req.FirstnameKana,
		PostalCode:     req.PostalCode,
		PrefectureCode: req.PrefectureCode,
		City:           req.City,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		PhoneNumber:    req.PhoneNumber,
		IsDefault:      req.IsDefault,
	}
	if err := h.user.UpdateAddress(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     住所削除
// @Description 指定されたIDの住所を削除します。
// @Tags        Address
// @Router      /addresses/{addressId} [delete]
// @Security    bearerauth
// @Param       addressId path string true "住所ID"
// @Success     204 "削除成功"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "他のユーザーのアドレス情報"
// @Failure     404 {object} util.ErrorResponse "アドレスが存在しない"
func (h *handler) DeleteAddress(ctx *gin.Context) {
	in := &user.DeleteAddressInput{
		AddressID: util.GetParam(ctx, "addressId"),
		UserID:    h.getUserID(ctx),
	}
	if err := h.user.DeleteAddress(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetAddressesByRevision(ctx context.Context, revisionIDs []int64) (service.Addresses, error) {
	if len(revisionIDs) == 0 {
		return service.Addresses{}, nil
	}
	in := &user.MultiGetAddressesByRevisionInput{
		AddressRevisionIDs: revisionIDs,
	}
	addresses, err := h.user.MultiGetAddressesByRevision(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAddresses(addresses), nil
}
