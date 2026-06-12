package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/dto"
	"github.com/wu-clan/lykn/internal/service"
	"github.com/wu-clan/lykn/pkg/response"
)

func ListProjectLicenses(c *gin.Context) {
	projectID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	licenses, err := service.ListProjectLicenses(c.Request.Context(), projectID)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, licenses)
}

func IssueLicense(c *gin.Context) {
	projectID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var req dto.IssueLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidRequestBody.Error())
		return
	}
	resp, _, err := service.IssueLicense(c.Request.Context(), projectID, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, resp)
}

func GetLicense(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	license, err := service.GetLicense(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Success(c, license)
}

func DownloadLicense(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	filename, content, err := service.DownloadLicense(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.String(http.StatusOK, content)
}
