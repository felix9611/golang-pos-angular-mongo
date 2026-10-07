package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one sys Menu
// @Description  Create one sys Menu
// @Tags         Sys Menu
// @Accept       json
// @Produce      json
// @Param        request  body      models.SysMenus  true  "Create Sys Menu Request Body"
// @Success      200      {object}  dto.SysUserCreateResponseDto
// @Router       /sys/menu/create [post]
func CreateSysMenu(c *gin.Context) {
	var data models.SysMenus
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.CreateSysMenu(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      List all sys Menu records
// @Description  List all sys Menu records
// @Tags         Sys Menu
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SysMenuList  true  "List Request Body"
// @Success      200      {object}  []models.SysMenus
// @Router       /sys/menu/list [post]
func ListAllSysMenu(c *gin.Context) {
	var query dto.SysMenuList
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListAllMenu(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      List all sys Menu main IDs
// @Description  List all sys Menu main IDs
// @Tags         Sys Menu
// @Accept       json
// @Produce      json
// @Success      200      {object}  []string
// @Router       /sys/menu/main-item [get]
func GetSysMenuMainIdList(c *gin.Context) {
	result, err := services.ListAllMainIdMenu()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu main IDs"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get one Sys Menu record by id
// @Description  Get one Sys Menu record by id
// @Tags         Sys Menu
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sys Menu ID"
// @Success      200      {object}  models.SysMenus
// @Router       /sys/menu/one/{id} [get]
func GetMenuItemById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneMenuItemById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Update one sys Menu
// @Description  Update one sys Menu
// @Tags         Sys Menu
// @Accept       json
// @Produce      json
// @Param        request  body      models.SysMenus  true  "Update Sys Menu Request Body"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /sys/menu/update [post]
func UpdateMenuItemById(c *gin.Context) {
	var data models.SysMenus
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.UpdateMenuItem(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get all sys Menu items
// @Description  Get all sys Menu items
// @Tags         Sys Menu
// @Accept       json
// @Produce      json
// @Success      200      {object}  []models.SysMenus
// @Router       /sys/menu/all-menu [get]
func GetAllMenuItems(c *gin.Context) {
	result, err := services.GetAllMenuItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get menus by IDs
// @Description  Get menus by IDs
// @Tags         Sys Menu
// @Accept       json
// @Produce      json
// @Param        request  body      dto.GetMenusByIds  true  "Get Menus By Ids Request Body"
// @Success      200      {object}  []dto.SysMenuChildrenSencond
// @Router       /sys/menu/user/tree-menu [post]
func GetMenusByIds(c *gin.Context) {
	var query dto.GetMenusByIds
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.GetMenusByIds(&query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func RegisterSysMenuRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	sysMenuRoute := rg.Group("/sys/menu", handle.MiddlewareFunc())
	{
		sysMenuRoute.POST("/create", CreateSysMenu)
		sysMenuRoute.POST("/list", ListAllSysMenu)
		sysMenuRoute.GET("/main-item", GetSysMenuMainIdList)
		sysMenuRoute.GET("/one/:id", GetMenuItemById)
		sysMenuRoute.POST("/update", UpdateMenuItemById)
		sysMenuRoute.GET("/all-menu", GetAllMenuItems)
		sysMenuRoute.POST("/user/tree-menu", GetMenusByIds)
	}
}
