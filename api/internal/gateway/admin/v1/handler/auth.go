package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

// @tag.name        Auth
// @tag.description 認証関連
func (h *handler) authRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/auth")

	r.GET("", h.GetAuth)
	r.POST("", h.SignIn)
	r.DELETE("", h.SignOut)
	r.GET("/providers", h.authentication, h.ListAuthProviders)
	r.GET("/google", h.authentication, h.AuthGoogleAccount)
	r.POST("/google", h.authentication, h.ConnectGoogleAccount)
	r.GET("/line", h.authentication, h.AuthLINEAccount)
	r.POST("/line", h.authentication, h.ConnectLINEAccount)
	r.POST("/refresh-token", h.RefreshAuthToken)
	r.POST("/device", h.authentication, h.RegisterDevice)
	r.PATCH("/email", h.authentication, h.UpdateAuthEmail)
	r.POST("/email/verified", h.VerifyAuthEmail)
	r.PATCH("/password", h.authentication, h.UpdateAuthPassword)
	r.POST("/forgot-password", h.ForgotAuthPassword)
	r.POST("/forgot-password/verified", h.ResetAuthPassword)
	r.GET("/user", h.authentication, h.GetAuthUser)
	r.GET("/coordinator", h.authentication, h.GetAuthCoordinator)
	r.PATCH("/coordinator", h.authentication, h.UpdateAuthCoordinator)
	r.GET("/coordinator/shippings", h.authentication, h.GetAuthShipping)
	r.PATCH("/coordinator/shippings", h.authentication, h.UpsertAuthShipping)
}

// @Summary     トークン検証
// @Description 認証トークンを検証し、認証情報を取得します。
// @Tags        Auth
// @Router      /v1/auth [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) GetAuth(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}

	in := &user.GetAdminAuthInput{
		AccessToken: token,
	}
	auth, err := h.user.GetAdminAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

type authUser interface {
	AuthUser() *service.AuthUser
}

