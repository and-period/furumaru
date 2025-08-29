package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
)

// @tag.name        VideoComment
// @tag.description 動画コメント関連
func (h *handler) videoCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/videos/:videoId/comments")

	r.GET("", h.createVideoViewerLog, h.ListVideoComments)
	r.POST("", h.authentication, h.createVideoViewerLog, h.CreateVideoComment)
}

// @Summary     動画コメント一覧取得
// @Description 動画のコメント一覧を取得します。
// @Tags        VideoComment
// @Router      /videos/{videoId}/comments [get]
// @Param       videoId path string true "動町ID"
// @Param       limit query int64 false "取得件数" default(20)
// @Produce     json
// @Success     200 {object} response.VideoCommentsResponse
func (h *handler) ListVideoComments(ctx *gin.Context) {
	const defaultLimit = 20

	video, err := h.getVideo(ctx, util.GetParam(ctx, "videoId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	startAt, err := util.GetQueryInt64(ctx, "start", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	endAt, err := util.GetQueryInt64(ctx, "end", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	nextToken := util.GetQuery(ctx, "next", "")

	in := &media.ListVideoCommentsInput{
		VideoID:      video.ID,
		CreatedAtGte: jst.ParseFromUnix(startAt),
		CreatedAtLt:  jst.ParseFromUnix(endAt),
		Limit:        limit,
		NextToken:    nextToken,
	}
	comments, token, err := h.media.ListVideoComments(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(comments) == 0 {
		res := &response.VideoCommentsResponse{
			Comments: []*response.VideoComment{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	users, err := h.multiGetUsers(ctx, comments.UserIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.VideoCommentsResponse{
		Comments:  service.NewVideoComments(comments, users.Map()).Response(),
		NextToken: token,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     動画コメント作成
// @Description 動画にコメントを投稿します。
// @Tags        VideoComment
// @Router      /videos/{videoId}/comments [post]
// @Security    bearerauth
// @Param       videoId path string true "動町ID"
// @Accept      json
// @Param       request body request.CreateVideoCommentRequest true "動画コメント作成"
// @Success     204 "作成成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "オンデマンド配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "オンデマンド配信が公開されていない"
func (h *handler) CreateVideoComment(ctx *gin.Context) {
	req := &request.CreateVideoCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.CreateVideoCommentInput{
		VideoID: util.GetParam(ctx, "videoId"),
		UserID:  h.getUserID(ctx),
		Content: req.Comment,
	}
	if _, err := h.media.CreateVideoComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
