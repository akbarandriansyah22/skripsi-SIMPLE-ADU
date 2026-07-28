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
		api.GET("/kategori", pengaduan.Categories)
		api.GET("/tiket/:kode", pengaduan.DetailByTicket)
		api.GET("/:id", pengaduan.Detail)
		api.GET("/:id/lampiran", pengaduan.Attachment)
		api.GET("/:id/respon/:responId/lampiran", pengaduan.ResponseAttachment)
		api.PUT("/:id", middleware.RoleMiddleware("mahasiswa"), pengaduan.Update)
		api.POST("/:id/respon", middleware.RoleMiddleware("mahasiswa", "admin_sistem", "admin_fakultas", "pimpinan_fakultas"), pengaduan.AddRespon)
		api.PATCH("/:id/selesai", middleware.RoleMiddleware("mahasiswa"), pengaduan.Finish)
	}
}