// @Summary     管理者情報取得
// @Description ログイン中の管理者情報を取得します。
// @Tags        Auth
// @Router      /v1/auth/user [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.AuthUserResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) GetAuthUser(ctx *gin.Context) {
	adminID := getAdminID(ctx)
	var (
		auth authUser
		err  error
	)
	switch getAdminType(ctx).Response() {
	case types.AdminTypeAdministrator:
		auth, err = h.getAdministrator(ctx, adminID)
	case types.AdminTypeCoordinator:
		auth, err = h.getCoordinator(ctx, adminID)
	case types.AdminTypeProducer:
		auth, err = h.getProducer(ctx, adminID)
	default:
		h.forbidden(ctx, errors.New("handler: unknown admin role"))
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthUserResponse{
		AuthUser: auth.AuthUser().Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     サインイン
// @Description メールアドレスとパスワードでサインインします。
// @Tags        Auth
// @Router      /v1/auth [post]
// @Accept      json
// @Param       request body types.SignInRequest true "サインイン"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) SignIn(ctx *gin.Context) {
	req := &types.SignInRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.SignInAdminInput{
		Key:      req.Username,
		Password: req.Password,
	}
	auth, err := h.user.SignInAdmin(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     認証済みプロバイダ一覧の取得
// @Description 連携済みの外部認証プロバイダ一覧を取得します。
// @Tags        Auth
// @Router      /v1/auth/providers [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.AuthProvidersResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) ListAuthProviders(ctx *gin.Context) {
	in := &user.ListAdminAuthProvidersInput{
		AdminID: getAdminID(ctx),
	}
	providers, err := h.user.ListAdminAuthProviders(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthProvidersResponse{
		Providers: service.NewAuthProviders(providers).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     Google認証用URLの発行
// @Description Googleアカウント連携用の認証URLを発行します。
// @Tags        Auth
// @Router      /v1/auth/google [get]
// @Security    bearerauth
// @Param       state query string true "CSRF対策用のstate" example("xxxxxxxxxx")
// @Param       redirectUri query string false "認証後のリダイレクト先（変更したいときのみ指定）" example("https://example.com")
// @Produce     json
// @Success     200 {object} types.AuthGoogleAccountResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     412 {object} util.ErrorResponse "すでに連携済み"
func (h *handler) AuthGoogleAccount(ctx *gin.Context) {
	in := &user.InitialGoogleAdminAuthInput{
		AdminID:     getAdminID(ctx),
		State:       util.GetQuery(ctx, "state", ""),
		RedirectURI: util.GetQuery(ctx, "redirectUri", ""),
	}
	authURL, err := h.user.InitialGoogleAdminAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthGoogleAccountResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     Googleアカウントの連携
// @Description Googleアカウントを連携します。
// @Tags        Auth
// @Router      /v1/auth/google [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.ConnectGoogleAccountRequest true "連携リクエスト"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) ConnectGoogleAccount(ctx *gin.Context) {
	req := &types.ConnectGoogleAccountRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.ConnectGoogleAdminAuthInput{
		AdminID:     getAdminID(ctx),
		Code:        req.Code,
		Nonce:       req.Nonce,
		RedirectURI: req.RedirectURI,
	}
	if err := h.user.ConnectGoogleAdminAuth(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     LINE認証用URLの発行
// @Description LINEアカウント連携用の認証URLを発行します。
// @Tags        Auth
// @Router      /v1/auth/line [get]
// @Security    bearerauth
// @Param       state query string true "CSRF対策用のstate" example("xxxxxxxxxx")
// @Param       redirectUri query string false "認証後のリダイレクト先（変更したいときのみ指定）" example("https://example.com")
// @Produce     json
// @Success     200 {object} types.AuthLINEAccountResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     412 {object} util.ErrorResponse "すでに連携済み"
func (h *handler) AuthLINEAccount(ctx *gin.Context) {
	in := &user.InitialLINEAdminAuthInput{
		AdminID:     getAdminID(ctx),
		State:       util.GetQuery(ctx, "state", ""),
		RedirectURI: util.GetQuery(ctx, "redirectUri", ""),
	}
	authURL, err := h.user.InitialLINEAdminAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthLINEAccountResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     LINEアカウントの連携
// @Description LINEアカウントを連携します。
// @Tags        Auth
// @Router      /v1/auth/line [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.ConnectLINEAccountRequest true "連携リクエスト"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) ConnectLINEAccount(ctx *gin.Context) {
	req := &types.ConnectLINEAccountRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.ConnectLINEAdminAuthInput{
		AdminID:     getAdminID(ctx),
		Code:        req.Code,
		Nonce:       req.Nonce,
		RedirectURI: req.RedirectURI,
	}
	if err := h.user.ConnectLINEAdminAuth(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     サインアウト
// @Description サインアウトします。
// @Tags        Auth
// @Router      /v1/auth [delete]
// @Security    bearerauth
// @Produce     json
// @Success     204
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) SignOut(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}

	in := &user.SignOutAdminInput{
		AccessToken: token,
	}
	if err := h.user.SignOutAdmin(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     トークン更新
// @Description リフレッシュトークンを使用してアクセストークンを更新します。
// @Tags        Auth
// @Router      /v1/auth/refresh-token [post]
// @Accept      json
// @Param       request body types.RefreshAuthTokenRequest true "トークン更新"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) RefreshAuthToken(ctx *gin.Context) {
	req := &types.RefreshAuthTokenRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.RefreshAdminTokenInput{
		RefreshToken: req.RefreshToken,
	}
	auth, err := h.user.RefreshAdminToken(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     デバイストークン登録
// @Description プッシュ通知用のデバイストークンを登録します。
// @Tags        Auth
// @Router      /v1/auth/device [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.RegisterAuthDeviceRequest true "デバイストークン"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) RegisterDevice(ctx *gin.Context) {
	req := &types.RegisterAuthDeviceRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.RegisterAdminDeviceInput{
		AdminID: getAdminID(ctx),
		Device:  req.Device,
	}
	if err := h.user.RegisterAdminDevice(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     メールアドレス更新
// @Description ログイン中のユーザーのメールアドレスを更新します。
// @Tags        Auth
// @Router      /v1/auth/email [patch]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpdateAuthEmailRequest true "メールアドレス"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     409 {object} util.ErrorResponse "現在すでに存在するメールアドレス"
// @Failure     412 {object} util.ErrorResponse "変更後のメールアドレスが変更前と同じ"
func (h *handler) UpdateAuthEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &types.UpdateAuthEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdminEmailInput{
		AccessToken: token,
		Email:       req.Email,
	}
	if err := h.user.UpdateAdminEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     メールアドレス更新 - コード検証
// @Description メールアドレス更新用の検証コードを確認します。
// @Tags        Auth
// @Router      /v1/auth/email/verified [post]
// @Accept      json
// @Param       request body types.VerifyAuthEmailRequest true "検証コード"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) VerifyAuthEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &types.VerifyAuthEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.VerifyAdminEmailInput{
		AccessToken: token,
		VerifyCode:  req.VerifyCode,
	}
	if err := h.user.VerifyAdminEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     パスワード更新
// @Description ログイン中のユーザーのパスワードを更新します。
// @Tags        Auth
// @Router      /v1/auth/password [patch]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpdateAuthPasswordRequest true "パスワード"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) UpdateAuthPassword(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &types.UpdateAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdminPasswordInput{
		AccessToken:          token,
		OldPassword:          req.OldPassword,
		NewPassword:          req.NewPassword,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.UpdateAdminPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     パスワードリセット
// @Description パスワードリセット用のメールを送信します。
// @Tags        Auth
// @Router      /v1/auth/forgot-password [post]
// @Accept      json
// @Param       request body types.ForgotAuthPasswordRequest true "メールアドレス"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ForgotAuthPassword(ctx *gin.Context) {
	req := &types.ForgotAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.ForgotAdminPasswordInput{
		Email: req.Email,
	}
	if err := h.user.ForgotAdminPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     パスワードリセット - コード検証
// @Description パスワードリセット用の検証コードを確認し、新しいパスワードを設定します。
// @Tags        Auth
// @Router      /v1/auth/forgot-password/verified [post]
// @Accept      json
// @Param       request body types.ResetAuthPasswordRequest true "パスワードリセット"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ResetAuthPassword(ctx *gin.Context) {
	req := &types.ResetAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.VerifyAdminPasswordInput{
		Email:                req.Email,
		VerifyCode:           req.VerifyCode,
		NewPassword:          req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.VerifyAdminPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     自身のコーディネータ情報取得
// @Description ログイン中のコーディネータの詳細情報を取得します。
// @Tags        Auth
// @Router      /v1/auth/coordinator [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.CoordinatorResponse
// @Failure     404 {object} util.ErrorResponse "コーディネータが存在しない"
func (h *handler) GetAuthCoordinator(ctx *gin.Context) {
	if getAdminType(ctx).Response() != types.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	in := &user.GetCoordinatorInput{
		CoordinatorID: getAdminID(ctx),
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

// @Summary     自身のコーディネータ情報更新
// @Description ログイン中のコーディネータの情報を更新します。
// @Tags        Auth
// @Router      /v1/auth/coordinator [patch]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpdateCoordinatorRequest true "コーディネータ情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) UpdateAuthCoordinator(ctx *gin.Context) {
	if getAdminType(ctx).Response() != types.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	req := &types.UpdateCoordinatorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateCoordinatorInput{
		CoordinatorID:     getAdminID(ctx),
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

// @Summary     自身の配送設定取得
// @Description ログイン中のコーディネータの配送設定を取得します。
// @Tags        Auth
// @Router      /v1/auth/coordinator/shippings [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.ShippingResponse
func (h *handler) GetAuthShipping(ctx *gin.Context) {
	if getAdminType(ctx).Response() != types.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	in := &store.GetShippingByShopIDInput{
		ShopID: getShopID(ctx),
	}
	shipping, err := h.store.GetShippingByShopID(ctx, in)
	if errors.Is(err, exception.ErrNotFound) {
		// 配送設定の登録をしていない場合、デフォルト設定を返却する
		in := &store.GetDefaultShippingInput{}
		shipping, err = h.store.GetDefaultShipping(ctx, in)
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.ShippingResponse{
		Shipping: service.NewShipping(shipping).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     自身の配送設定更新
// @Description ログイン中のコーディネータの配送設定を更新します。
// @Tags        Auth
// @Router      /v1/auth/coordinator/shippings [patch]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpsertShippingRequest true "配送設定"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) UpsertAuthShipping(ctx *gin.Context) {
	if getAdminType(ctx).Response() != types.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	req := &types.UpsertShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpsertShippingInput{
		CoordinatorID:     getAdminID(ctx),
		Box60Rates:        h.newShippingRatesForUpsert(req.Box60Rates),
		Box60Frozen:       req.Box60Frozen,
		Box80Rates:        h.newShippingRatesForUpsert(req.Box80Rates),
		Box80Frozen:       req.Box80Frozen,
		Box100Rates:       h.newShippingRatesForUpsert(req.Box100Rates),
		Box100Frozen:      req.Box100Frozen,
		HasFreeShipping:   req.HasFreeShipping,
		FreeShippingRates: req.FreeShippingRates,
	}
	if err := h.store.UpsertShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
