package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Batch Create location record
// @Description  Batch Create location record
// @Tags         Tax Information
// @Accept       json
// @Produce      json
// @Param        request  body      []models.Locations  true  "Locations Request Body"
// @Success      200      string
// @Router      /base/location/batch-create [post]
func BatchCreateLocation(c *gin.Context) {
	var locations []models.Locations
	if err := c.ShouldBindJSON(&locations); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.BatchCreateLocation(locations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Create one Location record
// @Description  Create one Location record
// @Tags         Location
// @Accept       json
// @Produce      json
// @Param        request  body      models.Locations  true  "Location Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /base/location/create [post]
func CreateLocation(c *gin.Context) {
	var location models.Locations
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateLocation(&location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one location record by id
// @Description  Get one location record by id
// @Tags         Location
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Location ID"
// @Success      200      {object}  models.Locations
// @Router       /base/location/one/{id} [get]
func GetLocationById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneLocationById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Delete one location record by id
// @Description  Delete one location record by id
// @Tags         Location
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Location ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router       /base/location/void/{id} [delete]
func InactiveLocationByID(c *gin.Context) {
	id := c.Param("id")
	result, err := services.InactiveLocationByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deactivated successfully", "id": id, "data": result})
}

// @Summary      Update one Location record
// @Description  Update one Location record
// @Tags         Location
// @Accept       json
// @Produce      json
// @Param        request  body      models.Locations  true  "Location Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /base/location/update/{id} [post]
func UpdateLocationById(c *gin.Context) {
	var location models.Locations
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateLocationById(&location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Locations
// @Description  List Locations
// @Tags         Location
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LocationPageDto  true  "List Location Request Body"
// @Success      200      {object}  dto.LocationList
// @Router       /base/location/list [post]
func ListLocations(c *gin.Context) {
	var locationPageDto dto.LocationPageDto
	if err := c.ShouldBindJSON(&locationPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.LocationList(&locationPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      List Locations Without pagination
// @Description  List Locations Without pagination
// @Tags         Location
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LocationPageDto  true  "List Location Request Body"
// @Success      200      {object}  []models.Locations
// @Router       /base/location/filter/list [post]
func LocationListWithFilter(c *gin.Context) {
	var locationPageDto dto.LocationPageDto
	if err := c.ShouldBindJSON(&locationPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.LocationListWithFilter(&locationPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get all location record
// @Description  Get all location record
// @Tags         Location
// @Produce      json
// @Success      200      {object}  []models.Locations
// @Router       /base/location/all [get]
func GetAllLocations(c *gin.Context) {
	locations, err := services.ListAllLocation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get locations"})
		return
	}
	c.JSON(http.StatusOK, locations)
}

func RegisterLocationRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	locationGroup := rg.Group("/base/location", handle.MiddlewareFunc())
	{
		locationGroup.POST("/create", CreateLocation)
		locationGroup.GET("/one/:id", GetLocationById)
		locationGroup.POST("/update", UpdateLocationById)
		locationGroup.DELETE("/void/:id", InactiveLocationByID)
		locationGroup.POST("/list", ListLocations)
		locationGroup.GET("/all", GetAllLocations)
		locationGroup.POST("/filter/list", LocationListWithFilter)
		locationGroup.POST("/batch-create", BatchCreateLocation)
	}
}
