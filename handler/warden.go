package handler

import (
	"git.ymt360.com/usercenter/ymt-ladon/domain"
	"git.ymt360.com/usercenter/ymt-ladon/log"
	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"
)

type WardenAllowReq = ladon.Request

func WardenAllow(c *gin.Context) {
	// json
	req := &WardenAllowReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Error(c, err)
		return
	}
	err = domain.Warden.IsAllowed(req)
	if err != nil {
		// 不允许
	}

	// err = nil 允许

}
