package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
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
	user := system.SysUser{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Password: req.Password,
	}
	_, err := userService.Create(user)
	if err != nil {
		// TODO: 处理具体错误
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}
	response.Created(ctx)
}

// Delete 删除用户
func (s *SysUserApi) Delete(ctx *gin.Context) {
	idString := ctx.Param("id")
	id := cast.ToInt64(idString)
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}

	if err := userService.Delete(id); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}
	response.NoContent(ctx)
}

// Update 更新用户
func (s *SysUserApi) Update(ctx *gin.Context) {
	idString := ctx.Param("id")
	id := cast.ToInt64(idString)
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}
	var req request.UpdateSysUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	user := system.SysUser{
		Email:    req.Email,
		Avatar:   req.Avatar,
		Nickname: req.Nickname,
	}

	if err := userService.Update(user, id); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}
	response.NoContent(ctx)
}

// Get 获取用户
func (s *SysUserApi) Get(ctx *gin.Context) {
	idString := ctx.Param("id")
	id := cast.ToInt64(idString)
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}
	user, err := userService.Get(id)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
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
	data, err := userService.GetList(req)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
	}
	response.Success(ctx, data)
}
