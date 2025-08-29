package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        AuthUser
// @tag.description 認証ユーザー関連
func (h *handler) authUserRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/users/me")

	r.POST("", h.CreateAuthUser)
	r.POST("/verified", h.VerifyAuthUser)
	r.POST("/google", h.CreateAuthUserWithGoogle)
	r.POST("/line", h.CreateAuthUserWithLINE)

	auth := r.Group("", h.authentication)
	auth.GET("", h.GetAuthUser)
	auth.DELETE("", h.DeleteAuthUser)
	auth.PATCH("/email", h.UpdateAuthUserEmail)
	auth.POST("/email/verified", h.VerifyAuthUserEmail)
	auth.PATCH("/username", h.UpdateAuthUserUsername)
	auth.PATCH("/account-id", h.UpdateAuthUserAccountID)
	auth.PATCH("/notification", h.UpdateAuthUserNotification)
	auth.PATCH("/thumbnail", h.UpdateAuthUserThumbnail)
}

// @Summary     認証ユーザー情報取得
// @Description ログイン中のユーザー情報を取得します。
// @Tags        AuthUser
// @Router      /users/me [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} response.AuthUserResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) GetAuthUser(ctx *gin.Context) {
	in := &user.GetUserInput{
		UserID: h.getUserID(ctx),
	}
	uuser, err := h.user.GetUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	notificationIn := &user.GetUserNotificationInput{
		UserID: uuser.ID,
	}
	notification, err := h.user.GetUserNotification(ctx, notificationIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(uuser, notification).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ユーザー登録
// @Description 新しいユーザーを登録します。
// @Tags        AuthUser
// @Router      /users/me [post]
// @Accept      json
// @Produce     json
// @Param       body body request.CreateAuthUserRequest true "ユーザー情報"
// @Success     200 {object} response.CreateAuthUserResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在するメールアドレス"
func (h *handler) CreateAuthUser(ctx *gin.Context) {
	req := &request.CreateAuthUserRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.CreateMemberInput{
		Username:             req.Username,
		AccountID:            req.AccountID,
		Lastname:             req.Lastname,
		Firstname:            req.Firstname,
		LastnameKana:         req.LastnameKana,
		FirstnameKana:        req.FirstnameKana,
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	userID, err := h.user.CreateMember(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.CreateAuthUserResponse{
		ID: userID,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ユーザー登録確認
// @Description ユーザー登録時に送信される確認コードで登録を確定します。
// @Tags        AuthUser
// @Router      /users/me/verified [post]
// @Accept      json
// @Param       body body request.VerifyAuthUserRequest true "確認コード"
// @Success     204 "確認成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) VerifyAuthUser(ctx *gin.Context) {
	req := &request.VerifyAuthUserRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.VerifyMemberInput{
		UserID:     req.ID,
		VerifyCode: req.VerifyCode,
	}
	if err := h.user.VerifyMember(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     Googleアカウントでユーザー登録
// @Description Googleアカウントを使用してユーザー登録を行います。
// @Tags        AuthUser
// @Router      /users/me/google [post]
// @Security    cookieauth
// @Accept      json
// @Produce     json
// @Param       body body request.CreateAuthUserWithGoogleRequest true "Googleアカウント情報"
// @Success     200 {object} response.AuthUserResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateAuthUserWithGoogle(ctx *gin.Context) {
	req := &request.CreateAuthUserWithGoogleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	userIn := &user.CreateMemberWithGoogleInput{
		CreateMemberDetailWithOAuth: user.CreateMemberDetailWithOAuth{
			SessionID:     h.getSessionID(ctx),
			Code:          req.Code,
			Nonce:         req.Nonce,
			RedirectURI:   req.RedirectURI,
			Username:      req.Username,
			AccountID:     req.AccountID,
			Lastname:      req.Lastname,
			Firstname:     req.Firstname,
			LastnameKana:  req.LastnameKana,
			FirstnameKana: req.FirstnameKana,
			PhoneNumber:   req.PhoneNumber,
		},
	}
	uuser, err := h.user.CreateMemberWithGoogle(ctx, userIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	notificationIn := &user.GetUserNotificationInput{
		UserID: uuser.ID,
	}
	notification, err := h.user.GetUserNotification(ctx, notificationIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(uuser, notification).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     LINEアカウントでユーザー登録
// @Description LINEアカウントを使用してユーザー登録を行います。
// @Tags        AuthUser
// @Router      /users/me/line [post]
// @Security    cookieauth
// @Accept      json
// @Produce     json
// @Param       body body request.CreateAuthUserWithLINERequest true "LINEアカウント情報"
// @Success     200 {object} response.AuthUserResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateAuthUserWithLINE(ctx *gin.Context) {
	req := &request.CreateAuthUserWithLINERequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	userIn := &user.CreateMemberWithLINEInput{
		CreateMemberDetailWithOAuth: user.CreateMemberDetailWithOAuth{
			SessionID:     h.getSessionID(ctx),
			Code:          req.Code,
			Nonce:         req.Nonce,
			RedirectURI:   req.RedirectURI,
			Username:      req.Username,
			AccountID:     req.AccountID,
			Lastname:      req.Lastname,
			Firstname:     req.Firstname,
			LastnameKana:  req.LastnameKana,
			FirstnameKana: req.FirstnameKana,
			PhoneNumber:   req.PhoneNumber,
		},
	}
	uuser, err := h.user.CreateMemberWithLINE(ctx, userIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	notificationIn := &user.GetUserNotificationInput{
		UserID: uuser.ID,
	}
	notification, err := h.user.GetUserNotification(ctx, notificationIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(uuser, notification).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     メールアドレス更新
// @Description ユーザーのメールアドレスを更新します。
// @Tags        AuthUser
// @Router      /users/me/email [patch]
// @Security    bearerauth
// @Accept      json
// @Param       body body request.UpdateAuthUserEmailRequest true "メールアドレス"
// @Success     204 "更新成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在するメールアドレス"
// @Failure     412 {object} util.ErrorResponse "変更後のメールアドレスが変更前と同じ"
func (h *handler) UpdateAuthUserEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAuthUserEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateMemberEmailInput{
		AccessToken: token,
		Email:       req.Email,
	}
	if err := h.user.UpdateMemberEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     メールアドレス更新確認
// @Description メールアドレス更新時に送信される確認コードで更新を確定します。
// @Tags        AuthUser
// @Router      /users/me/email/verified [post]
// @Security    bearerauth
// @Accept      json
// @Param       body body request.VerifyAuthUserEmailRequest true "確認コード"
// @Success     204 "確認成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) VerifyAuthUserEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.VerifyAuthUserEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.VerifyMemberEmailInput{
		AccessToken: token,
		VerifyCode:  req.VerifyCode,
	}
	if err := h.user.VerifyMemberEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ユーザー名更新
// @Description ユーザーの表示名を更新します。
// @Tags        AuthUser
// @Router      /users/me/username [patch]
// @Security    bearerauth
// @Accept      json
// @Param       body body request.UpdateAuthUserUsernameRequest true "ユーザー名"
// @Success     204 "更新成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) UpdateAuthUserUsername(ctx *gin.Context) {
	req := &request.UpdateAuthUserUsernameRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateMemberUsernameInput{
		UserID:   h.getUserID(ctx),
		Username: req.Username,
	}
	if err := h.user.UpdateMemberUsername(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     アカウントID更新
// @Description ユーザーのアカウントIDを更新します。
// @Tags        AuthUser
// @Router      /users/me/account-id [patch]
// @Security    bearerauth
// @Accept      json
// @Param       body body request.UpdateAuthUserAccountIDRequest true "アカウントID"
// @Success     204 "更新成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) UpdateAuthUserAccountID(ctx *gin.Context) {
	req := &request.UpdateAuthUserAccountIDRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateMemberAccountIDInput{
		UserID:    h.getUserID(ctx),
		AccountID: req.AccountID,
	}
	if err := h.user.UpdateMemberAccountID(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     通知設定更新
// @Description ユーザーの通知設定を更新します。
// @Tags        AuthUser
// @Router      /users/me/notification [patch]
// @Security    bearerauth
// @Accept      json
// @Param       body body request.UpdateAuthUserNotificationRequest true "通知設定"
// @Success     204 "更新成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) UpdateAuthUserNotification(ctx *gin.Context) {
	req := &request.UpdateAuthUserNotificationRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateUserNotificationInput{
		UserID:  h.getUserID(ctx),
		Enabled: req.Enabled,
	}
	if err := h.user.UpdateUserNotification(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     サムネイル更新
// @Description ユーザーのサムネイル画像を更新します。
// @Tags        AuthUser
// @Router      /users/me/thumbnail [patch]
// @Security    bearerauth
// @Accept      json
// @Param       body body request.UpdateAuthUserThumbnailRequest true "サムネイルURL"
// @Success     204 "更新成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) UpdateAuthUserThumbnail(ctx *gin.Context) {
	req := &request.UpdateAuthUserThumbnailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateMemberThumbnailURLInput{
		UserID:       h.getUserID(ctx),
		ThumbnailURL: req.ThumbnailURL,
	}
	if err := h.user.UpdateMemberThumbnailURL(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ユーザー削除
// @Description ユーザーアカウントを削除します。
// @Tags        AuthUser
// @Router      /users/me [delete]
// @Security    bearerauth
// @Success     204 "削除成功"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) DeleteAuthUser(ctx *gin.Context) {
	in := &user.DeleteUserInput{
		UserID: h.getUserID(ctx),
	}
	if err := h.user.DeleteUser(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetUsers(ctx context.Context, userIDs []string) (entity.Users, error) {
	if len(userIDs) == 0 {
		return entity.Users{}, nil
	}
	in := &user.MultiGetUsersInput{
		UserIDs: userIDs,
	}
	return h.user.MultiGetUsers(ctx, in)
}

func (h *handler) getMember(ctx context.Context, userID string) (*entity.User, error) {
	in := &user.GetUserInput{
		UserID: userID,
	}
	return h.user.GetUser(ctx, in)
}
