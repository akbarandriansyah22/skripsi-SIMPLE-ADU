package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	auth := controller.NewAuthController()

	api := r.Group("/api/auth")
	{
		api.POST("/register", auth.Register)
		api.POST("/login", auth.Login)
		api.GET("/profile", middleware.AuthMiddleware(), auth.Profile)
		api.PATCH("/change-password", middleware.AuthMiddleware(), auth.ChangePassword)
	}

}
