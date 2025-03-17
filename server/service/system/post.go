package system

import (
	"errors"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/pkg/utils/scope"
	"github.com/imehc/do-exercise/server/utils"
	"gorm.io/gorm"
)

type PostService struct{}

// 创建岗位
func (p *PostService) CreatePost(request request.PostRequest, createBy uint) (err error) {
	db := global.DB

	if !errors.Is(global.DB.Where("code = ?", request.Code).First(&system.Post{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同岗位编码")
	}

	post := system.Post{
		Name: request.Name,
		Code: request.Code,
		ControlWrapper: model.ControlWrapper{
			CreateBy: createBy,
		},
	}

	if request.Sort != 0 {
		post.Sort = request.Sort
	}
	if request.Status != 0 {
		post.Status = request.Status
	}
	if request.Remark != "" {
		post.Remark = request.Remark
	}

	return db.
		Create(&post).
		Error
}

// 删除岗位
func (p *PostService) DeletePost(param request.PostParam, deleteBy uint) (err error) {
	db := global.DB

	var post system.Post
	result := db.
		Unscoped().
		First(&post, param.PostId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("岗位不存在")
		}
		return result.Error
	}

	if !post.DeleteAt.Time.IsZero() {
		return errors.New("岗位已删除")
	}

	// 检查该岗位下是否有其它用户
	var user system.User
	result = db.
		Model(system.User{}).
		Where("post_id = ?", param.PostId).
		Find(&user)
	if result.RowsAffected > 0 {
		return errors.New("该岗位下存在用户，无法删除")
	}

	db.
		Model(system.Post{}).
		Where("post_id = ?", param.PostId).
		Update("delete_by", deleteBy).
		Delete(&post)
	return nil
}

// 更新岗位
func (p *PostService) UpdatePost(param request.PostParam, request request.PostRequest, updateBy uint) (err error) {
	db := global.DB

	var post system.Post
	result := db.
		First(&post, param.PostId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("岗位不存在")
		}
		return result.Error
	}

	post.Name = request.Name
	post.Code = request.Code
	post.Sort = request.Sort
	post.Status = request.Status
	post.Remark = request.Remark
	post.ControlWrapper = model.ControlWrapper{
		UpdateBy: updateBy,
	}

	db.
		Model(system.Post{}).
		Where("post_id = ?", param.PostId).
		Updates(&post).
		Omit("post_id", "created_at")

	return nil
}

// 根据id获取岗位信息
func (p *PostService) GetPost(param request.PostParam) (response sysRes.PostItem, err error) {
	db := global.DB

	var post system.Post
	result := db.
		First(&post, param.PostId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("岗位不存在")
		}
		return response, result.Error
	}

	response.ID = post.PostId
	response.PostRequest = request.PostRequest{
		Name:          post.Name,
		Code:          post.Code,
		SortWrapper:   post.SortWrapper,
		StatusWrapper: post.StatusWrapper,
		RemarkWrapper: post.RemarkWrapper,
	}

	response.ControlWrapper = post.ControlWrapper
	return
}

// 获取岗位列表
func (p *PostService) GetPostList(query request.PostQueryParams, s common.ScopeData) (response sysRes.PostResponse, err error) {
	db := global.DB
	// 应用数据权限过滤
	db = scope.GetDataScope(db, &s, "sys_post")

	var total int64
	var originPosts []system.Post
	err = db.
		Model(&system.Post{}).
		Order("sort Desc").
		Order("post_id ASC").
		Where(fmt.Sprintf("name LIKE '%%%s%%'", query.Name)).
		Count(&total).
		Scopes(utils.Paginate(query.PageSize, query.Page)).
		Find(&originPosts).
		Error
	if err != nil {
		return response, err
	}

	response.Meta.Page = query.Page
	response.Meta.PageSize = query.PageSize
	response.Meta.Total = total

	response.Data = make([]sysRes.PostItem, len(originPosts))
	for i, post := range originPosts {
		response.Data[i].ID = post.PostId
		response.Data[i].ControlWrapper = post.ControlWrapper
		response.Data[i].Name = post.Name
		response.Data[i].Code = post.Code
		response.Data[i].Sort = post.Sort
		response.Data[i].Status = post.Status
		response.Data[i].Remark = post.Remark
	}

	return
}
