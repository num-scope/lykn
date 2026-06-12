package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/pkg/response"
)

func parseUintParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, common.ErrInvalidParam.Error())
		return 0, false
	}
	return uint(value), true
}

func writeError(c *gin.Context, err error) {
	var appErr *common.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPStatus, appErr.Message)
		return
	}
	response.Error(c, http.StatusBadRequest, err.Error())
}
