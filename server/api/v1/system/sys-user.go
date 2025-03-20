package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type SysUserApi struct{}

func (s *SysUserApi) Create(ctx *gin.Context) {
	var req request.CreateSysUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"msg": req})
}
