package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/internal/service"
	"github.com/wu-clan/lykn/pkg/response"
)

func GetDashboardSummary(c *gin.Context) {
	data, err := service.GetDashboardSummary(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, data)
}
