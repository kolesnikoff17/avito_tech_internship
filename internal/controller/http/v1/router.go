package v1

import (
	"balance_api/internal/usecase"
	"balance_api/pkg/logger"
	"github.com/gin-gonic/gin"

	// Swagger docs
	_ "balance_api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter is an entry point to controller layer: it sets up middleware for "/" route
// and groups routers by version
func NewRouter(handler *gin.Engine, b usecase.Balance, l logger.Interface) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	h := handler.Group("/v1")
	{
		newBalanceRoutes(h, b, l)
	}
}
