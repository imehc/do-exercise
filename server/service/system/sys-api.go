package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type SysApiService struct{}

// Update 更新api
func (s *SysApiService) Update(db *gorm.DB, id uint, req request.UpdateSysApiReq) error {
	// 先检查api是否存在
	existApi := &system.SysApi{}
	err := db.
		Where("id = ?", id).
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
func (s *SysApiService) Get(db *gorm.DB, id uint) (*response.SysApiResp, error) {
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
func (s *SysApiService) GetList(db *gorm.DB, req common.Pagination) (*common.PageResult[response.SysApiResp], error) {
	var apis []system.SysApi
	var total int64

	// Count 用独立 builder，避免污染后续 Find 的状态
	countDB := db.Model(&system.SysApi{})
	if err := countDB.Count(&total).Error; err != nil {
		return nil, errors.New("getApiListFailed")
	}
	req.Normalize()
	err := db.Model(&system.SysApi{}).
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Order("disabled ASC").
		Order("sort DESC").
		Order("id ASC").
		Find(&apis).Error
	if err != nil {
		return nil, errors.New("getApiListFailed")
	}

	return &common.PageResult[response.SysApiResp]{
		Data: lo.Map(apis, func(api system.SysApi, index int) response.SysApiResp {
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
		}),
		Meta: common.PageMeta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	}, nil
}

// GetAll 查询所有api
func (s *SysApiService) GetAll(db *gorm.DB) ([]response.SysApiResp, error) {
	var apis []system.SysApi
	err := db.
		Order("disabled ASC").
		Order("sort DESC").
		Order("id ASC").
		Find(&apis).
		Error
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
func (s *SysApiService) GroupType(db *gorm.DB) ([]string, error) {
	var groups []string
	err := db.
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
