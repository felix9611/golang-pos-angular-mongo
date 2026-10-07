package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Batch Create Tax Information record
// @Description  Batch Create Tax Information record
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        request  body      []models.TaxInformations  true  "Tax Information Request Body"
// @Success      200      string
// @Router      /base/tax-information/batch-create [post]
func BatchCreateTaxInformation(c *gin.Context) {
	var taxInfo []models.TaxInformations
	if err := c.ShouldBindJSON(&taxInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.BatchCreateTaxInformation(taxInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Create one Tax Information record
// @Description  Create one Tax Information record
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        request  body      models.TaxInformations  true  "Tax Information Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /base/tax-information/create [post]
func CreateTaxInformation(c *gin.Context) {
	var taxInfo models.TaxInformations
	if err := c.ShouldBindJSON(&taxInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateTaxInformation(&taxInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one tax information record by id
// @Description  Get one tax information record by id
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "tax information ID"
// @Success      200      {object}  models.TaxInformations
// @Router       /base/tax-information/one/{id} [get]
func GetOneTaxInformation(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneTaxInformation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Void one tax information record by id
// @Description  Delete one tax information record by id
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tax Information ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router       /base/tax-information/void/{id} [delete]
func VoidOneTaxInformation(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidOneTaxInformation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tax Information inactivated successfully", "data": result})
}

// @Summary      List Tax Information
// @Description  List Tax Information
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        request  body      dto.TaxInformationListDTO  true  "List Tax Information Request Body"
// @Success      200      {object}  dto.TaxInformationList
// @Router       /base/tax-information/list [post]
func ListTaxInformation(c *gin.Context) {
	var taxInfoPageDto dto.TaxInformationListDTO
	if err := c.ShouldBindJSON(&taxInfoPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.TaxInformationList(&taxInfoPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Update one Tax Information record
// @Description  Update one Tax Information record
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        request  body      models.TaxInformations  true  "Tax Information Request Body"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /base/tax-information/update [post]
func UpdateTaxInformationById(c *gin.Context) {
	var taxInfo models.TaxInformations
	if err := c.ShouldBindJSON(&taxInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateTaxInformation(&taxInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Tax Information With filter
// @Description  List Tax Information With filter
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        request  body      dto.TaxInformationListDTO  true  "List Tax Information Request Body"
// @Success      200      {object}  []models.TaxInformations
// @Router       /base/tax-information/filter/list [post]
func ListTaxInformationsWithoutPagination(c *gin.Context) {
	var taxInfoPageDto dto.TaxInformationListDTO
	if err := c.ShouldBindJSON(&taxInfoPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.TaxInformationListWithoutPagination(&taxInfoPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Tax Information
// @Description  List Tax Information
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        request  body      dto.TaxInformationListDTO  true  "List Tax Information Request Body"
// @Success      200      {object}  dto.TaxInformationList
// @Router       /base/tax-information/list [post]
func ListTaxInformations(c *gin.Context) {
	var taxInfoPageDto dto.TaxInformationListDTO
	if err := c.ShouldBindJSON(&taxInfoPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.TaxInformationList(&taxInfoPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get all Tax Information record
// @Description  Get all Tax Information record
// @Tags         Tax Information
// @Produce      json
// @Success      200      {object}  []models.TaxInformations
// @Router       /base/tax-information/all [get]
func ListAllTaxInformation(c *gin.Context) {
	result, err := services.ListAllTaxInformation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func RegisterTaxInformationRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	TaxInformationRoute := rg.Group("/base/tax-information", handle.MiddlewareFunc())
	{
		TaxInformationRoute.POST("/create", CreateTaxInformation)
		TaxInformationRoute.GET("/one/:id", GetOneTaxInformation)
		TaxInformationRoute.POST("/update", UpdateTaxInformationById)
		TaxInformationRoute.DELETE("/void/:id", VoidOneTaxInformation)
		TaxInformationRoute.POST("/list", ListTaxInformations)
		TaxInformationRoute.GET("/all", ListAllTaxInformation)
		TaxInformationRoute.POST("/filter/list", ListTaxInformationsWithoutPagination)
		TaxInformationRoute.POST("/batch-create", BatchCreateTaxInformation)
	}

}
