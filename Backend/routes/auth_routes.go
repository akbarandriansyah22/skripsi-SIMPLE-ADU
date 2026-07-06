package routes

import (
	controller "backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	auth := controller.NewAuthController()

	api := r.Group("/api/auth")
	{
		api.POST("/register", auth.Register)
		api.POST("/login", auth.Login)
	}

}