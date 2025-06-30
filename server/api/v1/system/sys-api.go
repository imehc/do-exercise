package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/spf13/cast"
)

type SysApiApi struct{}

// Update 更新api
func (s *SysApiApi) Update(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.UpdateSysApiReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	req.Id = id

	if err := apiService.Update(req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Get 获取api详情
func (s *SysApiApi) Get(ctx *gin.Context) {
	id := cast.ToUint(ctx.Param("id"))
	if id == 0 {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	api, err := apiService.Get(id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, api)
}

// GetList 获取api列表
func (s *SysApiApi) GetList(ctx *gin.Context) {
	var req common.Pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := apiService.GetList(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// GetAll 获取所有api
func (s *SysApiApi) GetAll(ctx *gin.Context) {
	data, err := apiService.GetAll()
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// GetGroupType 获取分组类型
func (s *SysApiApi) GetGroupType(ctx *gin.Context) {
	data, err := apiService.GroupType()
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
