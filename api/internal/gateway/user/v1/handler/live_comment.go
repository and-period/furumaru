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

func (h *handler) liveCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules/:scheduleId/comments")

	r.GET("", h.createBroadcastViewerLog, h.ListLiveComments)
	r.POST("", h.authentication, h.createBroadcastViewerLog, h.CreateLiveComment)
}

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
	nextToken := util.GetQuery(ctx, "nextToken", "")

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

func (h *handler) CreateLiveComment(ctx *gin.Context) {
	req := &request.CreateLiveCommentRequest{}
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
