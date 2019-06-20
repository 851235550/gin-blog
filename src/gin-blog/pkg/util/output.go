package util

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"gin-blog/src/gin-blog/pkg/e"
)

func OutJson(c *gin.Context, code int, data interface{}) {
	//var c *gin.Context
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
