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

type DeptApi struct{}

func (d *DeptApi) CreateDept(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.DeptRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := deptService.CreateDept(req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (d *DeptApi) DeleteDept(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var deptParam request.DeptParam
	var err error
	deptIdStr := ctx.Param("deptId")
	deptParam.DeptId, err = strconv.Atoi(deptIdStr)
	if err != nil {
		response.BadRequest("dept_id 类型错误", ctx)
		return
	}

	err = deptService.DeleteDept(deptParam, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (d *DeptApi) UpdateDept(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var deptParam request.DeptParam
	var err error
	deptIdStr := ctx.Param("deptId")
	deptParam.DeptId, err = strconv.Atoi(deptIdStr)
	if err != nil {
		response.BadRequest("dept_id 类型错误", ctx)
		return
	}

	req := request.DeptRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = deptService.UpdateDept(deptParam, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (d *DeptApi) GetDept(ctx *gin.Context) {
	var deptParam request.DeptParam
	var err error
	deptIdStr := ctx.Param("deptId")
	deptParam.DeptId, err = strconv.Atoi(deptIdStr)
	if err != nil {
		response.BadRequest("dept_id 类型错误", ctx)
		return
	}

	dept, err := deptService.GetDep(deptParam)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(dept, ctx)
}

func (d *DeptApi) GetDeptList(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	depts, err := deptService.GetDeptTreeList(common.ScopeData{Claims: claims})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(depts, ctx)
}

func (d *DeptApi) GetDeptTree(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	depts, err := deptService.GetDeptTree(&common.ScopeData{Claims: claims})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(depts, ctx)
}
