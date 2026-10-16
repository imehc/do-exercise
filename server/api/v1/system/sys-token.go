package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
	systemService "github.com/imehc/do-exercise/server/service/system"
)

type SysTokenApi struct{}

// tokenScope 从已通过 AuthMiddleware 校验的会话里取出可见范围。
// 会话数据存在 Redis，不走 GORM，租户插件无法介入，只能在服务层显式裁剪。
func tokenScope(ctx *gin.Context) systemService.TokenScope {
	return systemService.TokenScope{
		TenantId:     ctx.GetString("tenantId"),
		IsSuperAdmin: ctx.GetBool("isSuperAdmin"),
	}
}

// GetList 获取token列表
func (s *SysTokenApi) FindAll(ctx *gin.Context) {
	data, err := sysTokenService.FindAll(tokenScope(ctx))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// Delete 删除token
func (s *SysTokenApi) Delete(ctx *gin.Context) {
	var req request.SysTokenDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := sysTokenService.Delete(req, tokenScope(ctx))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// ModityStatus 修改token状态
func (s *SysTokenApi) ModityStatus(ctx *gin.Context) {
	var req request.SysTokenModityStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := sysTokenService.ModityStatus(req, tokenScope(ctx))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}
