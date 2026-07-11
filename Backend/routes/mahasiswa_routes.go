package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func MahasiswaRoutes(r *gin.Engine) {
	mahasiswa := controller.NewMahasiswaController()

	api := r.Group("/api/mahasiswa")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("mahasiswa"))
	{
		api.GET("/profile", mahasiswa.Profile)
		api.PUT("/profile", mahasiswa.UpdateProfile)
	}
}
