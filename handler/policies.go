package handler

import (
	"github.com/azi-v/ladon-api/domain"
	"github.com/azi-v/ladon-api/log"
	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"
)

type PolicyAddReq = ladon.DefaultPolicy

func PolicyAdd(c *gin.Context) {
	// 校验参数
	var pol *PolicyAddReq
	err := c.ShouldBindJSON(pol)
	if err != nil {
		log.Error(c, err)
		return
	}
	// 添加策略
	err = domain.Warden.Manager.Create(pol)
	if err != nil {
		log.Error(c, err)
		return
	}
}
