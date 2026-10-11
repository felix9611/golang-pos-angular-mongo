package controllers

import (
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Return Invoice record
// @Description  Create one Return Invoice record
// @Tags         Return Invoice
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateReturnInvoiceDto  true  "Code Type Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /return-invoice/create [post]
func CreateReturnInvoice(c *gin.Context) {
	var req dto.CreateReturnInvoiceDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateReturnInvoice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one Return Invoicee record by id
// @Description  Get one Return Invoice record by id
// @Tags         Return Invoice
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "data ID"
// @Success      200      {object}  dto.ReturnInvoiceDetailResponse
// @Router       /return-invoice/one/{id} [get]
func GetReturnInvoiceById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetReturnInvoiceById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      List Return Invoice records
// @Description  List Return Invoice records
// @Tags         Return Invoice
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateReturnInvoiceDto  true  "Code Type Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /return-invoice/list [post]
func ListReturnInvoice(c *gin.Context) {
	var req dto.InvoiceListRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListReturnInvoice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterReturnInvoiceRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	returnInvoiceRoutes := rg.Group("/return-invoice", handle.MiddlewareFunc())
	{
		returnInvoiceRoutes.POST("/create", CreateReturnInvoice)
		returnInvoiceRoutes.GET("/one/:id", GetReturnInvoiceById)
		returnInvoiceRoutes.POST("/list", ListReturnInvoice)
	}
}
