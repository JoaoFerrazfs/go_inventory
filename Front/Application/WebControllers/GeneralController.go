package webControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeneralController struct{}

type Rack struct {
	ID              string
	Nome            string
	Localizacao     string
	CapacidadeTotal int
	SlotsOcupados   int
	// NOVO CAMPO
	PercentualOcupacao float64
}

func NewGeneralController() *GeneralController {
	return &GeneralController{}
}

func (controller *GeneralController) Register(group *gin.RouterGroup) {
	group.GET("/", controller.Home)
	group.GET("/racks", controller.Racks)
}

func (controller *GeneralController) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/home", gin.H{
		"title": "Home Page",
	})
}

func (controller *GeneralController) Racks(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/racks", gin.H{
		"title": "Lista de Racks",
	})
}
