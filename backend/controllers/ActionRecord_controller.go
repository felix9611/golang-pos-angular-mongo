package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// PingExample godoc
// @Summary      List Action Records
// @Description  List Action Records
// @Tags         Action Records
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListActionRecordReqDto  true  "List Action Record Request Body"
// @Success      200      {object}  dto.ActionReocrdList
// @Router       /action-records/list [post]
func ListPageActionRecords(c *gin.Context) {
	var pageDto dto.ListActionRecordReqDto
	if err := c.ShouldBindJSON(&pageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actionRecords, err := services.ListActionRecords(&pageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list action records"})
		return
	}
	c.JSON(http.StatusOK, actionRecords)
}

func RegisterActionRecordsRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	ActionRecordGroup := r.Group("/action-records", authMiddleware.MiddlewareFunc())
	{
		ActionRecordGroup.POST("/list", ListPageActionRecords)
	}
}
