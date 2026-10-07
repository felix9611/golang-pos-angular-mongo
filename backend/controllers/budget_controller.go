package controllers

import (
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create one Budget record
// @Description  Create one Budget record
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        request  body      models.Budgets  true  "Budget Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /base/budget/create [post]
func CreateBudgetRecord(c *gin.Context) {
	var budgetDto models.Budgets
	if err := c.ShouldBindJSON(&budgetDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateBudgetRecord(&budgetDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get one budget record by id
// @Description  Get one budget record by id
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "budget ID"
// @Success      200      {object}  models.Budgets
// @Router       /base/budget/one/{id} [get]
func GetBudgetRecordById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetOneBudgetRecord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      List Budgets
// @Description  List Budgets
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListBudgetRecordsDto  true  "List Budget Request Body"
// @Success      200      {object}  dto.BudgetList
// @Router       /sys/budget/list [post]
func ListPageBudgetRecords(c *gin.Context) {
	var budgetPageDto dto.ListBudgetRecordsDto
	if err := c.ShouldBindJSON(&budgetPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budgets, err := services.ListBudgetRecords(&budgetPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list budgets"})
		return
	}
	c.JSON(http.StatusOK, budgets)
}

// @Summary      List Budgets With Filter
// @Description  List Budgets With Filter
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListBudgetRecordsDto  true  "List Budget Request Body"
// @Success      200      {object}  []models.Budgets
// @Router       /sys/budget/filter/list [post]
func ListPageBudgetRecordsWithFilter(c *gin.Context) {
	var budgetPageDto dto.ListBudgetRecordsDto
	if err := c.ShouldBindJSON(&budgetPageDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budgets, err := services.ListBudgetRecordsFilter(&budgetPageDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list budgets"})
		return
	}
	c.JSON(http.StatusOK, budgets)
}

// @Summary      Get budget summary
// @Description  Get budget summary
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.BudgetSummary
// @Router      /base/budget/getBudgetSummary [get]
func GetBudgetSummary(c *gin.Context) {
	result, err := services.GetBudgetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Update one budget record
// @Description  Update one budget record
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        request  body      models.Budgets  true  "Budget Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /base/budget/update [post]
func UpdateBudgetRecord(c *gin.Context) {
	var budgetDto models.Budgets
	if err := c.ShouldBindJSON(&budgetDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateBudgetRecord(&budgetDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Delete one budget record by id
// @Description  Delete one budget record by id
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Budget ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router      /base/budget/void/{id} [delete]
func InactiveBudgetRecord(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidBudgetRecord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Budget record inactivated successfully", "data": result})
}

func RegisterBudgetRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	BudgetGroup := rg.Group("/base/budget", handle.MiddlewareFunc())
	{
		BudgetGroup.POST("/create", CreateBudgetRecord)
		BudgetGroup.GET("/one/:id", GetBudgetRecordById)
		BudgetGroup.POST("/list", ListPageBudgetRecords)
		BudgetGroup.GET("/getBudgetSummary", GetBudgetSummary)
		BudgetGroup.POST("/update", UpdateBudgetRecord)
		BudgetGroup.DELETE("/void/:id", InactiveBudgetRecord)
		BudgetGroup.POST("/filter/list", ListPageBudgetRecordsWithFilter)
	}
}
