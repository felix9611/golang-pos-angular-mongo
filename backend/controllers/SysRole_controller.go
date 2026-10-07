package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one sys role
// @Description  Create one sys role
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Param        request  body      models.SysRoles  true  "List Action Record Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router       /sys/role/create [post]
func CreateSysRoleApi(c *gin.Context) {
	var role models.SysRoles
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := services.CreateSysRole(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary      Get one sys role record by id
// @Description  Get one sys role record by id
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "sys role ID"
// @Success      200      {object}  models.SysRoles
// @Router       /sys/role/one/{id} [get]
func GetRoleById(c *gin.Context) {
	id := c.Param("id")
	role, err := services.GetOneSysRoleById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": role})
}

// @Summary      Get all sys role
// @Description  Get all sys role
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Success      200      {object}  []dto.ListAllSysRoleDto
// @Router       /sys/role/all [get]
func GetAllRoles(c *gin.Context) {
	roles, err := services.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get roles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// @Summary      Update sys role
// @Description  Update sys role
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Param        request  body      models.SysRoles  true  "Update Request Body"
// @Success      200      {object}  models.SysRoles
// @Router       /sys/role/update [post]
func UpdateSysRole(c *gin.Context) {
	// id := c.Param("id")
	var role models.SysRoles
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateRoleById(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully", "id": role.ID, "data": result})
}

// @Summary      Delete one sys role record by id
// @Description  Delete one sys role record by id
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sys Role ID"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router       /sys/role/void/{id} [delete]
func InactiveSysRoleByID(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidRoleById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deactivated successfully", "id": id, "data": result})
}

// @Summary      List sys roles with pagination
// @Description  List sys roles with pagination
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RolesPageDto  true  "List Action Record Request Body"
// @Success      200      {object}  dto.RolesList
// @Router       /sys/role/list [post]
func ListPageRoles(c *gin.Context) {
	var rolePageDto dto.RolesPageDto
	if err := c.ShouldBindJSON(&rolePageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.RolesList(&rolePageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// @Summary      Update sys role menu permission
// @Description  Update sys role menu permission
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Param        request  body      dto.MenuItemPermissionBody  true  "Update Request Body"
// @Success      200      {object}  dto.MessageResponse
// @Router       /sys/role/update-permission [post]
func HandleMeunPermission(c *gin.Context) {
	var menuPermissionDto dto.MenuItemPermissionBody
	if err := c.ShouldBindJSON(&menuPermissionDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.HandleMenuPermission(&menuPermissionDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu permissions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Menu permissions updated successfully", "data": result})
}

// @Summary      List sys role menu permission
// @Description  List sys role menu permission
// @Tags         Sys Role
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RoleIdsBody  true  "Update Request Body"
// @Success      200      {object}  []dto.SysRolesWithMenus
// @Router       /sys/role/list-permission [post]
func LoadRoleWithMenu(c *gin.Context) {
	var roleIdsDto dto.RoleIdsBody
	if err := c.ShouldBindJSON(&roleIdsDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.LoadRoleWithMenu(&roleIdsDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load roles with menus"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func RegisterSysRoleRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	roleGroup := rg.Group("/sys/role", handle.MiddlewareFunc())
	{
		roleGroup.POST("/create", CreateSysRoleApi)
		roleGroup.GET("/one/:id", GetRoleById)
		roleGroup.POST("/update", UpdateSysRole)
		roleGroup.DELETE("/void/:id", InactiveSysRoleByID)
		roleGroup.POST("/list", ListPageRoles)
		roleGroup.POST("/update-permission", HandleMeunPermission)
		roleGroup.GET("/all", GetAllRoles)
		roleGroup.POST("/list-permission", LoadRoleWithMenu)
	}
}
