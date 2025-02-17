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

type PostApi struct{}

func (p *PostApi) CreatePost(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	req := request.PostRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err := postService.CreatePost(req, claims.ID)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Created(ctx)
}

func (p *PostApi) DeletePost(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var postParam request.PostParam
	var err error
	postIdStr := ctx.Param("postId")
	postParam.PostId, err = strconv.Atoi(postIdStr)
	if err != nil {
		response.BadRequest("post_id 类型错误", ctx)
		return
	}

	err = postService.DeletePost(postParam, claims.ID)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (p *PostApi) UpdatePost(ctx *gin.Context) {
	data, _ := ctx.Get(global.CLAIMS)
	claims := data.(system.Claims)

	var postParam request.PostParam
	var err error
	postIdStr := ctx.Param("postId")
	postParam.PostId, err = strconv.Atoi(postIdStr)
	if err != nil {
		response.BadRequest("post_id 类型错误", ctx)
		return
	}

	req := request.PostRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	err = postService.UpdatePost(postParam, req, claims.ID)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.NoContent(ctx)
}

func (p *PostApi) GetPost(ctx *gin.Context) {
	var postParam request.PostParam
	var err error
	postIdStr := ctx.Param("postId")
	postParam.PostId, err = strconv.Atoi(postIdStr)
	if err != nil {
		response.BadRequest("post_id 类型错误", ctx)
		return
	}

	post, err := postService.GetPost(postParam)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(post, ctx)
}

func (p *PostApi) GetPostList(ctx *gin.Context) {
	query := request.PostQueryParams{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(utils.GetValidMsg(query, err), ctx)
		return
	}

	posts, err := postService.GetPostList(query)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}

	response.Success(posts, ctx)
}
