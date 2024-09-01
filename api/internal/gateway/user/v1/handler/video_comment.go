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

func (h *handler) videoCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules/:scheduleId/comments")

	r.GET("", h.createVideoViewerLog, h.ListVideoComments)
	r.POST("", h.authentication, h.createVideoViewerLog, h.CreateVideoComment)
}

func (h *handler) ListVideoComments(ctx *gin.Context) {
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

func (h *handler) CreateVideoComment(ctx *gin.Context) {
	req := &request.CreateVideoCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.CreateBroadcastCommentInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		UserID:     h.getUserID(ctx),
		Content:    req.Comment,
	}
	if _, err := h.media.CreateBroadcastComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
