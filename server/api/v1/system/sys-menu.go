package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/spf13/cast"
)

type SysMenuApi struct{}

func (s *SysMenuApi) Create(ctx *gin.Context) {
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
			Message: err.Error(),
		})
		return
	}
	response.Created(ctx)
}

func (s *SysMenuApi) Delete(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}

	if err := menuService.Delete(id); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}
	response.NoContent(ctx)
}

func (s *SysMenuApi) Update(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
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
			Message: err.Error(),
		})
		return
	}
	response.NoContent(ctx)
}

func (s *SysMenuApi) Get(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: "id is required",
		})
		return
	}
	menu, err := menuService.Get(id)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}
	response.Success(ctx, menu)
}

func (s *SysMenuApi) GetTree(ctx *gin.Context) {
	tree, err := menuService.GetTree()
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
	}
	response.Success(ctx, tree)
}
