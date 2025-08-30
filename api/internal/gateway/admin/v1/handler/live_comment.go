package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
)

// @tag.name        LiveComment
// @tag.description ライブコメント関連
func (h *handler) liveCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules/:scheduleId/comments")

	r.GET("", h.ListLiveComments)
	r.PATCH("/:commentId", h.filterAccessLiveComment, h.UpdateLiveComment)
}

func (h *handler) filterAccessLiveComment(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			schedule, err := h.getSchedule(ctx, util.GetParam(ctx, "scheduleId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, schedule.CoordinatorID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     ライブコメント一覧取得
// @Description 指定されたスケジュールのライブコメント一覧を取得します。ページネーションと期間フィルタリングに対応しています。
// @Tags        LiveComment
// @Router      /v1/schedules/{scheduleId}/comments [get]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       start query integer false "検索開始日時（unixtime）" example("1640962800")
// @Param       end query integer false "検索終了日時（unixtime）" example("1640962800")
// @Param       next query string false "次ページトークン" example("token123")
// @Produce     json
// @Success     200 {object} response.LiveCommentsResponse
func (h *handler) ListLiveComments(ctx *gin.Context) {
	const defaultLimit = 20

	schedule, err := h.getSchedule(ctx, util.GetParam(ctx, "scheduleId"))
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

	in := &media.ListBroadcastCommentsInput{
		ScheduleID:   schedule.ID,
		WithDisabled: true,
		CreatedAtGte: jst.ParseFromUnix(startAt),
		CreatedAtLt:  jst.ParseFromUnix(endAt),
		Limit:        limit,
		NextToken:    nextToken,
	}
	comments, token, err := h.media.ListBroadcastComments(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(comments) == 0 {
		res := &response.LiveCommentsResponse{
			Comments: []*response.LiveComment{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	users, err := h.multiGetUsers(ctx, comments.UserIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.LiveCommentsResponse{
		Comments:  service.NewLiveComments(comments, users.Map()).Response(),
		NextToken: token,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ライブコメント更新
// @Description ライブコメントの状態（無効/有効）を更新します。
// @Tags        LiveComment
// @Router      /v1/schedules/{scheduleId}/comments/{commentId} [patch]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       commentId path string true "コメントID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body request.UpdateLiveCommentRequest true "コメント情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "コメントの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "コメントが存在しない"
func (h *handler) UpdateLiveComment(ctx *gin.Context) {
	req := &request.UpdateLiveCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.UpdateBroadcastCommentInput{
		CommentID: util.GetParam(ctx, "commentId"),
		Disabled:  req.Disabled,
	}
	if err := h.media.UpdateBroadcastComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
