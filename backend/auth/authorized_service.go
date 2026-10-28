package auth

import (
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/services"
	"golang-pos-mongo-backend/backend/tools"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthUser struct {
	Username   string               `json:"username"`
	RoleIds    []primitive.ObjectID `json:"roleIds"`
	PlaceIds   []primitive.ObjectID `json:"placeIds"`
	DeptIds    []primitive.ObjectID `json:"deptIds"`
	SuperAdmin bool                 `json:"superAdmin"`
}

func AuthPayloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*AuthUser); ok {
			return jwt.MapClaims{
				"username":   v.Username,
				"roleIds":    v.RoleIds,
				"placeIds":   v.PlaceIds,
				"deptIds":    v.DeptIds,
				"superAdmin": v.SuperAdmin,
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
				Username:   user.Username,
				RoleIds:    user.Roles,
				PlaceIds:   user.PlaceIds,
				DeptIds:    user.DeptIds,
				SuperAdmin: user.SuperAdmin,
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
