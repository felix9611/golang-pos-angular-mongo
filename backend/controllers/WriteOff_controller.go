package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Batch Create Write Off record
// @Description  Batch Create Write Off record
// @Tags         Write Off
// @Accept       json
// @Produce      json
// @Param        request  body      []models.Department  true  "Department Request Body"
// @Success      200      []models.WriteOff
// @Router      /asset/write-off/create [post]
func BatchCreateWriteOffRecord(c *gin.Context) {
	var assetData []dto.WriteOffPureList
	if err := c.ShouldBindJSON(&assetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.BatchCreateWriteOff(assetData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Create one Write Off record
// @Description  Create one Write Off record
// @Tags         Write Off
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateWriteOffRecrod  true  "Write Off Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /asset/write-off/create [post]
func CreateWriteOffRecord(c *gin.Context) {
	var writeOffRecordDto dto.CreateWriteOffRecrod
	if err := c.ShouldBindJSON(&writeOffRecordDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateWriteOff(writeOffRecordDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Write Off Records
// @Description  List Write Off Records
// @Tags         Write Off
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListWriteOffReqDto  true  "List Write Off Request Body"
// @Success      200      {object}  dto.WriteOffList
// @Router       /asset/write-off/list [post]
func ListPageWriteOff(c *gin.Context) {
	var writeOffPageDto dto.ListWriteOffReqDto
	if err := c.ShouldBindJSON(&writeOffPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.ListPageWriteOff(writeOffPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list write off records"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      List Write Off Records With Filter
// @Description  List Write Off Records With Filter
// @Tags         Write Off
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListWriteOffReqDto  true  "List Write Off Request Body"
// @Success      200      {object}  dto.WriteOffPureList
// @Router       /asset/write-off/filter/list [post]
func ListPageWriteOffWithFilter(c *gin.Context) {
	var writeOffPageDto dto.ListWriteOffReqDto
	if err := c.ShouldBindJSON(&writeOffPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.ListPageWriteOffWithFilter(writeOffPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list write off records"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func RegisterWriteOffsRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	writeOffGroup := rg.Group("/asset/write-off", handle.MiddlewareFunc())
	{
		writeOffGroup.POST("/create", CreateWriteOffRecord)
		writeOffGroup.POST("/list", ListPageWriteOff)
		writeOffGroup.POST("/filter/list", ListPageWriteOffWithFilter)
		writeOffGroup.POST("/batch-create", BatchCreateWriteOffRecord)
	}
}
