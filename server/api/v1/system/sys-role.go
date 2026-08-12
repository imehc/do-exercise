package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/cast"
)

type SysRoleApi struct{}

// Create 创建角色
func (s *SysRoleApi) Create(ctx *gin.Context) {
	var req request.CreateSysRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	_, err := roleService.Create(util.DB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Created(ctx)
}

// Delete 删除角色
func (s *SysRoleApi) Delete(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}

	if err := roleService.Delete(util.DB(ctx), id); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Update 更新角色
func (s *SysRoleApi) Update(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.UpdateSysRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := roleService.Update(util.DB(ctx), id, req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Get 获取角色详情
func (s *SysRoleApi) Get(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	user, err := roleService.Get(util.DB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, user)
}

// GetList 获取角色列表
func (s *SysRoleApi) GetList(ctx *gin.Context) {
	var req request.QuerySysRoleReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := roleService.GetList(util.DB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// GetAll 获取所有角色
func (s *SysRoleApi) GetAll(ctx *gin.Context) {
	data, err := roleService.GetAll(util.DB(ctx))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
