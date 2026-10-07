package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create Asset Item record
// @Description  Create Asset Item record
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        request  body      []models.AssetLists true  "Asset List Request Body"
// @Success      200      {object}  []models.AssetLists
// @Router       /asset/asset-list/create [post]
func CreateAssetList(c *gin.Context) {
	var assetItem models.AssetLists
	if err := c.ShouldBindJSON(&assetItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to create the asset item
	createdAsset, err := services.CreateAssetItem(&assetItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdAsset)
}

// @Summary      List Asset Items With Filter
// @Description  List Asset Items With Filter
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        request  body      []dto.ListAssetReqDto true  "Asset List Request Body"
// @Success      200      {object}  []dto.AssetListsPureDetails
// @Router       /asset/asset-list/filter/list [post]
func ListAssetItemsWithFilter(c *gin.Context) {
	var req dto.ListAssetReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to list the asset items
	assetItems, err := services.ListAsseetItemsWithFilter(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assetItems)
}

// @Summary      List Asset Items
// @Description  List Asset Items
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        request  body      []dto.ListAssetReqDto true  "Asset List Request Body"
// @Success      200      {object}  []dto.AssetListList
// @Router       /asset/asset-list/list [post]
func ListAssetItems(c *gin.Context) {
	var req dto.ListAssetReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to list the asset items
	assetItems, err := services.ListAssetItems(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assetItems)
}

// @Summary      List All Asset Items
// @Description  List All Asset Items
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Success      200      {object}  []models.AssetLists
// @Router       /asset/asset-list/list-all [get]
func ListAllAssetItems(c *gin.Context) {

	// Call the service to list the asset items
	assetItems, err := services.ListAllAssetItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assetItems)
}

// @Summary      Get one asset record by id
// @Description  Get one asset record by id
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "asset list ID"
// @Success      200      {object}  models.AssetLists
// @Router       /base/asset-list/one/{id} [get]
func GetOneAssetItem(c *gin.Context) {
	id := c.Param("id")
	assetItem, err := services.GetOneAssetItemByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assetItem)
}

// @Summary      Get one asset record by code
// @Description  Get one asset record by code
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "asset code"
// @Success      200      {object}  models.AssetLists
// @Router       /base/asset-list/code/{code} [get]
func GetOneAssetItemByAssetCode(c *gin.Context) {
	id := c.Param("code")
	assetItem, err := services.GetOneAssetItemByAssetCode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assetItem)
}

// @Summary      Update Asset Item record
// @Description  Update Asset Item record
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        request  body      models.AssetLists true  "Asset List Request Body"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /asset/asset-list/update [post]
func UpdateAssetItem(c *gin.Context) {
	var updateData models.AssetLists
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to update the asset item
	updatedAsset, err := services.UpdateAssetItem(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAsset)
}

// @Summary      Query Asset Data by Data Type
// @Description  Query Asset Data by Data Type
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        request  body      dto.DashboardReqDto true  "Dashboard Request Body"
// @Success      200      {object}	[]string
// @Router       /asset/asset-list/chart-query-date [post]
func QueryAssetDataByDataType(c *gin.Context) {
	var req dto.DashboardReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to query asset data
	assetData, err := services.QueryMakerForData(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assetData)
}

// @Summary      Get one asset file records by id
// @Description  Get one asset file records by id
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "asset list ID"
// @Success      200  {object}  []models.AssetListFiles
// @Router       /base/asset-list/load-file/{id} [get]
func LoadFilesByAssetId(c *gin.Context) {
	id := c.Param("id")
	assetFiles, err := services.GetListAssetFiles(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assetFiles)
}

// @Summary      Remove Asset File record
// @Description  Remove Asset File record
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "asset file ID"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /asset/asset-list/file-remove/{id} [delete]
func RemoveAssetFile(c *gin.Context) {
	id := c.Param("id")
	res, err := services.DeleteAssetFile(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary      Batch upload Asset item record
// @Description  Batch upload Asset item record
// @Tags         AssetList
// @Accept       json
// @Produce      json
// @Param        request  body      []models.AssetTypes true  "Asset Type Request Body"
// @Success      200      string  "Batch upload Asset Type record"
// @Router       /asset/asset-list/batch-create [post]
func BatchUploadAssetList(c *gin.Context) {
	var codeTypes []dto.AssetListsPureDetails
	if err := c.ShouldBindJSON(&codeTypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.BatchCreateAssetItems(codeTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterAssetListRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	assetListGroup := rg.Group("/asset/asset-list", handle.MiddlewareFunc())
	{
		assetListGroup.POST("/create", CreateAssetList)
		assetListGroup.POST("/list", ListAssetItems)
		assetListGroup.GET("/list-all", ListAllAssetItems)
		assetListGroup.GET("/one/:id", GetOneAssetItem)
		assetListGroup.POST("/update", UpdateAssetItem)
		assetListGroup.GET("/load-file/:id", LoadFilesByAssetId)
		assetListGroup.GET("/code/:code", GetOneAssetItemByAssetCode)
		assetListGroup.DELETE("/file-remove/:id", RemoveAssetFile)
		assetListGroup.POST("/chart-query-date", QueryAssetDataByDataType)
		assetListGroup.POST("/filter/list", ListAssetItemsWithFilter)
		assetListGroup.POST("/batch-create", BatchUploadAssetList)
	}
}
