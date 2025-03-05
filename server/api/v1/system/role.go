package system

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/utils"
)

type RoleApi struct{}

func (r *RoleApi) CreateRole(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.CreateRoleRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := roleService.Create(req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (r *RoleApi) DeleteRole(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.RoleParam
	var err error
	idStr := ctx.Param("roleId")
	param.RoleId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("role_id 类型错误", ctx)
		return
	}

	err = roleService.Delete(param, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (r *RoleApi) UpdateRole(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.RoleParam
	var err error
	idStr := ctx.Param("roleId")
	param.RoleId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("role_id 类型错误", ctx)
		return
	}

	req := request.UpdateRoleRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = roleService.Update(param, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (r *RoleApi) GetRole(ctx *gin.Context) {
	var param request.RoleParam
	var err error
	idStr := ctx.Param("roleId")
	param.RoleId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("role_id 类型错误", ctx)
		return
	}

	data, err := roleService.Find(param)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}

func (r *RoleApi) GetRoleList(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	query := request.RoleQueryParams{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(utils.GetValidMsg(query, err), ctx)
		return
	}

	data, err := roleService.FindList(query, common.ScopeData{Claims: claims})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}

func (r *RoleApi) UpdateRoleDataScope(ctx *gin.Context) {
	// data, _ := ctx.Get(global.CLAIMS)
	// claims := data.(system.Claims)

	var param request.RoleParam
	var err error
	idStr := ctx.Param("roleId")
	param.RoleId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("role_id 类型错误", ctx)
		return
	}

	req := request.UpdateRoleDataScope{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	// TODO: 暂时不支持更新角色数据范围, 最后来实现该功能
	response.NotImplemented(ctx)
	// err = roleService.UpdateDataScope(param, req, claims.UserId)
	// if err != nil {
	// 	response.BadRequest(err.Error(), ctx)
	// 	return
	// }

	// response.NoContent(ctx)
}

func (r *RoleApi) UpdateMenuDataScope(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.RoleParam
	var err error
	idStr := ctx.Param("roleId")
	param.RoleId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("role_id 类型错误", ctx)
		return
	}

	req := request.UpdateMenuDataScope{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = roleService.UpdateMenuScope(param, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (r *RoleApi) UpdateApiDataScope(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.RoleParam
	var err error
	idStr := ctx.Param("roleId")
	param.RoleId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("role_id 类型错误", ctx)
		return
	}

	req := request.UpdateApiDataScope{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = roleService.UpdateApiScope(param, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}
