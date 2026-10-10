package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Invoice record
// @Description  Create one Invoice record
// @Tags         Invoice
// @Accept       json
// @Produce      json
// @Param        request  body      dto.InvoiceDto  true  "Invoice Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /invoice/create [post]
func CreateInvoice(c *gin.Context) {
	var req dto.InvoiceDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateInvoice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Invoice
// @Description  List Invoice
// @Tags         Invoice
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListInvoiceRequestDto  true  "List Code Type Request Body"
// @Success      200      {object}  dto.ListInvoices
// @Router       /invoice/list [post]
func ListInvoice(c *gin.Context) {
	var req dto.ListInvoiceRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListInvoices(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one Invoice record by id
// @Description  Get one Invoice record by id
// @Tags         Invoice
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Invoice ID"
// @Success      200      {object}  dto.InvoiceBody
// @Router       /invoice/one/{id} [get]
func GetInvoiceById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneInvoiceDetailById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get one Invoice record by number
// @Description  Get one Invoice record by number
// @Tags         Invoice
// @Accept       json
// @Produce      json
// @Param        number   path      string  true  "Invoice Number"
// @Success      200      {object}  dto.InvoiceBody
// @Router       /invoice/number/{number} [get]
func GetInvoiceByNo(c *gin.Context) {
	id := c.Param("number")
	result, err := services.GetOneInvoiceDetailByInvNo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Void one Invoice record by id
// @Description  Void one Invoice record by id
// @Tags         Invoice
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Invoicen ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router       /base/location/void/{id} [delete]
func VoidInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidOneInvoice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Void Invoice successfully", "id": id, "data": result})
}

func RegisterInvoiceRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	invoiceGroup := rg.Group("/invoice", handle.MiddlewareFunc())
	{
		invoiceGroup.POST("/create", CreateInvoice)
		invoiceGroup.POST("/list", ListInvoice)
		invoiceGroup.GET("/one/:id", GetInvoiceById)
		invoiceGroup.DELETE("/void/:id", VoidInvoiceByID)
		invoiceGroup.GET("/number/:number", GetInvoiceByNo)
	}
}
