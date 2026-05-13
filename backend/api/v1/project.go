package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/internal/dto"
	"github.com/wu-clan/lykn/backend/internal/service"
	"github.com/wu-clan/lykn/backend/pkg/response"
)

func ListProjects(c *gin.Context) {
	data, err := service.ListProjects(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, data)
}

func CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, err := service.CreateProject(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func GetProject(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	project, err := service.GetProject(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, project)
}

func UpdateProject(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, err := service.UpdateProject(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func DeleteProject(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := service.DeleteProject(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func DownloadProjectPublicKey(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	project, err := service.GetProject(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.Header("Content-Type", "application/x-pem-file")
	c.Header("Content-Disposition", "attachment; filename=public.pem")
	c.String(http.StatusOK, project.PublicKey)
}
