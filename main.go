package main

import (
	"net/http"

	"git.ymt360.com/usercenter/ymt-ladon/api"
	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: 增加优雅退出等配置
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "OK"}) })
	api.HTTPRouter(r)
	r.Run()
}
