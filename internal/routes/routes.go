package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/histopathai/image-catalog-service/internal/handlers"
)

func SetupRouter(imageHandler *handlers.ImageHandler) *gin.Engine {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/images", imageHandler.UploadImage)
		apiV1.GET("/images/:image_id", imageHandler.GetImageByID)
		apiV1.PUT("/images/:image_id", imageHandler.UpdateImageByID)
		apiV1.DELETE("/images/:image_id", imageHandler.DeleteImageByID)
		apiV1.GET("/images", imageHandler.GetImages)
	}
	return router
}
