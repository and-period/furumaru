package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
)

// @tag.name        VideoComment
// @tag.description 動画コメント関連
func (h *handler) videoCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/videos/:videoId/comments", h.authentication)

	r.GET("", h.ListVideoComments)
	r.PATCH("/:commentId", h.filterAccessVideoComment, h.UpdateVideoComment)
}

func (h *handler) filterAccessVideoComment(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			video, err := h.getVideo(ctx, util.GetParam(ctx, "videoId"))
			if err != nil {
				return false, err
			}
			return video.CoordinatorID == getAdminID(ctx), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     動画コメント一覧取得
// @Description 指定された動画のコメント一覧を取得します。ページネーションと期間フィルタリングに対応しています。
// @Tags        VideoComment
// @Router      /v1/videos/{videoId}/comments [get]
// @Security    bearerauth
// @Param       videoId path string true "動画ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       start query integer false "検索開始日時（unixtime）" example("1640962800")
// @Param       end query integer false "検索終了日時（unixtime）" example("1640962800")
// @Param       next query string false "次ページトークン" example("token123")
// @Produce     json
// @Success     200 {object} types.VideoCommentsResponse
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
		WithDisabled: true,
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
		res := &types.VideoCommentsResponse{
			Comments: []*types.VideoComment{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	users, err := h.multiGetUsers(ctx, comments.UserIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.VideoCommentsResponse{
		Comments:  service.NewVideoComments(comments, users.Map()).Response(),
		NextToken: token,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     動画コメント更新
// @Description 動画コメントの状態（無効/有効）を更新します。
// @Tags        VideoComment
// @Router      /v1/videos/{videoId}/comments/{commentId} [patch]
// @Security    bearerauth
// @Param       videoId path string true "動画ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       commentId path string true "コメントID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateVideoCommentRequest true "コメント情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "コメントの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "コメントが存在しない"
func (h *handler) UpdateVideoComment(ctx *gin.Context) {
	req := &types.UpdateVideoCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.UpdateVideoCommentInput{
		CommentID: util.GetParam(ctx, "commentId"),
		Disabled:  req.Disabled,
	}
	if err := h.media.UpdateVideoComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
