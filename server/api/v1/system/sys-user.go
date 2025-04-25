package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/spf13/cast"
)

type SysUserApi struct{}

// Create 创建用户
func (s *SysUserApi) Create(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	var req request.CreateSysUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	_, err := userService.Create(req)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.Created(ctx)
}

// Delete 删除用户
func (s *SysUserApi) Delete(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	id := cast.ToInt64(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("idCannotBeEmpty", lang),
		})
		return
	}

	if err := userService.Delete(id); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// Update 更新用户
func (s *SysUserApi) Update(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	id := cast.ToInt64(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("idCannotBeEmpty", lang),
		})
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

	if err := userService.Update(user); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// Get 获取用户详情
func (s *SysUserApi) Get(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	id := cast.ToInt64(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("idCannotBeEmpty", lang),
		})
		return
	}
	user, err := userService.Get(id)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.Success(ctx, user)
}

// GetList 获取用户列表
func (s *SysUserApi) GetList(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	var req common.Pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := userService.GetList(req)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
	}
	response.Success(ctx, data)
}

// ResetPassword 重置密码
func (s *SysUserApi) ResetPassword(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	id := cast.ToInt64(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("idCannotBeEmpty", lang),
		})
		return
	}
	var req request.UpdateSysUserPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	err := userService.ResetPassword(req, nil)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}
