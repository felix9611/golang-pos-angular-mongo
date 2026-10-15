package controllers

import (
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"golang-pos-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Product Type record
// @Description  Create one Product Type record
// @Tags         Product Type
// @Accept       json
// @Produce      json
// @Param        request  body      models.ProductTypes  true  "Product Type Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /product/product-type/create [post]
func CreateProductType(c *gin.Context) {
	var productType models.ProductTypes
	if err := c.ShouldBindJSON(&productType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateProductType(&productType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one Product Type record by id
// @Description  Get one Product Type record by id
// @Tags         Product Type
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product Type ID"
// @Success      200      {object}  models.ProductTypes
// @Router       /product/product-type/one/{id} [get]
func GetProductTypeById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneProductType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Update one Product Types record
// @Description  Update one Product Types record
// @Tags         Product Types
// @Accept       json
// @Produce      json
// @Param        request  body      models.ProductTypes  true  "Product Types Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /product/product-type/update [post]
func UpdateProductType(c *gin.Context) {
	var productType models.ProductTypes
	if err := c.ShouldBindJSON(&productType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateProductType(&productType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get all records
// @Description  Get all records
// @Tags         Product Type
// @Accept       json
// @Produce      json
// @Success      200      {object}  []models.ProductTypes
// @Router       /product/product-type/all [get]
func GetProductTypes(c *gin.Context) {
	productTypes, err := services.ListAlProductType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get asset types"})
		return
	}
	c.JSON(http.StatusOK, productTypes)
}

// @Summary      List Product Types With Filter
// @Description  List Product Types With Filter
// @Tags         Product Type
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ProductTypeListDto  true  "List Product Type Request Body"
// @Success      200      {object}  []models.ProductTypes
// @Router       /product/product-type/list [post]
func ListProductTypeWithFilter(c *gin.Context) {
	var productTypePageDto dto.ProductTypeListDto
	if err := c.ShouldBindJSON(&productTypePageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetTypes, err := services.ListProductTypeWithFilter(&productTypePageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list asset types"})
		return
	}
	c.JSON(http.StatusOK, assetTypes)
}

// @Summary      List Product Types
// @Description  List Product Types
// @Tags         Product Type
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ProductTypeListDto  true  "List Product Type Request Body"
// @Success      200      {object}  dto.ProductTypeList
// @Router       /product/product-type/list [post]
func ListProductType(c *gin.Context) {
	var productTypePageDto dto.ProductTypeListDto
	if err := c.ShouldBindJSON(&productTypePageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetTypes, err := services.ListProductType(&productTypePageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list asset types"})
		return
	}
	c.JSON(http.StatusOK, assetTypes)
}

// @Summary      Void one Product Types record by id
// @Description  Void one Product Types record by id
// @Tags         Product Types
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product Types ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /product/product-type/void/{id} [delete]
func InactiveProductType(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidOneProductType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset Type inactivated successfully", "data": result})
}

// @Summary      Batch upload Product Type record
// @Description  Batch upload Product Type record
// @Tags         Product Type
// @Accept       json
// @Produce      json
// @Param        request  body      []models.ProductTypes true  "Product Type Request Body"
// @Success      200      string  "Batch upload Product Type record"
// @Router       /product/product-type/batch-create [post]
func BatchUploadProductType(c *gin.Context) {
	var productTypes []models.ProductTypes
	if err := c.ShouldBindJSON(&productTypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.BatchCreateProductTypes(productTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterProductTypeRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	productType := rg.Group("/product/product-type", handle.MiddlewareFunc())
	{
		productType.POST("/create", CreateProductType)
		productType.GET("/one/:id", GetProductTypeById)
		productType.POST("/list", ListProductType)
		productType.POST("/update", UpdateProductType)
		productType.GET("/all", GetProductTypes)
		productType.DELETE("/void/:id", InactiveProductType)
		productType.POST("/filter/list", ListProductTypeWithFilter)
		productType.POST("/batch-create", BatchUploadProductType)

	}
}
