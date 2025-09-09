package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
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
	r.POST("/refresh-token", h.RefreshAuthToken)
	r.PATCH("/password", h.authentication, h.UpdateAuthPassword)
	r.POST("/forgot-password", h.ForgotAuthPassword)
	r.POST("/forgot-password/verified", h.ResetAuthPassword)
	r.GET("/google", h.AuthGoogleAccount)
	r.GET("/line", h.AuthLINEAccount)
}

// @Summary     トークン検証
// @Description 認証トークンを検証し、認証情報を取得します。
// @Tags        Auth
// @Router      /auth [get]
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

	in := &user.GetUserAuthInput{
		AccessToken: token,
	}
	auth, err := h.user.GetUserAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     サインイン
// @Description ユーザー名/メールアドレスとパスワードでサインインします。
// @Tags        Auth
// @Router      /auth [post]
// @Accept      json
// @Param       request body types.SignInRequest true "サインイン"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) SignIn(ctx *gin.Context) {
	req := &types.SignInRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.SignInUserInput{
		Key:      req.Username,
		Password: req.Password,
	}
	auth, err := h.user.SignInUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     サインアウト
// @Description ふるマルからサインアウトします。
// @Tags        Auth
// @Router      /auth [delete]
// @Security    bearerauth
// @Produce     json
// @Success     204
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) SignOut(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
	}

	in := &user.SignOutUserInput{
		AccessToken: token,
	}
	if err := h.user.SignOutUser(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     トークンリフレッシュ
// @Description リフレッシュトークンを使用してアクセストークンを更新します。
// @Tags        Auth
// @Router      /auth/refresh-token [post]
// @Accept      json
// @Param       request body types.RefreshAuthTokenRequest true "トークンリフレッシュ"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) RefreshAuthToken(ctx *gin.Context) {
	req := &types.RefreshAuthTokenRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.RefreshUserTokenInput{
		RefreshToken: req.RefreshToken,
	}
	auth, err := h.user.RefreshUserToken(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     パスワード変更
// @Description 現在のパスワードを使用して新しいパスワードに変更します。
// @Tags        Auth
// @Router      /auth/password [patch]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpdateAuthPasswordRequest true "パスワード変更"
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

	in := &user.UpdateMemberPasswordInput{
		AccessToken:          token,
		OldPassword:          req.OldPassword,
		NewPassword:          req.NewPassword,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.UpdateMemberPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     パスワード再設定リクエスト
// @Description パスワード再設定のための検証コードをメールで送信します。
// @Tags        Auth
// @Router      /auth/forgot-password [post]
// @Accept      json
// @Param       request body types.ForgotAuthPasswordRequest true "パスワード再設定リクエスト"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ForgotAuthPassword(ctx *gin.Context) {
	req := &types.ForgotAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.ForgotMemberPasswordInput{
		Email: req.Email,
	}
	if err := h.user.ForgotMemberPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     パスワード再設定実行
// @Description 検証コードを使用してパスワードを再設定します。
// @Tags        Auth
// @Router      /auth/forgot-password/verified [post]
// @Accept      json
// @Param       request body types.ResetAuthPasswordRequest true "パスワード再設定実行"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ResetAuthPassword(ctx *gin.Context) {
	req := &types.ResetAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.VerifyMemberPasswordInput{
		Email:                req.Email,
		VerifyCode:           req.VerifyCode,
		NewPassword:          req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.VerifyMemberPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     Google認証URL取得
// @Description Google OAuth認証のための認証URLを取得します。
// @Tags        Auth
// @Router      /auth/google [get]
// @Security    cookieauth
// @Param       state query string false "ステート"
// @Param       redirectUri query string false "リダイレクトURI"
// @Produce     json
// @Success     200 {object} types.AuthGoogleAccountResponse
func (h *handler) AuthGoogleAccount(ctx *gin.Context) {
	in := &user.AuthMemberWithGoogleInput{
		AuthMemberDetailWithOAuth: user.AuthMemberDetailWithOAuth{
			SessionID:   h.getSessionID(ctx),
			State:       util.GetQuery(ctx, "state", ""),
			RedirectURI: util.GetQuery(ctx, "redirectUri", ""),
		},
	}
	authURL, err := h.user.AuthMemberWithGoogle(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthGoogleAccountResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     LINE認証URL取得
// @Description LINE OAuth認証のための認証URLを取得します。
// @Tags        Auth
// @Router      /auth/line [get]
// @Security    cookieauth
// @Param       state query string false "ステート"
// @Param       redirectUri query string false "リダイレクトURI"
// @Produce     json
// @Success     200 {object} types.AuthLINEAccountResponse
func (h *handler) AuthLINEAccount(ctx *gin.Context) {
	in := &user.AuthMemberWithLINEInput{
		AuthMemberDetailWithOAuth: user.AuthMemberDetailWithOAuth{
			SessionID:   h.getSessionID(ctx),
			State:       util.GetQuery(ctx, "state", ""),
			RedirectURI: util.GetQuery(ctx, "redirectUri", ""),
		},
	}
	authURL, err := h.user.AuthMemberWithLINE(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthLINEAccountResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}
