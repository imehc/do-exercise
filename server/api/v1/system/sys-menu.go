package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/spf13/cast"
)

type SysMenuApi struct{}

// Create 创建菜单
func (s *SysMenuApi) Create(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	var req request.CreateSysMenuReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	_, err := menuService.Create(req)
	if err != nil {
		// TODO: 处理具体错误
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.Created(ctx)
}

// Delete 删除菜单
func (s *SysMenuApi) Delete(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("idCannotBeEmpty", lang),
		})
		return
	}

	if err := menuService.Delete(id); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// Update 更新菜单
func (s *SysMenuApi) Update(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("idCannotBeEmpty", lang),
		})
		return
	}
	var req request.UpdateSysMenuReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	req.Id = id

	if err := menuService.Update(req); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// GetList 获取菜单详情
func (s *SysMenuApi) Get(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("idCannotBeEmpty", lang),
		})
		return
	}
	menu, err := menuService.Get(id)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.Success(ctx, menu)
}

// GetList 获取菜单树
func (s *SysMenuApi) GetTree(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	tree, err := menuService.GetTree()
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.Success(ctx, tree)
}
