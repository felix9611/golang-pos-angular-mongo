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

// @Summary      Create one member record
// @Description  Create one member record
// @Tags         Member List
// @Accept       json
// @Produce      json
// @Param        request  body      models.MemberLists  true  "Member List Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /member/member-list/create [post]
func CreateMember(c *gin.Context) {
	var member models.MemberLists
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateMember(&member)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one member record by id
// @Description  Get one member record by id
// @Tags         Member List
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "member ID"
// @Success      200      {object}  models.MemberLists
// @Router       /member/member-list/{id} [get]
func GetOneMemberByID(c *gin.Context) {
	id := c.Param("id")
	assetItem, err := services.GetOneMemberById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assetItem)
}

// @Summary      Update member record
// @Description  Update member record
// @Tags         Member List
// @Accept       json
// @Produce      json
// @Param        request  body      models.ProductLists true  "Member List Request Body"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /member/member-list/update [post]
func UpdateMember(c *gin.Context) {
	var updateData models.MemberLists
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to update the asset item
	updatedAsset, err := services.UpdateMember(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAsset)
}

// @Summary      Remove member record
// @Description  Remove member record
// @Tags         Member List
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "asset file ID"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /member/member-list/void/{id} [delete]
func RemoveMember(c *gin.Context) {
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
		res, err := services.VoidOneMember(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

// @Summary      List Member Items
// @Description  List Member Items
// @Tags         Member List
// @Accept       json
// @Produce      json
// @Param        request  body      []dto.MemberListsDto true  "Asset List Request Body"
// @Success      200      {object}  []dto.MemberList
// @Router       /member/member-list/list [post]
func ListMemberLists(c *gin.Context) {
	var req dto.MemberListsDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to list the asset items
	final, err := services.ListMemberLists(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, final)
}

// @Summary      List Member by name
// @Description  List Member by name
// @Tags         Member List
// @Accept       json
// @Produce      json
// @Param        request  body      dto.MemberListsByNameDto true  "Request Body"
// @Success      200      {object}  []dto.MemberList
// @Router       /member/member-list/list-member [post]
func FindMembersByName(c *gin.Context) {
	var req dto.MemberListsByNameDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to list the asset items
	final, err := services.FindMembersByName(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, final)
}

// @Summary      Remove member special day record
// @Description  Remove member special day record
// @Tags         Member List
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "special day data ID"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /member/member-list/special-day/remove/{id} [delete]
func RemoveMemberSpecialDay(c *gin.Context) {
	id := c.Param("id")
	res, err := services.VoidMemberSpecialDay(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func RegisterMemberListRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	memberList := rg.Group("/member/member-list", handle.MiddlewareFunc())
	{
		memberList.POST("/create", CreateMember)
		memberList.GET("/one/:id", GetOneMemberByID)
		memberList.POST("/update", UpdateMember)
		memberList.DELETE("/void/:id", RemoveMember)
		memberList.POST("/list", ListMemberLists)
		memberList.DELETE("/special-day/remove/:id", RemoveMemberSpecialDay)
		memberList.POST("/list-member", FindMembersByName)
	}
}
