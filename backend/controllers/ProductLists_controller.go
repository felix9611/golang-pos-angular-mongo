package controllers

import (
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"golang-pos-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Product List record
// @Description  Create one Product List record
// @Tags         Product List
// @Accept       json
// @Produce      json
// @Param        request  body      models.ProductLists  true  "Product List Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /product/product-list/create [post]
func CreateProductList(c *gin.Context) {
	var productList models.ProductLists
	if err := c.ShouldBindJSON(&productList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateProductItem(&productList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one product record by id
// @Description  Get one product record by id
// @Tags         Product List
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "product list ID"
// @Success      200      {object}  models.ProductLists
// @Router       /product/product-list/{id} [get]
func GetOneProductItemByID(c *gin.Context) {
	id := c.Param("id")
	assetItem, err := services.GetOneProductItemByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assetItem)
}

// @Summary      Get one product record by code
// @Description  Get one product record by code
// @Tags         ProductList
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "product code"
// @Success      200      {object}  models.ProductLists
// @Router       /product/product-list/{code} [get]
func GetOneProductItemByCode(c *gin.Context) {
	id := c.Param("code")
	productItem, err := services.GetOneProductItemByCode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, productItem)
}

// @Summary      Update Product Item record
// @Description  Update Product Item record
// @Tags         Product List
// @Accept       json
// @Produce      json
// @Param        request  body      models.ProductLists true  "Asset List Request Body"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /product/product-list/update [post]
func UpdateProductItem(c *gin.Context) {
	var updateData models.ProductLists
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to update the asset item
	updatedAsset, err := services.UpdateProductItem(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAsset)
}

// @Summary      Get one product file records by id
// @Description  Get one product file records by id
// @Tags         Product List
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "product list ID"
// @Success      200  {object}  []models.ProductListFiles
// @Router       /product/load-file/{id} [get]
func LoadFilesByProductId(c *gin.Context) {
	id := c.Param("id")
	assetFiles, err := services.GetListProductFiles(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assetFiles)
}

// @Summary      Remove Product File record
// @Description  Remove Product File record
// @Tags         Product List
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "asset file ID"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /product/product-list/file-remove/{id} [delete]
func RemoveProductFile(c *gin.Context) {
	id := c.Param("id")
	res, err := services.DeleteProductFile(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary      List Product Items
// @Description  List Product Items
// @Tags         Product List
// @Accept       json
// @Produce      json
// @Param        request  body      []dto.ListProductReqDto true  "Request Body"
// @Success      200      {object}  []dto.ProductListList
// @Router       /product/product-list/list [post]
func ListProductItems(c *gin.Context) {
	var req dto.ListProductReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to list the asset items
	productItems, err := services.ListProductItems(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productItems)
}

// @Summary      List Product Items with filters
// @Description  List Product Items with filters
// @Tags         Product List
// @Accept       json
// @Produce      json
// @Param        request  body      []dto.ListProductReqDto true  "Request Body"
// @Success      200      {object}  []dto.ListProductDetailPure
// @Router       /product/product-list/filter/list [post]
func ListProductItemsWithFilter(c *gin.Context) {
	var req dto.ListProductReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to list the asset items
	productItems, err := services.ListProductItemsWithFilters(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productItems)
}

// @Summary      Void one Product Items record by id
// @Description  Void one Product Items record by id
// @Tags         Product Items
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product Items ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /product/product-list/void/{id} [delete]
func InactiveProductItems(c *gin.Context) {
	id := c.Param("id")
	result, err := services.InactiveProductItemsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset Type inactivated successfully", "data": result})
}

func RegisterProductListRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	productList := rg.Group("/product/product-list", handle.MiddlewareFunc())
	{
		productList.POST("/create", CreateProductList)
		productList.GET("/one/:id", GetOneProductItemByID)
		productList.GET("/code/:code", GetOneProductItemByCode)
		productList.POST("/update", UpdateProductItem)
		productList.GET("/load-file/:id", LoadFilesByProductId)
		productList.DELETE("/file-remove/:id", RemoveProductFile)
		productList.POST("/list", ListProductItems)
		productList.DELETE("/void/:id", InactiveProductItems)
		productList.POST("/filter/list", ListProductItemsWithFilter)
	}
}
