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
