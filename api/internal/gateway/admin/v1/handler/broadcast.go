package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

func (h *handler) broadcastRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules/:scheduleId/broadcasts", h.authentication, h.filterAccessSchedule)

	r.GET("", h.GetBroadcast)
	r.POST("", h.StartBroadcast)
	r.DELETE("", h.StopBroadcast)
	r.POST("/archive-video", h.UploadBroadcastArchive)
	r.POST("/static-image", h.ActivateBroadcastStaticImage)
	r.DELETE("/static-image", h.DeactivateBroadcastStaticImage)
	r.POST("/rtmp", h.ActivateBroadcastRTMP)
	r.POST("/mp4", h.ActivateBroadcastMP4)
}

func (h *handler) GetBroadcast(ctx *gin.Context) {
	in := &media.GetBroadcastByScheduleIDInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	broadcast, err := h.media.GetBroadcastByScheduleID(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.BroadcastResponse{
		Broadcast: service.NewBroadcast(broadcast).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) StartBroadcast(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) StopBroadcast(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) UploadBroadcastArchive(ctx *gin.Context) {
	file, header, err := h.parseFile(ctx, "video")
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	in := &media.UpdateBroadcastArchiveInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		File:       file,
		Header:     header,
	}
	if err := h.media.UpdateBroadcastArchive(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ActivateBroadcastRTMP(ctx *gin.Context) {
	in := &media.ActivateBroadcastRTMPInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.ActivateBroadcastRTMP(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ActivateBroadcastMP4(ctx *gin.Context) {
	file, header, err := h.parseFile(ctx, "video")
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	in := &media.ActivateBroadcastMP4Input{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		File:       file,
		Header:     header,
	}
	if err := h.media.ActivateBroadcastMP4(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ActivateBroadcastStaticImage(ctx *gin.Context) {
	in := &media.ActivateBroadcastStaticImageInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.ActivateBroadcastStaticImage(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeactivateBroadcastStaticImage(ctx *gin.Context) {
	in := &media.DeactivateBroadcastStaticImageInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.DeactivateBroadcastStaticImage(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}
