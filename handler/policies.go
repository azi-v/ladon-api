package handler

import (
	"git.ymt360.com/usercenter/ymt-ladon/domain"
	"git.ymt360.com/usercenter/ymt-ladon/log"
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
