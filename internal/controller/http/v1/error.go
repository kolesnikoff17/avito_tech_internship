package v1

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Msg string `json:"error"`
}

// ErrorResponse aborts current request with err
func ErrorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{Msg: msg})
}
