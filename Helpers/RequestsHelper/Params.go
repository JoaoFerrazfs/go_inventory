package requestsHelper

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIDParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return 0, err
	}
	return uint(idUint64), nil
}

func GetParam(c *gin.Context, paramName string) string {
	param := c.Param(paramName)

	return param
}

func GetParamAsInt(c *gin.Context, paramName string) (int, error) {
	param := c.Param(paramName)
	intParam, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return 0, err
	}

	return intParam, nil
}
