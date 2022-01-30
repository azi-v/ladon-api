package api

import (
	"github.com/azi-v/ladon-api/handler"
	"github.com/gin-gonic/gin"
)

func HTTPRouter(e *gin.Engine) {

	e.POST("/policy/add", handler.PolicyAdd)     // 添加策略
	e.POST("/warden/allow", handler.WardenAllow) // 检测权限

}
