package v1

import (
	"balance_api/internal/usecase"
	"balance_api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// NewRouter is an entry point to controller layer: it sets up middleware for "/" route
// and groups routers by version
func NewRouter(handler *gin.Engine, b usecase.Balance, l logger.Interface) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	h := handler.Group("/v1")
	{
		newBalanceRoutes(h, b, l)
	}
}
