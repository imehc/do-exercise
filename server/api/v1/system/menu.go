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

type MenuApi struct{}

func (m *MenuApi) CreateMenu(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.MenuRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := menuService.CreateMenu(req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (m *MenuApi) DeleteMenu(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.MenuParam
	var err error
	idStr := ctx.Param("menuId")
	param.MenuId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("menu_id 类型错误", ctx)
		return
	}

	err = menuService.DeleteMenu(param, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (m *MenuApi) UpdateMenu(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.MenuParam
	var err error
	idStr := ctx.Param("menuId")
	param.MenuId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("menu_id 类型错误", ctx)
		return
	}

	req := request.MenuRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = menuService.UpdateMenu(param, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (m *MenuApi) GetMenu(ctx *gin.Context) {
	var param request.MenuParam
	var err error
	idStr := ctx.Param("menuId")
	param.MenuId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("menu_id 类型错误", ctx)
		return
	}

	data, err := menuService.GetMenu(param)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}

func (m *MenuApi) GetMenuList(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	data, err := menuService.GetMenuTreeList(common.ScopeData{Claims: claims})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}

func (m *MenuApi) GetMenuTree(ctx *gin.Context) {
	data, err := menuService.GetMenuTree()
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}

func (m *MenuApi) GetMenuTreeCompact(ctx *gin.Context) {
	data, err := menuService.GetMenuTreeCompact()
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}
