package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/histopathai/image-catalog-service/config"
	"github.com/histopathai/image-catalog-service/internal/handlers"
)

func SetupRouter(imageHandler *handlers.ImageHandler, gcsProxyHandler *handlers.GCSProxyHandler, cfg *config.Config) *gin.Engine {

	gin.SetMode(cfg.Server.GINMode)
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/images/:image_id", imageHandler.GetImageByID)
		apiV1.PUT("/images/:image_id", imageHandler.UpdateImageByID)
		apiV1.DELETE("/images/:image_id", imageHandler.DeleteImageByID)
		apiV1.GET("/images", imageHandler.GetImages)

		// 🔥 Wildcard route to proxy all GCS objects
		apiV1.GET("/proxy/*objectPath", gcsProxyHandler.ProxyObject)
	}

	return router
}
