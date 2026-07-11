package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func PimpinanRoutes(r *gin.Engine) {
	pimpinan := controller.NewPimpinanController()

	api := r.Group("/api/pimpinan")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("pimpinan", "pimpinan_fakultas", "Pimpinan Fakultas"))
	{
		api.GET("/dashboard", pimpinan.Dashboard)
		api.GET("/pengaduan/urgensi-tinggi", pimpinan.GetUrgensiTinggi)
		api.POST("/pengaduan/:id/disposisi", pimpinan.CreateDisposisi)
		api.GET("/disposisi", pimpinan.MyDisposisi)
	}
}
