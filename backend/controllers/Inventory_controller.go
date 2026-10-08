package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Stock In
// @Description  Stock In
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request  body      dto.StockInOutProductLocationDto  true  "Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /inventory/stock-in [post]
func StockIn(c *gin.Context) {
	var data dto.StockInOutProductLocationDto
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.StockIn(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Stock Move
// @Description  Stock Move
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request  body      dto.StockMoveProductLocationDto  true  "Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /inventory/stock-move [post]
func StockMove(c *gin.Context) {
	var data dto.StockMoveProductLocationDto
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.StockMove(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Stock Out
// @Description  Stock Out
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request  body      dto.StockInOutProductLocationDto  true  "Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /inventory/stock-out [post]
func StockOut(c *gin.Context) {
	var data dto.StockInOutProductLocationDto
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.StockOut(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Product Locations
// @Description  List Product Locations
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListProductLocationtRequestDto  true  "List Location Request Body"
// @Success      200      {object}  dto.ListProductLocationts
// @Router       /inventory/product-location-list [post]
func ListProductLocations(c *gin.Context) {
	var req dto.ListProductLocationtRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ProductLocationList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterInventoryRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	productLocation := rg.Group("/inventory", handle.MiddlewareFunc())
	{
		productLocation.POST("/stock-in", StockIn)
		productLocation.POST("/stock-move", StockMove)
		productLocation.POST("/stock-out", StockOut)
		productLocation.POST("/product-location-list", ListProductLocations)
	}
}
