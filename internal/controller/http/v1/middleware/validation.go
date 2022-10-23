package mw

import (
	"balance_api/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Msg string `json:"error"`
}

// ValidateJSONBody binds request body to given json struct
func ValidateJSONBody[BodyType any](l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body BodyType
		err := c.ShouldBindJSON(&body)
		if err != nil {
			l.Infof("validation err: %s", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Invalid request body format"})
			return
		}
		c.Set("jsonBody", body)
		c.Next()
	}
}

// GetJSONBody returns bound into given json struct request body
func GetJSONBody[BodyType any](c *gin.Context) BodyType {
	return c.MustGet("jsonBody").(BodyType)
}

// ValidateQuery binds request query to given json struct
func ValidateQuery[QueryType any](l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query QueryType
		err := c.ShouldBindQuery(&query)
		if err != nil {
			l.Infof("validation err: %s", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Invalid request query"})
			return
		}
		c.Set("queryParams", query)
		c.Next()
	}
}

// GetQueryParams returns bound into given json struct request query
func GetQueryParams[QueryType any](c *gin.Context) QueryType {
	return c.MustGet("queryParams").(QueryType)
}
