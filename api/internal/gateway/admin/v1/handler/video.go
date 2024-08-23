package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) videoRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/videos", h.authentication)

	r.GET("", h.ListVideos)
	r.POST("", h.CreateVideo)
	r.GET("/:videoId", h.GetVideo)
	r.PATCH("/:videoId", h.UpdateVideo)
	r.DELETE("/:videoId", h.DeleteVideo)
}

func (h *handler) ListVideos(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.VideosResponse{
		Videos:       []*response.Video{},
		Coordinators: []*response.Coordinator{},
		Products:     []*response.Product{},
		Experiences:  []*response.Experience{},
		Total:        0,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetVideo(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.VideoResponse{
		Video:       &response.Video{},
		Coordinator: &response.Coordinator{},
		Products:    []*response.Product{},
		Experiences: []*response.Experience{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateVideo(ctx *gin.Context) {
	req := &request.CreateVideoRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.VideoResponse{
		Video:       &response.Video{},
		Coordinator: &response.Coordinator{},
		Products:    []*response.Product{},
		Experiences: []*response.Experience{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateVideo(ctx *gin.Context) {
	req := &request.UpdateVideoRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteVideo(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}
