package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Code Type record
// @Description  Create one Code Type record
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        request  body      models.CodeTypes  true  "Code Type Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /base/code-type/create [post]
func CreateCodeType(c *gin.Context) {
	var codeType models.CodeTypes
	if err := c.ShouldBindJSON(&codeType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateCodeType(&codeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one code type record by id
// @Description  Get one code type record by id
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "code type ID"
// @Success      200      {object}  models.CodeTypes
// @Router       /base/code-type/one/{id} [get]
func GetCodeTypeById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneCodeType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Delete one code type record by id
// @Description  Delete one code type record by id
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Code Type ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /base/code-type/void/{id} [delete]
func InactiveCodeTypeByID(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidOneCodeType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CodeType inactivated successfully", "data": result})
}

// @Summary      List Code Types without pagination
// @Description  List Code Types without pagination
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CodeTypeListDto  true  "List Code Type Request Body"
// @Success      200      {object}  dto.CodeTypeList
// @Router       /base/code-type/filter/list [post]
func ListCodeTypeWithoutPagination(c *gin.Context) {
	var codeTypePageDto dto.CodeTypeListDto
	if err := c.ShouldBindJSON(&codeTypePageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CodeTypeListWithoutPagination(&codeTypePageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Code Types
// @Description  List Code Types
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CodeTypeListDto  true  "List Code Type Request Body"
// @Success      200      {object}  dto.CodeTypeList
// @Router       /base/code-type/list [post]
func ListCodeType(c *gin.Context) {
	var codeTypePageDto dto.CodeTypeListDto
	if err := c.ShouldBindJSON(&codeTypePageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CodeTypeList(&codeTypePageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Update one Code Type record
// @Description  Update one Code Type record
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        request  body      models.CodeTypes  true  "Code Type Request Body"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /base/code-type/update [post]
func UpdateCodeTypeById(c *gin.Context) {
	var codeType models.CodeTypes
	if err := c.ShouldBindJSON(&codeType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateCodeType(&codeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one code type record by type
// @Description  Get one code type record by type
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        type   path      string  true  "code type"
// @Success      200      {object}  []models.CodeTypes
// @Router       /base/code-type/get-type/{type} [get]
func ListCodeTypeByType(c *gin.Context) {
	typeString := c.Param("type")

	result, err := services.ListCodeTypeByType(typeString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Batch upload code type record
// @Description  Batch upload code type record
// @Tags         Code Type
// @Accept       json
// @Produce      json
// @Param        request  body      []models.CodeTypes true  "Code Type Request Body"
// @Success      200      string  "Batch upload Asset Type record"
// @Router       /base/code-type/batch-create [post]
func BatchUploadCodeTypes(c *gin.Context) {
	var codeTypes []models.CodeTypes
	if err := c.ShouldBindJSON(&codeTypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.BatchInsertCodeTypes(codeTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterCodeTypeRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	codeTypeGroup := rg.Group("/base/code-type", handle.MiddlewareFunc())
	{
		codeTypeGroup.POST("/create", CreateCodeType)
		codeTypeGroup.GET("/one/:id", GetCodeTypeById)
		codeTypeGroup.DELETE("/void/:id", InactiveCodeTypeByID)
		codeTypeGroup.POST("/list", ListCodeType)
		codeTypeGroup.POST("/update", UpdateCodeTypeById)
		codeTypeGroup.GET("/get-type/:type", ListCodeTypeByType)
		codeTypeGroup.POST("/batch-create", BatchUploadCodeTypes)
		codeTypeGroup.POST("/filter/list", ListCodeTypeWithoutPagination)
	}
}
