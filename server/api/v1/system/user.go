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

// 创建用户
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

// 删除用户
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

// 更新用户
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

// 查询单个用户
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

// 查询用户列表
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

// 查询当前用户信息
func (u *UserApi) GetUserInfo(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)
	user, err := userService.GetUser(request.UserParam{UserId: int(claims.ID)})
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(user, ctx)
}

// 更新当前用户信息
func (u *UserApi) UpdateUserInfo(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.UserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := userService.UpdateUser(request.UserParam{UserId: int(claims.ID)}, req, claims.ID)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}
