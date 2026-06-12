package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/dto"
	"github.com/wu-clan/lykn/internal/service"
	"github.com/wu-clan/lykn/pkg/response"
)

func ListFeatures(c *gin.Context) {
	data, err := service.ListFeatures(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, data)
}

func CreateFeature(c *gin.Context) {
	var req dto.CreateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, err := service.CreateFeature(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func UpdateFeature(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, err := service.UpdateFeature(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func DeleteFeature(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := service.DeleteFeature(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
