package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      List ExcelFieldMatch records
// @Description  List ExcelFieldMatch records
// @Tags         Excel Field Match
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ExcelFieldMatchPageDto  true  "List ExcelFieldMatch Request Body"
// @Success      200      {object}  dto.ExcelFieldMatchList
// @Router       /sys/excel-field-match/list [post]
func ListPageExcelFieldMatch(c *gin.Context) {
	var param dto.ExcelFieldMatchPageDto
	result, err := services.ExcelFieldMatchListAndPage(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Create ExcelFieldMatch record
// @Description  Create ExcelFieldMatch record
// @Tags         Excel Field Match
// @Accept       json
// @Produce      json
// @Param        request  body      models.ExcelFieldMatchs true  "Excel Field Match Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router       /sys/excel-field-match/create [post]
func CreateExcelFieldMatch(c *gin.Context) {
	var param models.ExcelFieldMatchs
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.CreateExcelFieldMatch(&param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get one ExcelFieldMatch record by id
// @Description  Get one ExcelFieldMatch record by id
// @Tags         Excel Field Match
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ExcelFieldMatch ID"
// @Success      200      {object}  models.ExcelFieldMatchs
// @Router       /sys/excel-field-match/one/{id} [get]
func GetOneExcelFieldMatchById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneExcelFieldMatch(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get one ExcelFieldMatch record by code
// @Description  Get one ExcelFieldMatch record by code
// @Tags         Excel Field Match
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Code"
// @Success      200      {object}  models.ExcelFieldMatchs
// @Router       /sys/excel-field-match/one/{id} [get]
func GetOneExcelFieldMatchByCode(c *gin.Context) {
	code := c.Param("code")
	result, err := services.GetOneExcelFieldMatchByCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Void one ExcelFieldMatch record by id
// @Description  Void one ExcelFieldMatch record by id
// @Tags         Excel Field Match
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ExcelFieldMatch ID"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /sys/excel-field-match/void/{id} [delete]
func VoidOneExcelFieldMatchById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.InactiveOneExcelFieldMatch(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Update ExcelFieldMatch record
// @Description  Update ExcelFieldMatch record
// @Tags         Excel Field Match
// @Accept       json
// @Produce      json
// @Param        request  body      models.ExcelFieldMatchs true  "Excel Field Match Request Body"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /sys/excel-field-match/update [post]
func UpdateExcelFieldMatchById(c *gin.Context) {
	var param models.ExcelFieldMatchs
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.UpdateOneExcelFieldMatch(&param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func RegisterExcelFieldMatchRoutes(r *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	excelFieldMatchGroup := r.Group("/sys/excel-field-match", handle.MiddlewareFunc())
	{
		excelFieldMatchGroup.GET("/one/:id", GetOneExcelFieldMatchById)
		excelFieldMatchGroup.POST("/list", ListPageExcelFieldMatch)
		excelFieldMatchGroup.POST("/create", CreateExcelFieldMatch)
		excelFieldMatchGroup.DELETE("/void/:id", VoidOneExcelFieldMatchById)
		excelFieldMatchGroup.POST("/update", UpdateExcelFieldMatchById)
		excelFieldMatchGroup.GET("/code/:code", GetOneExcelFieldMatchByCode)
	}
}
