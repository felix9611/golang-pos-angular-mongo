package controllers

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"golang-fixedasset-mongo-backend/backend/dto"
	"net/http"
//	"time"

	"github.com/gin-gonic/gin"

)

// var tongsCollection = config.GetCollection("tongs")

var u models.Tongs

func GetTongs(c *gin.Context) {
	// Placeholder for the actual logic to get a Tong
	// This should interact with the database to retrieve the Tong data
	c.JSON(http.StatusOK, gin.H{
		"message": "Hey Tongs",
	})
}

func CreateTongs(c *gin.Context) {
	
	var tong models.Tongs
	if err := c.ShouldBindJSON(&tong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := services.CreateTongs(&tong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func GetOneTongsRun(c *gin.Context) {
	id := c.Param("id")
	tongs, err := services.GetOneTongsRun(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tongs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tongs })
}

func UpdateTongsByID(c *gin.Context) {
	id := c.Param("id")
	var tong models.Tongs
	result, err := services.UpdateTongsByID(id, &tong)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tongs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update Tongs by ID", "id": id, "data": result})
}


func VoidTongs(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidTongs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to void tongs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tong voided successfully", "result": result})
}

func ListPageTongs(c *gin.Context) {
	var tongReq dto.TongsPageDto
	if err := c.ShouldBindJSON(&tongReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tongs, err := services.TongsListPage(&tongReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tongs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tongs})
	
}

func RegisterTongsRoutes(rg *gin.RouterGroup) {
	tongs := rg.Group("/tongs")
	{
		tongs.GET("/get", GetTongs)
		tongs.GET("/one/:id", GetOneTongsRun)
		tongs.POST("/create", CreateTongs)
		tongs.POST("/update/:id", UpdateTongsByID) // Uncomment when implemented
		tongs.DELETE("/void/:id", VoidTongs) // Uncomment when implemented
		tongs.POST("/list", ListPageTongs) // Uncomment when implemented
		// tongs.DELETE("/delete/:id", DeleteTongsByID) // Uncomment
		
	}
}

