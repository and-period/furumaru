package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

// @tag.name        GuestVideoComment
// @tag.description ゲスト動画コメント関連
func (h *handler) guestVideoCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/guests/videos/:videoId/comments")

	r.POST("", h.createVideoViewerLog, h.CreateGuestVideoComment)
}

// @Summary     ゲスト動画コメント作成
// @Description ゲストユーザーとして動画にコメントを投稿します。
// @Tags        GuestVideoComment
// @Router      /guests/videos/{videoId}/comments [post]
// @Param       videoId path string true "動町ID"
// @Accept      json
// @Param       request body types.CreateGuestVideoCommentRequest true "ゲスト動画コメント作成"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateGuestVideoComment(ctx *gin.Context) {
	req := &types.CreateGuestVideoCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.CreateVideoGuestCommentInput{
		VideoID: util.GetParam(ctx, "videoId"),
		Content: req.Comment,
	}
	if _, err := h.media.CreateVideoGuestComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
