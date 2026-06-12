package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/dto"
	"github.com/wu-clan/lykn/internal/service"
	"github.com/wu-clan/lykn/pkg/response"
)

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, err := service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func GetCurrentUser(c *gin.Context) {
	data, err := service.GetCurrentUser(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, data)
}
