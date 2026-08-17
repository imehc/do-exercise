package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/cast"
)

type SysUserApi struct{}

// Create 创建用户
func (s *SysUserApi) Create(ctx *gin.Context) {
	var req request.CreateSysUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	_, err := userService.Create(util.DB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Created(ctx)
}

// Delete 删除用户
func (s *SysUserApi) Delete(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}

	if err := userService.Delete(util.DB(ctx), id); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Update 更新用户
func (s *SysUserApi) Update(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.UpdateSysUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	user := request.UpdateSysUserReq{
		Avatar:   req.Avatar,
		Nickname: req.Nickname,
		RoleIds:  req.RoleIds,
	}

	if err := userService.Update(util.DB(ctx), id, user); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Get 获取用户详情
func (s *SysUserApi) Get(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	user, err := userService.Get(util.DB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, user)
}

// UsernameExists 检查用户名是否已存在
func (s *SysUserApi) UsernameExists(ctx *gin.Context) {
	username := cast.ToString(ctx.Query("username"))
	if username == "" {
		response.BadRequest(ctx, "usernameRequired")
		return
	}
	exists, err := userService.UsernameExists(util.DB(ctx), username)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, gin.H{"exists": exists})
}

// GetList 获取用户列表
func (s *SysUserApi) GetList(ctx *gin.Context) {
	var req common.Pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := userService.GetList(util.DB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// ResetPassword 重置密码
func (s *SysUserApi) ResetPassword(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.UpdateSysUserPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := userService.ResetPassword(util.DB(ctx), id, req, nil, "")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// AssignTenant 给用户分配（复制到）目标租户，仅平台超级管理员可调用
func (s *SysUserApi) AssignTenant(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.AssignUserTenantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	err := userService.AssignUserToTenant(util.BypassTenantDB(ctx), id, req.TenantId)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// ListAssignableTenants 获取指定用户可分配的候选租户列表
func (s *SysUserApi) ListAssignableTenants(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	tenants, err := userService.ListAssignableTenants(util.BypassTenantDB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, tenants)
}
