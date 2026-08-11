package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/cast"
)

type SysJobApi struct{}

// Create 创建定时任务
func (s *SysJobApi) Create(ctx *gin.Context) {
	var req request.CreateSysJobReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	_, err := jobService.Create(util.DB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Created(ctx)
}

// Delete 删除定时任务
func (s *SysJobApi) Delete(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}

	if err := jobService.Delete(util.DB(ctx), id); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Update 更新定时任务
func (s *SysJobApi) Update(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.UpdateSysJobReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := jobService.Update(util.DB(ctx), id, req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Get 获取单个定时任务
func (s *SysJobApi) Get(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	menu, err := jobService.Get(util.DB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, menu)
}

// GetList 获取定时任务列表
func (s *SysJobApi) GetList(ctx *gin.Context) {
	var req request.QuerySysJobReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := jobService.GetList(util.DB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// Start 启动任务
func (s *SysJobApi) Start(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	err := jobService.Start(util.DB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Stop 停止任务
func (s *SysJobApi) Stop(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	err := jobService.Stop(util.DB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Execute 立即执行一次
func (s *SysJobApi) Execute(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	err := jobService.Execute(util.DB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}
