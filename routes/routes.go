package routes

import (
	"di_importer/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", handlers.HomeHandler)
	r.POST("/upload", handlers.UploadHandler)
}
