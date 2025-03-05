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

type DictApi struct{}

func (d *DictApi) CreateDict(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.CreateDictRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := dictService.CreateDict(req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (d *DictApi) DeleteDict(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var dictParam request.DictParam
	var err error
	dicIdStr := ctx.Param("dictId")
	dictParam.DictId, err = strconv.Atoi(dicIdStr)
	if err != nil {
		response.BadRequest("dict_id 类型错误", ctx)
		return
	}

	err = dictService.DeleteDict(dictParam, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (d *DictApi) UpdateDict(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var dictParam request.DictParam
	var err error
	dicIdStr := ctx.Param("dictId")
	dictParam.DictId, err = strconv.Atoi(dicIdStr)
	if err != nil {
		response.BadRequest("dict_id 类型错误", ctx)
		return
	}

	req := request.UpdateDictRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = dictService.UpdateDict(dictParam, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (d *DictApi) GetDict(ctx *gin.Context) {
	var dictParam request.DictParam
	var err error
	dicIdStr := ctx.Param("dictId")
	dictParam.DictId, err = strconv.Atoi(dicIdStr)
	if err != nil {
		response.BadRequest("dict_id 类型错误", ctx)
		return
	}

	dict, err := dictService.GetDict(dictParam)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(dict, ctx)
}

func (d *DictApi) GetDictList(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	query := request.DictQueryParams{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(utils.GetValidMsg(query, err), ctx)
		return
	}

	dicts, err := dictService.GetDictList(query, common.ScopeData{Claims: claims})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(dicts, ctx)
}

// 字典详情

func (d *DictApi) CreateDictData(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.CreateDictDataRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := dictService.CreateDictData(req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (d *DictApi) DeleteDictData(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var dictDataParam request.DictDataParam
	var err error
	dictDataIdStr := ctx.Param("dictDataId")
	dictDataParam.DictDataId, err = strconv.Atoi(dictDataIdStr)
	if err != nil {
		response.BadRequest("dict_data_id 类型错误", ctx)
		return
	}

	err = dictService.DeleteDictData(dictDataParam, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (d *DictApi) UpdateDictData(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var dictDataParam request.DictDataParam
	var err error
	dictDataIdStr := ctx.Param("dictDataId")
	dictDataParam.DictDataId, err = strconv.Atoi(dictDataIdStr)
	if err != nil {
		response.BadRequest("dict_data_id 类型错误", ctx)
		return
	}

	req := request.UpdateDictDataRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = dictService.UpdateDictData(dictDataParam, req, claims.UserId)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (d *DictApi) GetDictData(ctx *gin.Context) {
	var dictDataParam request.DictDataParam
	var err error
	dictDataIdStr := ctx.Param("dictDataId")
	dictDataParam.DictDataId, err = strconv.Atoi(dictDataIdStr)
	if err != nil {
		response.BadRequest("dict_data_id 类型错误", ctx)
		return
	}

	dictData, err := dictService.GetDictData(dictDataParam)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(dictData, ctx)
}

func (d *DictApi) GetDictDataList(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	query := request.DictDataQueryParams{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(utils.GetValidMsg(query, err), ctx)
		return
	}

	dictData, err := dictService.GetDictDataList(query, common.ScopeData{Claims: claims})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(dictData, ctx)
}
