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

type ApiApi struct{}

func (r *ApiApi) CreateApi(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.ApiRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := apiService.Create(req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (r *ApiApi) DeleteApi(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.ApiParam
	var err error
	idStr := ctx.Param("apiId")
	param.ApiId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("api_id 类型错误", ctx)
		return
	}

	err = apiService.Delete(param, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (r *ApiApi) UpdateApi(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var param request.ApiParam
	var err error
	idStr := ctx.Param("apiId")
	param.ApiId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("api_id 类型错误", ctx)
		return
	}

	req := request.ApiRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = apiService.Update(param, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (r *ApiApi) GetApi(ctx *gin.Context) {
	var param request.ApiParam
	var err error
	idStr := ctx.Param("apiId")
	param.ApiId, err = strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest("api_id 类型错误", ctx)
		return
	}

	data, err := apiService.Find(param)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}

func (r *ApiApi) GetApis(ctx *gin.Context) {
	data, err := apiService.FindAll()
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}

func (r *ApiApi) GetApiList(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	query := request.ApiQueryParams{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(utils.GetValidMsg(query, err), ctx)
		return
	}

	data, err := apiService.FindList(query, common.ScopeData{Claims: claims})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(data, ctx)
}
