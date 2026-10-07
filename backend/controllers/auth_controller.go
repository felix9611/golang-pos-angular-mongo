package controllers

import (
	"golang-fixedasset-mongo-backend/backend/auth"
	"golang-fixedasset-mongo-backend/backend/services"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	identityKey = "username"
	port        string
)

func HeyTongsHandler(c *gin.Context) {
	userRaw, exists := c.Get(identityKey)
	if !exists {
		c.JSON(401, gin.H{
			"message": "Unauthorized: identity not found",
		})
		return
	}

	user, ok := userRaw.(*auth.AuthUser)
	if !ok {
		c.JSON(401, gin.H{
			"message": "Unauthorized: invalid identity type",
		})
		return
	}

	c.JSON(200, gin.H{
		"username": user.Username,
		"text":     "Hey Tongs.",
	})
}

// @Summary      Get user details by token
// @Description  Get user details by token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.AuthDataResponse
// @Router       /asset/type/list [get]
func GetUserDetails(c *gin.Context) {
	userRaw, exists := c.Get(identityKey)
	if !exists {
		c.JSON(401, gin.H{
			"message": "Unauthorized: identity not found",
		})
		return
	}

	authUser, ok := userRaw.(*auth.AuthUser)
	if !ok {
		c.JSON(401, gin.H{
			"message": "Unauthorized: invalid identity type",
		})
		return
	}

	username := authUser.Username
	userDetails, err := services.GetUserDetailByUsername(username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error: failed to get user details",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": userDetails,
	})
}

func RegisterAuthRoute(r *gin.RouterGroup, handle *jwt.GinJWTMiddleware) {
	r.POST("/login", handle.LoginHandler)
	auth := r.Group("/auth", handle.MiddlewareFunc())
	{
		auth.GET("/refresh_token", handle.RefreshHandler)
		auth.GET("/verify-token", HeyTongsHandler)
		auth.GET("/user-profile", GetUserDetails)
	}
}
