package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

// @tag.name        GuestLiveComment
// @tag.description ゲストライブコメント関連
func (h *handler) guestLiveCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/guests/schedules/:scheduleId/comments")

	r.POST("", h.createBroadcastViewerLog, h.CreateGuestLiveComment)
}

// @Summary     ゲストライブコメント作成
// @Description ゲストユーザーとしてライブ配信にコメントを投稿します。
// @Tags        GuestLiveComment
// @Router      /guests/schedules/{scheduleId}/comments [post]
// @Param       scheduleId path string true "スケジュールID"
// @Accept      json
// @Param       request body types.CreateGuestLiveCommentRequest true "ゲストライブコメント作成"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateGuestLiveComment(ctx *gin.Context) {
	req := &types.CreateGuestLiveCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.CreateBroadcastGuestCommentInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		Content:    req.Comment,
	}
	if _, err := h.media.CreateBroadcastGuestComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
