package system

import (
	"errors"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/utils"
	"gorm.io/gorm"
)

type ApiService struct{}

// 创建api
func (r ApiService) Create(request request.ApiRequest, createdBy uint) (err error) {
	db := global.DB

	api := system.Api{
		Handle: request.Handle,
		Title:  request.Title,
		Path:   request.Path,
		Type:   request.Type,
		Action: request.Action,
		ControlWrapper: model.ControlWrapper{
			CreatedBy: createdBy,
		},
	}

	err = db.Create(&api).Error
	return
}

// 删除api
func (r ApiService) Delete(param request.ApiParam, deletedBy uint) (err error) {
	db := global.DB

	var api system.Api
	result := db.
		Unscoped().
		First(&api, param.ApiId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("接口不存在")
		}
		return result.Error
	}

	if !api.DeletedAt.Time.IsZero() {
		return errors.New("接口已删除")
	}

	db.
		Model(system.Api{}).
		Where("id = ?", param.ApiId).
		Update("deleted_by", deletedBy).
		Delete(&api)
	return nil
}

// 更新api
func (r ApiService) Update(param request.ApiParam, request request.ApiRequest, updatedBy uint) (err error) {
	db := global.DB

	var api system.Api
	result := db.
		First(&api, param.ApiId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("接口不存在")
		}
		return result.Error
	}

	api.Handle = request.Handle
	api.Title = request.Title
	api.Path = request.Path
	api.Type = request.Type
	api.Action = request.Action
	api.ControlWrapper = model.ControlWrapper{
		UpdatedBy: updatedBy,
	}

	db.
		Model(system.Api{}).
		Where("id = ?", param.ApiId).
		Updates(&api).
		Omit("id", "created_at")

	return nil
}

// 查询api
func (r ApiService) Find(param request.ApiParam) (response sysRes.ApiItem, err error) {
	db := global.DB

	var api system.Api
	result := db.
		First(&api, param.ApiId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("接口不存在")
		}
		return response, result.Error
	}

	response.ID = api.ID
	response.ControlWrapper = api.ControlWrapper
	response.Handle = api.Handle
	response.Title = api.Title
	response.Path = api.Path
	response.Type = api.Type
	response.Action = api.Action
	return
}

// 查询所有api
func (r ApiService) FindAll() (response []sysRes.ApiItem, err error) {
	db := global.DB

	var apis []system.Api
	err = db.
		Order("id ASC").
		Find(&apis).
		Error
	if err != nil {
		return response, err
	}

	response = make([]sysRes.ApiItem, len(apis))
	for i, api := range apis {
		response[i].ID = api.ID
		response[i].ControlWrapper = api.ControlWrapper
		response[i].Handle = api.Handle
		response[i].Title = api.Title
		response[i].Path = api.Path
		response[i].Type = api.Type
		response[i].Action = api.Action
	}

	return
}

// 查询api列表
func (r ApiService) FindList(query request.ApiQueryParams) (response sysRes.ApiResponse, err error) {
	db := global.DB

	var total int64
	var originApis []system.Api
	err = db.
		Model(&system.Api{}).
		Order("id ASC").
		Where(fmt.Sprintf("title LIKE '%%%s%%'", query.Name)).
		Count(&total).
		Scopes(utils.Paginate(query.PageSize, query.Page)).
		Find(&originApis).
		Error
	if err != nil {
		return response, err
	}

	response.Meta.Page = query.Page
	response.Meta.PageSize = query.PageSize
	response.Meta.Total = total

	response.Data = make([]sysRes.ApiItem, len(originApis))
	for i, api := range originApis {
		response.Data[i].ID = api.ID
		response.Data[i].ControlWrapper = api.ControlWrapper
		response.Data[i].Handle = api.Handle
		response.Data[i].Title = api.Title
		response.Data[i].Path = api.Path
		response.Data[i].Type = api.Type
		response.Data[i].Action = api.Action
	}

	return
}
