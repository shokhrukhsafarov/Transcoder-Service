package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

var GlobalTemp = 0

func (h *handlerV1) Ping(ctx *gin.Context) {
	result := &models.PipelineWebhookResponse{
		Success: GlobalTemp > 0,
	}
	GlobalTemp++

	ctx.JSON(http.StatusOK, result)
}
