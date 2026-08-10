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
		Id:       id,
		Avatar:   req.Avatar,
		Nickname: req.Nickname,
		RoleIds:  req.RoleIds,
	}

	if err := userService.Update(util.DB(ctx), user); err != nil {
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
	req.Id = id

	err := userService.ResetPassword(util.DB(ctx), req, nil, "")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}
