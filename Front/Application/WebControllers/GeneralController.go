package webControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeneralController struct{}

func NewGeneralController() *GeneralController {
	return &GeneralController{}
}

func (controller *GeneralController) Register(group *gin.RouterGroup) {
	group.GET("/", controller.Home)
}

func (controller *GeneralController) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"title": "Home Page",
	})
}
