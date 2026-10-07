package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	//	"time"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Sys User
// @Description  Create one Sys User
// @Tags         Sys User
// @Accept       json
// @Produce      json
// @Param        request  body      models.SysUsers  true  "Sys User Request Body"
// @Success      200      {object}  dto.SysUserCreateResponseDto
// @Router      /sys/user/create [post]
func CreateSysUser(c *gin.Context) {
	var user models.SysUsers
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := services.CreateSysUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary      Update one Sys User
// @Description  Update one Sys User
// @Tags         Sys User
// @Accept       json
// @Produce      json
// @Param        request  body      models.SysUsers  true  "Sys User Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /sys/user/update [post]
func UpdateSysUserByID(c *gin.Context) {
	var user models.SysUsers
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateSysUserByID(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary      Get one Sys User record by id
// @Description  Get one Sys User record by id
// @Tags         Sys User
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sys User ID"
// @Success      200      {object}  models.SysUsers
// @Router       /sys/user/one/{id} [get]
func GetSysUsersById(c *gin.Context) {
	id := c.Param("id")
	user, err := services.GetOneSysUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary      Void one Sys User by id
// @Description  Void one Sys User by id
// @Tags         Sys User
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /sys/user/{id} [delete]
func InactiveUserByID(c *gin.Context) {
	id := c.Param("id")
	result, err := services.InactiveUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inactive user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary      List Sys Users
// @Description  List Sys Users
// @Tags         Sys User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SysUserPageDto  true  "List Sys Users Request Body"
// @Success      200      {object}  dto.SysUserList
// @Router       /sys/user/list [post]
func SysUserLists(c *gin.Context) {
	var pageDto dto.SysUserPageDto
	if err := c.ShouldBindJSON(&pageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.SysUserList(&pageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// @Summary      Update one Sys User Avatar
// @Description  Update one Sys User Avatar
// @Tags         Sys User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SysUserAvatarUpdateDto  true  "Sys User Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /sys/user/user-self/update-avatar [post]
func SysUserAvatarUpdate(c *gin.Context) {
	var avatarDto dto.SysUserAvatarUpdateDto
	if err := c.ShouldBindJSON(&avatarDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UserUpdateAvatar(avatarDto.Username, avatarDto.PhotoBase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary      Update one Sys User Password
// @Description  Update one Sys User Password
// @Tags         Sys User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SysUserUpdatePasswordDto  true  "Sys User Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /sys/user/user-self/update-password [post]
func SysUserPasswordUpdate(c *gin.Context) {
	var passwordDto dto.SysUserUpdatePasswordDto
	if err := c.ShouldBindJSON(&passwordDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateUserPassword(passwordDto.Username, passwordDto.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func RegisterSysUserRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	userGroup := rg.Group("/sys/user", handle.MiddlewareFunc())
	{
		userGroup.POST("/create", CreateSysUser)
		userGroup.POST("/update", UpdateSysUserByID)
		userGroup.GET("/one/:id", GetSysUsersById)
		userGroup.DELETE("/:id", InactiveUserByID)
		userGroup.POST("/list", SysUserLists)
		userGroup.POST("/user-self/update-avatar", SysUserAvatarUpdate)
		userGroup.POST("/user-self/update-password", SysUserPasswordUpdate)
	}
}
