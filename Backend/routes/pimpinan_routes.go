package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func PimpinanRoutes(r *gin.Engine) {
	pimpinan := controller.NewPimpinanController()

	api := r.Group("/api/pimpinan")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("pimpinan_fakultas"))
	{
		api.GET("/dashboard", pimpinan.Dashboard)
		api.GET("/pengaduan/urgensi-tinggi", pimpinan.GetUrgensiTinggi)
		api.GET("/pengaduan/riwayat", pimpinan.GetRiwayatUrgensiTinggi)
		api.GET("/pengaduan/riwayat/:id", pimpinan.GetRiwayatUrgensiTinggiDetail)
		api.GET("/pengaduan/:id", pimpinan.GetPengaduan)
		api.POST("/pengaduan/:id/disposisi", pimpinan.CreateDisposisi)
		api.GET("/disposisi", pimpinan.MyDisposisi)
		api.GET("/monitoring", pimpinan.Monitoring)
		api.GET("/unit/penanganan", pimpinan.GetAssignmentUnits)
	}
}
