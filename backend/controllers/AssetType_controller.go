package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Asset Type record
// @Description  Create one Asset Type record
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Param        request  body      models.AssetTypes  true  "Asset Type Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /asset/type/create [post]
func CreateAssetType(c *gin.Context) {
	var assetType models.AssetTypes
	if err := c.ShouldBindJSON(&assetType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateAssetType(&assetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one Asset Type record by id
// @Description  Get one Asset Type record by id
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset Type ID"
// @Success      200      {object}  models.AssetTypes
// @Router       /asset/type/one/{id} [get]
func GetAssetTypeById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneAssetType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Void one Asset Type record by id
// @Description  Void one Asset Type record by id
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset Type ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /asset/type/void/{id} [delete]
func InactiveAssetType(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidOneAssetType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset Type inactivated successfully", "data": result})
}

// @Summary      Update one Asset Type record
// @Description  Update one Asset Type record
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Param        request  body      models.AssetTypes  true  "Asset Type Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /asset/type/update [post]
func UpdateAssetType(c *gin.Context) {
	var assetType models.AssetTypes
	if err := c.ShouldBindJSON(&assetType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateAssetType(&assetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one Asset Type records
// @Description  Get one Asset Type records
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Success      200      {object}  []models.AssetTypes
// @Router       /asset/type/all [get]
func GetAssetTypes(c *gin.Context) {
	assetTypes, err := services.ListAllAssetType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get asset types"})
		return
	}
	c.JSON(http.StatusOK, assetTypes)
}

// @Summary      List Asset Types
// @Description  List Asset Types
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AssetTypeListDto  true  "List Asset Types Request Body"
// @Success      200      {object}  dto.AssetTypeList
// @Router       /asset/type/list [post]
func ListPageAssetTypes(c *gin.Context) {
	var assetTypePageDto dto.AssetTypeListDto
	if err := c.ShouldBindJSON(&assetTypePageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetTypes, err := services.ListAssetType(&assetTypePageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list asset types"})
		return
	}
	c.JSON(http.StatusOK, assetTypes)

}

// @Summary      List Asset Types
// @Description  List Asset Types
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AssetTypeListDto  true  "List Asset Types Request Body"
// @Success      200      {object}  []models.AssetTypes
// @Router       /asset/type/filter/list [post]
func ListAssetTypes(c *gin.Context) {
	var assetTypePageDto dto.AssetTypeListDto
	if err := c.ShouldBindJSON(&assetTypePageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetTypes, err := services.ListAssetTypeNoPaging(&assetTypePageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list asset types"})
		return
	}
	c.JSON(http.StatusOK, assetTypes)
}

// @Summary      Batch upload Asset Type record
// @Description  Batch upload Asset Type record
// @Tags         Asset Type
// @Accept       json
// @Produce      json
// @Param        request  body      []models.AssetTypes true  "Asset Type Request Body"
// @Success      200      string  "Batch upload Asset Type record"
// @Router       /asset/type/batch-create [post]
func BatchUploadAssetTypes(c *gin.Context) {
	var codeTypes []models.AssetTypes
	if err := c.ShouldBindJSON(&codeTypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.BatchCreateAssetTypes(codeTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterAssetTypeRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	AssetTypeGroup := rg.Group("/asset/type", handle.MiddlewareFunc())
	{
		AssetTypeGroup.POST("/create", CreateAssetType)
		AssetTypeGroup.GET("/one/:id", GetAssetTypeById)
		AssetTypeGroup.DELETE("/void/:id", InactiveAssetType)
		AssetTypeGroup.POST("/update", UpdateAssetType)
		AssetTypeGroup.GET("/all", GetAssetTypes)
		AssetTypeGroup.POST("/list", ListPageAssetTypes)
		AssetTypeGroup.POST("filter/list", ListAssetTypes)
		AssetTypeGroup.POST("/batch-create", BatchUploadAssetTypes)
	}
}
