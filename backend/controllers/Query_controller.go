package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Budget record
// @Description  Create one Budget record
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        request  body      models.Budgets  true  "Budget Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
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

// @Summary      Create one Budget record
// @Description  Create one Budget record
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        request  body      models.Budgets  true  "Budget Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
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

func RegisterQueryRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	QueryGroup := rg.Group("/query", handle.MiddlewareFunc())
	{
		QueryGroup.POST("/invoice", IvoivceItemQuery)
		QueryGroup.POST("/product-location", ProductLocationQuery)
	}
}
