package routes

import (
	"web_app/controller"
	_ "web_app/docs"
	"web_app/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func RegisterRouter() *gin.Engine {
	route := gin.New()
	route.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.GinLogger())
	route.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	userGroup := route.Group("/user")
	{
		userGroup.GET("/list", middlewares.JWTAuth(), middlewares.IsAdmin(), controller.GetUserListController)
		userGroup.POST("/register", controller.RegisterUserController) //
		userGroup.POST("/login", controller.LoginController)
	}
	return route
}
