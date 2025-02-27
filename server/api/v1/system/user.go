package system

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/utils"
)

type UserApi struct{}

func (u *UserApi) CreateUser(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.UserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := userService.CreateUser(req, claims.ID)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (u *UserApi) DeleteUser(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var userParam request.UserParam
	var err error
	userIdStr := ctx.Param("userId")
	userParam.UserId, err = strconv.Atoi(userIdStr)
	if err != nil {
		response.BadRequest("user_id 类型错误", ctx)
		return
	}

	err = userService.DeleteUser(userParam, claims.ID)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (u *UserApi) UpdateUser(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var userParam request.UserParam
	var err error
	userIdStr := ctx.Param("userId")
	userParam.UserId, err = strconv.Atoi(userIdStr)
	if err != nil {
		response.BadRequest("user_id 类型错误", ctx)
		return
	}

	req := request.UserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = userService.UpdateUser(userParam, req, claims.ID)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (u *UserApi) GetUser(ctx *gin.Context) {
	var userParam request.UserParam
	var err error
	userIdStr := ctx.Param("userId")
	userParam.UserId, err = strconv.Atoi(userIdStr)
	if err != nil {
		response.BadRequest("user_id 类型错误", ctx)
		return
	}

	user, err := userService.GetUser(userParam)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(user, ctx)
}

func (u *UserApi) GetUserList(ctx *gin.Context) {
	query := request.UserQueryParams{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(utils.GetValidMsg(query, err), ctx)
		return
	}

	users, err := userService.GetUserList(query)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(users, ctx)
}
