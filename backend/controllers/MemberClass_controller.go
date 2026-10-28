package controllers

import (
	"golang-pos-mongo-backend/backend/auth"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"golang-pos-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Member Class record
// @Description  Create one Member Class record
// @Tags         Member Class
// @Accept       json
// @Produce      json
// @Param        request  body      models.MemberClass  true  "Member Class Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /member/member-class/create [post]
func CreateMemberClass(c *gin.Context) {
	value, _ := c.Get(identityKey)
	user, ok := value.(*auth.AuthUser)
	if !ok || user == nil {
		c.JSON(401, gin.H{"error": "Invalid user"})
		return
	}

	if !user.SuperAdmin {
		c.JSON(403, gin.H{
			"error": "Access denied",
		})
	} else {
		var MemberClass models.MemberClass
		if err := c.ShouldBindJSON(&MemberClass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := services.CreateMemberClass(&MemberClass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// @Summary      Get one Member Class record by id
// @Description  Get one Member Class record by id
// @Tags         Member Class
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Member Class ID"
// @Success      200      {object}  models.MemberClass
// @Router       /member/member-class/one/{id} [get]
func GetCreateMemberClassById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneMemberClassById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      List Member Class
// @Description  List Member Class
// @Tags         Member Class
// @Accept       json
// @Produce      json
// @Param        request  body      dto.MemberClassListDto  true  "List Member Class Request Body"
// @Success      200      {object}  dto.MemberClassList
// @Router       /member/member-class/list [post]
func ListMemberClass(c *gin.Context) {
	var memberClassPageDto dto.MemberClassListDto
	if err := c.ShouldBindJSON(&memberClassPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	memberClass, err := services.ListMemberClass(&memberClassPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list asset types"})
		return
	}
	c.JSON(http.StatusOK, memberClass)
}

// @Summary      Get all records
// @Description  Get all records
// @Tags         Member Class
// @Accept       json
// @Produce      json
// @Success      200      {object}  []models.MemberClass
// @Router       /member/member-class/all [get]
func GetMemberClasses(c *gin.Context) {
	productTypes, err := services.ListAllMemberClass()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get asset types"})
		return
	}
	c.JSON(http.StatusOK, productTypes)
}

// @Summary      Update one Member Class record
// @Description  Update one Member Class record
// @Tags         Member Class
// @Accept       json
// @Produce      json
// @Param        request  body      models.MemberClass  true  "Member Class Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /member/member-class/update [post]
func UpdateMemberClass(c *gin.Context) {
	value, _ := c.Get(identityKey)
	user, ok := value.(*auth.AuthUser)
	if !ok || user == nil {
		c.JSON(401, gin.H{"error": "Invalid user"})
		return
	}

	if !user.SuperAdmin {
		c.JSON(403, gin.H{
			"error": "Access denied",
		})
	} else {
		var memberClass models.MemberClass
		if err := c.ShouldBindJSON(&memberClass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := services.UpdateMemberClass(&memberClass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// @Summary      Void one Product Types record by id
// @Description  Void one Product Types record by id
// @Tags         Product Types
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product Types ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /member/member-class/void/{id} [delete]
func InactiveMemberClass(c *gin.Context) {
	value, _ := c.Get(identityKey)
	user, ok := value.(*auth.AuthUser)
	if !ok || user == nil {
		c.JSON(401, gin.H{"error": "Invalid user"})
		return
	}

	if !user.SuperAdmin {
		c.JSON(403, gin.H{
			"error": "Access denied",
		})
	} else {
		id := c.Param("id")
		result, err := services.VoidOneMemberClass(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Asset Type inactivated successfully", "data": result})
	}
}

func RegisterMemberClassRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	MemberClassGroup := rg.Group("/member/member-class", handle.MiddlewareFunc())
	{
		MemberClassGroup.POST("/list", ListMemberClass)
		MemberClassGroup.POST("/create", CreateMemberClass)
		MemberClassGroup.GET("/one/:id", GetCreateMemberClassById)
		MemberClassGroup.GET("/all", GetMemberClasses)
		MemberClassGroup.POST("/update", UpdateMemberClass)
		MemberClassGroup.DELETE("/void/:id", InactiveMemberClass)
	}
}
