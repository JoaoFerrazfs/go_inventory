package webRoutes

import (
	baseContainer "go_inventory/Container"
	webControllers "go_inventory/Front/Application/WebControllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterWebRoutes(router *gin.Engine, dbInstance *gorm.DB) {
	container := baseContainer.BuildContainer(dbInstance)

	web := router.Group("/")

	container.Invoke(func(webcontrollers *webControllers.GeneralController) {
		webcontrollers.Register(web.Group(""))
	})
}
