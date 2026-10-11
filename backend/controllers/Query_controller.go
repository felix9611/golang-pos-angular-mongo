package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Get query for invoice item by data type
// @Description  Get query for invoice item by data type
// @Tags         Query
// @Accept       json
// @Produce      json
// @Param        request  body      dto.InvoiceDashboardReqDto  true  "Request Body"
// @Success      200      {object}
// @Router      /query/invoice [post]
func IvoivceItemQuery(ctx *gin.Context) {
	var req dto.InvoiceDashboardReqDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.InvoiceGetByDataType(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// @Summary      Get query for product location by data type
// @Description  Get query for product location by data type
// @Tags         Query
// @Accept       json
// @Produce      json
// @Param        request  body      dto.InvoiceDashboardReqDto  true  "Request Body"
// @Success      200      {object}
// @Router      /query/product-location [post]
func ProductLocationQuery(ctx *gin.Context) {
	var req dto.InvoiceDashboardReqDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ProductLocationGetByDataType(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// @Summary      Get query for inventory by data type
// @Description  Get query for inventory by data type
// @Tags         Query
// @Accept       json
// @Produce      json
// @Param        request  body      dto.InvoiceDashboardReqDto  true  "Request Body"
// @Success      200      {object}
// @Router      /query/inventory_records [post]
func InvRecordQuery(ctx *gin.Context) {
	var req dto.InvoiceDashboardReqDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.GetInvRecordByDataType(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func RegisterQueryRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	QueryGroup := rg.Group("/query", handle.MiddlewareFunc())
	{
		QueryGroup.POST("/invoice", IvoivceItemQuery)
		QueryGroup.POST("/product-location", ProductLocationQuery)
		QueryGroup.POST("/inventory_records", InvRecordQuery)
	}
}
