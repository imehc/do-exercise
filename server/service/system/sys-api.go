package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
)

type SysApiService struct{}

// Update 更新api
func (s *SysApiService) Update(req request.UpdateSysApiReq) error {
	db := global.DB
	// 先检查api是否存在
	existApi := &system.SysApi{}
	err := db.
		Where("id = ?", req.Id).
		First(existApi).
		Error
	if err != nil {
		return errors.New("allApisNotFound")
	}
	existApi.Description = req.Description
	existApi.Group = req.Group
	existApi.Disabled = req.Disabled
	existApi.Sort = req.Sort

	err = db.
		Model(existApi).
		Select("Description", "Group", "Disabled", "Sort").
		Updates(existApi).
		Error
	if err != nil {
		return errors.New("updateApiFailed")
	}
	return nil
}

// Get 查询单个api
func (s *SysApiService) Get(id uint) (*response.SysApiResp, error) {
	db := global.DB
	// 先检查api是否存在
	existApi := &system.SysApi{}
	err := db.
		Where("id = ?", id).
		First(existApi).
		Error
	if err != nil {
		return nil, errors.New("allApisNotFound")
	}

	return &response.SysApiResp{
		Id:          existApi.Id,
		Path:        existApi.Path,
		Description: existApi.Description,
		Group:       existApi.Group,
		Disabled:    existApi.Disabled,
		Sort:        existApi.Sort,
		CreatedAt:   existApi.CreatedAt,
		UpdatedAt:   existApi.UpdatedAt,
	}, nil
}

// GetList 查询api列表
func (s *SysApiService) GetList(req common.Pagination) (*common.PageResult[response.SysApiResp], error) {
	var apis []system.SysApi
	var total int64
	db := global.DB.Model(&system.SysApi{})
	db.Count(&total)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	db = db.Scopes(util.Paginate(req.PageSize, req.Page))
	err := db.Find(&apis).Error
	if err != nil {
		return nil, errors.New("getApiListFailed")
	}

	return &common.PageResult[response.SysApiResp]{
		Data: lo.Map(apis, func(api system.SysApi, index int) response.SysApiResp {
			return response.SysApiResp{
				Id:          api.Id,
				Method:      api.Method,
				Description: api.Description,
				Group:       api.Group,
				Disabled:    api.Disabled,
				Sort:        api.Sort,
				CreatedAt:   api.CreatedAt,
				UpdatedAt:   api.UpdatedAt,
			}
		}),
		Meta: common.PageMeta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	}, nil

}

// GetAll 查询所有api
func (s *SysApiService) GetAll() ([]response.SysApiResp, error) {
	var apis []system.SysApi
	err := global.DB.Find(&apis).Error
	if err != nil {
		return nil, errors.New("getAllApiListFailed")
	}
	return lo.Map(apis, func(api system.SysApi, index int) response.SysApiResp {
		return response.SysApiResp{
			Id:          api.Id,
			Path:        api.Path,
			Method:      api.Method,
			Description: api.Description,
			Group:       api.Group,
			Disabled:    api.Disabled,
			Sort:        api.Sort,
			CreatedAt:   api.CreatedAt,
			UpdatedAt:   api.UpdatedAt,
		}
	}), nil
}

// GroupType 查询api分组类型
func (s *SysApiService) GroupType() ([]string, error) {
	var groups []string
	err := global.DB.
		Model(&system.SysApi{}).
		Select("COALESCE(\"group\", 'Default') as group").
		Distinct().
		Pluck("group", &groups).
		Error
	if err != nil {
		return nil, errors.New("getApiGroupTypeFailed")
	}
	return groups, nil
}
