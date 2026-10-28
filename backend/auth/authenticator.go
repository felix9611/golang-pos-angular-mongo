package auth

import (
	"golang-pos-mongo-backend/backend/tools"
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	identityKey = "username"
	port        string
)

func HandlerMiddleWare(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	if err := authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}
	return authMiddleware.MiddlewareFunc()
}

func InitAuthParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(os.Getenv("JWT_SECRET_KEY")),
		Timeout:     time.Hour * 24 * 1,
		MaxRefresh:  time.Hour * 24 * 1,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*AuthUser); ok {
				return jwt.MapClaims{
					identityKey:  v.Username,
					"roleIds":    v.RoleIds,
					"placeIds":   v.PlaceIds,
					"deptIds":    v.DeptIds,
					"superAdmin": v.SuperAdmin,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			username, ok := claims[identityKey].(string)
			if !ok {
				return nil
			}
			user := &AuthUser{
				Username: username,
			}
			if v, ok := claims["superAdmin"].(bool); ok {
				user.SuperAdmin = v
			}

			user.RoleIds = tools.ClaimStrings(claims["roleIds"])
			user.PlaceIds = tools.ClaimStrings(claims["placeIds"])
			user.DeptIds = tools.ClaimStrings(claims["deptIds"])

			return user
		},
		Authenticator: LoginAuthenticator(),
		Authorizator:  Authorizator(),
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}
