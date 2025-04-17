package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system/request"
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
	_, err := roleService.Create(req)
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

// Delete 删除角色
func (s *SysRoleApi) Delete(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}

	if err := roleService.Delete(id); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}
	response.NoContent(ctx)
}

// Update 更新角色
func (s *SysRoleApi) Update(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}
	var req request.UpdateSysRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	req.Id = id

	if err := roleService.Update(req); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}
	response.NoContent(ctx)
}

// Get 获取角色
func (s *SysRoleApi) Get(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}
	user, err := roleService.Get(id)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
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
	data, err := roleService.GetList(req)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
	}
	response.Success(ctx, data)
}
