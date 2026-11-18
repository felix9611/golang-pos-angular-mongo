package controllers

import (
	"encoding/json"
	"golang-fixedasset-mongo-backend/backend/auth"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/services"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary      Create job task
// @Description  Create job task
// @Tags         Task Board
// @Accept       json
// @Produce      json
// @Param        request  body      models.JobDashboard  true  "Job Dashboard Request Body"
// @Success      200      {object}  dto.GeneralCreateResponseBody
// @Router      /tasks/create [post]
func CreateJobBoard(c *gin.Context) {
	value, _ := c.Get(identityKey)
	user, ok := value.(*auth.AuthUser)
	if !ok || user == nil {
		c.JSON(401, gin.H{"error": "Invalid user"})
		return
	}

	var task models.JobDashboard
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !(user.SuperAdmin || user.Manager) && task.Status >= 3 {
		c.JSON(403, gin.H{
			"error": "You do not allow to create this task on this status",
		})
	} else {
		task.CreatedBy = user.Username
		id, err := services.CreateTasks(&task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

// @Summary      List all job task
// @Description  List all job task
// @Tags         Task Board
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListJobDashboardRequestDto  true  "Job Dashboard Request Body"
// @Success      200      {object}  models.JobDashboard
// @Router      /tasks/list [post]
func ListJobBoard(c *gin.Context) {
	value, _ := c.Get(identityKey)
	data, err := json.Marshal(value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
		return
	}

	var user dto.AuthUser

	if err := json.Unmarshal(data, &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user structure"})
		return
	}

	var task dto.ListJobDashboardRequestDto
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListJobDashboard(task, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Load failed"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get one task by id
// @Description  Get one task by id
// @Tags         Task Board
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200      {object}  models.JobDashboard
// @Router       /tasksget-one/{id} [get]
func GetTaskById(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetTaskById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Update one task
// @Description  Update one task
// @Tags         Job Board
// @Accept       json
// @Produce      json
// @Param        request  body      models.JobDashboard  true  "Location Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /tasks/update [post]
func UpdateTaskById(c *gin.Context) {
	value, _ := c.Get(identityKey)
	data, err := json.Marshal(value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
		return
	}

	var user dto.AuthUser

	if err := json.Unmarshal(data, &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user structure"})
		return
	}

	var task models.JobDashboard
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !(user.SuperAdmin || user.Manager) && task.Status >= 3 {
		c.JSON(403, gin.H{
			"error": "You do not allow to update this task on this status",
		})
	} else {
		result, err := services.UpdateTask(&task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// @Summary      Delete task by id
// @Description  Delete task by id
// @Tags         Job Board
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200      {object}  dto.GeneralUpdateInactiveResponseBody
// @Router       /tasks/void/{id} [delete]
func VoidTaskByID(c *gin.Context) {
	id := c.Param("id")
	result, err := services.VoidTaskById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deactivated successfully", "id": id, "data": result})
}

// @Summary      Moving one task
// @Description  Moving one task
// @Tags         Job Board
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdateJobDashboardStatusDto  true  "Job Dashboard Request Body for update"
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /tasks/change-status [post]
func UpdateJobDashboardStatus(c *gin.Context) {
	value, _ := c.Get(identityKey)
	data, err := json.Marshal(value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
		return
	}

	var user dto.AuthUser

	if err := json.Unmarshal(data, &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user structure"})
		return
	}

	var task *dto.UpdateJobDashboardStatusDto
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !(user.SuperAdmin || user.Manager) && task.Status > 3 {
		c.JSON(403, gin.H{
			"error": "You do not allow to update this task on this status",
		})
	} else {
		result, err := services.UpdateJobDashboardStatus(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// @Summary      Sync tasks
// @Description  Sync tasks
// @Tags         Job Board
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.GeneralUpdateResponse
// @Router      /tasks/sync [post]
func SyncJobDashboard(c *gin.Context) {
	value, _ := c.Get(identityKey)
	data, err := json.Marshal(value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
		return
	}

	var user dto.AuthUser

	if err := json.Unmarshal(data, &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user structure"})
		return
	}

	result, err := services.SyncTasks(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterTaskBoardsRoutes(rg *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	taskBoards := rg.Group("/tasks", handle.MiddlewareFunc())
	{
		taskBoards.POST("/create", CreateJobBoard)
		taskBoards.POST("/list", ListJobBoard)
		taskBoards.GET("/get-one/:id", GetTaskById)
		taskBoards.POST("/update", UpdateTaskById)
		taskBoards.POST("/change-status", UpdateJobDashboardStatus)
		taskBoards.DELETE("/void/:id", VoidTaskByID)
		taskBoards.GET("/sync", SyncJobDashboard)
	}
}
