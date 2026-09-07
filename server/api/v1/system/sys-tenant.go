package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/cast"
)

type SysTenantApi struct{}

// Create 创建租户（平台层）
func (s *SysTenantApi) Create(ctx *gin.Context) {
	var req request.CreateSysTenantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	_, err := tenantService.Create(util.BypassTenantDB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Created(ctx)
}

// Update 更新租户
func (s *SysTenantApi) Update(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.UpdateSysTenantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := tenantService.Update(util.BypassTenantDB(ctx), id, req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Delete 删除租户
func (s *SysTenantApi) Delete(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	if err := tenantService.Delete(util.BypassTenantDB(ctx), id); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// Get 获取租户详情
func (s *SysTenantApi) Get(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	tenant, err := tenantService.Get(util.BypassTenantDB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, tenant)
}

// GetList 获取租户列表
func (s *SysTenantApi) GetList(ctx *gin.Context) {
	var req request.QuerySysTenantReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := tenantService.GetList(util.BypassTenantDB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// ListAssignableAdmins 获取可被选作租户管理员的现有用户列表（创建租户时用）
func (s *SysTenantApi) ListAssignableAdmins(ctx *gin.Context) {
	users, err := tenantService.ListAssignableAdmins(util.BypassTenantDB(ctx))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, users)
}

// ListAssignableUsers 获取可分配给指定租户的现有用户列表
func (s *SysTenantApi) ListAssignableUsers(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	users, err := tenantService.ListAssignableUsers(util.BypassTenantDB(ctx), id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, users)
}

// AssignUsers 分配现有用户到指定租户
func (s *SysTenantApi) AssignUsers(ctx *gin.Context) {
	id := cast.ToString(ctx.Param("id"))
	if id == "" {
		response.BadRequest(ctx, "idCannotBeEmpty")
		return
	}
	var req request.AssignTenantUsersReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := tenantService.AssignUsers(util.BypassTenantDB(ctx), id, req.UserIds); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}
