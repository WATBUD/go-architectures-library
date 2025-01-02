package handler

import (
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

// Response sends a JSON response with the given HTTP status code, custom code, message, and data.
// httpCode: the HTTP status code to be sent in the response
// code: a custom application-specific code
// message: a message to be included in the response
// data: the data to be included in the response
func (g *Gin) Response(httpCode, code int, message string, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"http_status": httpCode,
		"code":        code,
		"msg":         message,
		"data":        data,
	})
}
