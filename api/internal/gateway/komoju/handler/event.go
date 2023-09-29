package handler

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *handler) Event(ctx *gin.Context) {
	req, err := httputil.DumpRequest(ctx.Request, true)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	res, err := httputil.DumpResponse(ctx.Request.Response, true)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	h.logger.Debug("Received Event", zap.Any("request", req), zap.Any("response", res))
	ctx.Status(http.StatusNoContent)
}
