package main

import (
	"golang-fixedasset-mongo-backend/backend/auth"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/controllers"
	"time"

	_ "golang-fixedasset-mongo-backend/backend/docs"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Asset Management API
// @version         1.0
// @description     This is an API server for asset management system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:6500
func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	config.ConnectDatabase()

	authServiceMiddleware, errAuthService := jwt.New(auth.InitAuthParams())
	if errAuthService != nil {
		panic("JWT Middleware initialization failed: " + errAuthService.Error())
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/")
	controllers.RegisterTongsRoutes(api)
	controllers.RegisterDepartmentRoutes(api, authServiceMiddleware)
	controllers.RegisterSysUserRoutes(api, authServiceMiddleware)
	// example.RegisterAuthExampleRoute(api, authMiddleware)
	controllers.RegisterAuthRoute(api, authServiceMiddleware)
	controllers.RegisterSysRoleRoutes(api, authServiceMiddleware)
	controllers.RegisterVendorRoutes(api, authServiceMiddleware)
	controllers.RegisterLocationRoutes(api, authServiceMiddleware)
	controllers.RegisterCodeTypeRoutes(api, authServiceMiddleware)
	controllers.RegisterTaxInformationRoutes(api, authServiceMiddleware)
	controllers.RegisterAssetTypeRoutes(api, authServiceMiddleware)
	controllers.RegisterInvRecordRoutes(api, authServiceMiddleware)
	controllers.RegisterAssetListRoutes(api, authServiceMiddleware)
	controllers.RegisterStockTakeRoutes(api, authServiceMiddleware)
	controllers.RegisterBudgetRoutes(api, authServiceMiddleware)
	controllers.RegisterActionRecordsRoutes(api, authServiceMiddleware)
	controllers.RegisterSysMenuRoutes(api, authServiceMiddleware)
	controllers.RegisterRepairRecordRoutes(api, authServiceMiddleware)
	controllers.RegisterExcelFieldMatchRoutes(api, authServiceMiddleware)
	controllers.RegisterWriteOffsRoutes(api, authServiceMiddleware)

	router.Run(":6500")
}
