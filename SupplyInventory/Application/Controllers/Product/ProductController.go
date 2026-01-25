package controllers

import (
	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests/Product"
	services "go_inventory/SupplyInventory/Application/Services/Product"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (controller *ProductController) Register(group *gin.RouterGroup) {
	group.POST("", controller.CreateProduct)
	group.GET("/:ean", controller.GetProductByEAN)
	group.DELETE("/:ean", controller.DeleteProduct)
}

// @Summary Create a new product
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body requests.CreateProductRequest true "Product Data (EAN must be 13 digits)"
// @Success 201 {object} apiContracts.ProductResponse
// @Failure 422 {object} map[string]string
// @Router /api/v1/products [post]
func (controller *ProductController) CreateProduct(c *gin.Context) {
	var req requests.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, requestsHelper.FormatValidationErrors(err))
		return
	}

	product, appErr := controller.productService.CreateProduct(req.EAN, req.Name)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(201, gin.H{"data": gin.H{"product": product}})
}

// @Summary Get product by EAN
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ean path string true "Product EAN (13 digits)"
// @Success 200 {object} apiContracts.ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/{ean} [get]
func (controller *ProductController) GetProductByEAN(c *gin.Context) {
	ean := c.Param("ean")

	product, appErr := controller.productService.GetProductByEAN(ean)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(200, gin.H{"data": gin.H{"product": product}})
}

// @Summary Delete product by EAN
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ean path string true "Product EAN (13 digits)"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/{ean} [delete]
func (controller *ProductController) DeleteProduct(c *gin.Context) {
	ean := c.Param("ean")

	deleted, appErr := controller.productService.DeleteProduct(ean)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	if !deleted {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(200, gin.H{"data": "Product deleted successfully"})
}
