package auth

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/tools"
	"golang-fixedasset-mongo-backend/backend/services"
)

type AuthUser struct {
	Username string `json:"username"`
}

func AuthPayloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*AuthUser); ok {
			return jwt.MapClaims{
				"username": v.Username,
			}
		}
		return jwt.MapClaims{}
	}
}

func AuthIdentityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		username, ok := claims[identityKey].(string)
		if !ok {
			return nil // 導致 Unauthorized
		}
		return &AuthUser{
			Username: username,
		}
	}
}

func AuthUnauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func LoginAuthenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals dto.Login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", err
		}
		
		userID := loginVals.Username
		password := loginVals.Password

		if userID == "" || password == "" {
			return nil, jwt.ErrMissingLoginValues
		}

		hashedPassword := tools.HashPassword(password, tools.Salt)
		user, err := services.GetUserByUsername(userID)
		if err != nil || user.Status != 1 {
			return nil, jwt.ErrFailedAuthentication
		}

		if user.Password != hashedPassword {
			services.AddLoginRecord(user.Username, loginVals.IpAddress, "failed")
			return nil, jwt.ErrFailedAuthentication
		} else {
			services.AddLoginRecord(user.Username, loginVals.IpAddress, "success")
			return &AuthUser{
				Username: user.Username,
			}, nil
		}
	}
}

func Authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		v := data.(*AuthUser)
		user, err := services.GetUserByUsername(v.Username) // 🔁
		if err != nil || user.Status != 1 {
			return false
		}
		return true
	}
}
