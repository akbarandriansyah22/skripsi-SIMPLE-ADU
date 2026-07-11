package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func PengaduanRoutes(r *gin.Engine) {
	pengaduan := controller.NewPengaduanController()

	api := r.Group("/api/pengaduan")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("", middleware.RoleMiddleware("mahasiswa"), pengaduan.Create)
		api.GET("", middleware.RoleMiddleware("mahasiswa"), pengaduan.MyPengaduan)
		api.GET("/tiket/:kode", pengaduan.DetailByTicket)
		api.GET("/:id", pengaduan.Detail)
		api.PUT("/:id", middleware.RoleMiddleware("mahasiswa"), pengaduan.Update)
		api.POST("/:id/respon", middleware.RoleMiddleware("mahasiswa", "admin", "admin_fakultas", "Admin Fakultas", "pimpinan", "pimpinan_fakultas", "Pimpinan Fakultas"), pengaduan.AddRespon)
		api.PATCH("/:id/selesai", middleware.RoleMiddleware("mahasiswa"), pengaduan.Finish)
	}
}
