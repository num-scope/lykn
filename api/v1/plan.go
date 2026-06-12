package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/dto"
	"github.com/wu-clan/lykn/internal/service"
	"github.com/wu-clan/lykn/pkg/response"
)

func ListPlans(c *gin.Context) {
	data, err := service.ListPlans(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, data)
}

func CreatePlan(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, err := service.CreatePlan(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func UpdatePlan(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, err := service.UpdatePlan(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func DeletePlan(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := service.DeletePlan(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
