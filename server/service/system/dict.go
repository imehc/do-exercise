package system

import (
	"errors"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
	"gorm.io/gorm"
)

type DictService struct{}

func (d *DictService) CreateDict(request request.CreateDictRequest, createdBy uint) (err error) {
	db := global.DB

	dict := system.Dict{
		Name:   request.Name,
		Type:   request.Type,
		Remark: request.Remark,
		ControlWrapper: model.ControlWrapper{
			CreatedBy: createdBy,
		},
	}
	dict.Status = request.Status

	err = db.
		Create(&dict).
		Error
	return
}

func (d *DictService) DeleteDict(param request.DictParam, deletedBy uint) (err error) {
	db := global.DB

	var dict system.Dict
	result := db.
		Unscoped().
		First(&dict, param.DictId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("字典不存在")
		}
		return result.Error
	}

	if !dict.DeletedAt.Time.IsZero() {
		return errors.New("字典已删除")
	}

	db.
		Model(system.Dict{}).
		Where("id = ?", param.DictId).
		Update("deleted_by", deletedBy).
		Delete(&dict)
	return nil
}

func (d *DictService) UpdateDict(param request.DictParam, request request.UpdateDictRequest, updatedBy uint) (err error) {
	db := global.DB

	var dict system.Dict
	result := db.
		First(&dict, param.DictId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("字典不存在")
		}
		return result.Error
	}

	dict.Status = request.Status
	dict.Remark = request.Remark
	dict.ControlWrapper = model.ControlWrapper{
		UpdatedBy: updatedBy,
	}

	db.
		Model(system.Dict{}).
		Where("id = ?", param.DictId).
		Updates(&dict).
		Select("status", "remark")

	return nil
}

func (d *DictService) GetDict(param request.DictParam) (response sysRes.DictItem, err error) {
	db := global.DB

	var dict system.Dict
	result := db.
		First(&dict, param.DictId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("字典不存在")
		}
		return response, result.Error
	}

	response.ID = dict.ID
	response.CreateDictRequest = request.CreateDictRequest{
		Name: dict.Name,
		Type: dict.Type,
		DictRequest: request.DictRequest{
			RemarkWrapper: model.RemarkWrapper{
				Remark: dict.Remark,
			},
			StatusWrapper: model.StatusWrapper{
				Status: dict.Status,
			},
		},
	}
	response.ControlWrapper = dict.ControlWrapper
	return
}

func (d *DictService) GetDictList(query request.DictQueryParams) (response sysRes.DictResponse, err error) {
	db := global.DB

	var total int64
	var originDicts []system.Dict
	offset := (query.Page - 1) * query.PageSize
	err = db.
		Model(&system.Dict{}).
		Order("id ASC").
		Where(fmt.Sprintf("name LIKE '%%%s%%'", query.Name)).
		Count(&total).
		Offset(offset).
		Limit(query.PageSize).
		Find(&originDicts).
		Error
	if err != nil {
		return response, err
	}
	response.Meta.Page = query.Page
	response.Meta.PageSize = query.PageSize
	response.Meta.Total = total

	response.Data = make([]sysRes.DictItem, len(originDicts))
	for i, dict := range originDicts {
		response.Data[i].ID = dict.ID
		response.Data[i].ControlWrapper = dict.ControlWrapper
		response.Data[i].Name = dict.Name
		response.Data[i].Type = dict.Type
		response.Data[i].Remark = dict.Remark
		response.Data[i].Status = dict.Status
	}

	return
}

// 字典数据

func (d *DictService) CreateDictData(request request.CreateDictDataRequest, createdBy uint) (err error) {
	db := global.DB

	var dict system.Dict
	result := db.
		Unscoped().
		Where("type = ?", request.DictType).
		First(&dict)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("字典不存在")
		}
		return result.Error
	}

	if !dict.DeletedAt.Time.IsZero() {
		return errors.New("字典已删除")
	}

	var dictData system.DictData
	result = db.
		Unscoped().
		Where("dict_type = ?", request.DictType).
		Where("value = ?", request.Value).
		First(&dictData)
	if result.Error == nil {
		return errors.New("数据键值已存在")
	}

	data := system.DictData{
		DictType: request.DictType,
		Label:    request.Label,
		Value:    request.Value,
		ControlWrapper: model.ControlWrapper{
			CreatedBy: createdBy,
		},
	}
	dictData.Status = request.Status
	dictData.Remark = request.Remark

	err = db.
		Create(&data).
		Error
	return
}

func (d *DictService) DeleteDictData(param request.DictDataParam, deletedBy uint) (err error) {
	db := global.DB

	var dictData system.DictData
	result := db.
		Unscoped().
		First(&dictData, param.DictDataId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("字典数据不存在")
		}
		return result.Error
	}

	if !dictData.DeletedAt.Time.IsZero() {
		return errors.New("字典数据已删除")
	}

	db.
		Model(system.DictData{}).
		Where("id = ?", param.DictDataId).
		Update("deleted_by", deletedBy).
		Delete(&dictData)
	return nil
}

func (d *DictService) UpdateDictData(param request.DictDataParam, request request.UpdateDictDataRequest, updatedBy uint) (err error) {
	db := global.DB

	var dictData system.DictData
	result := db.
		First(&dictData, param.DictDataId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("字典数据不存在")
		}
		return result.Error
	}

	dictData.Remark = request.Remark
	dictData.Sort = request.Sort
	dictData.Status = request.Status
	dictData.ControlWrapper = model.ControlWrapper{
		UpdatedBy: updatedBy,
	}

	db.
		Model(system.DictData{}).
		Where("id = ?", param.DictDataId).
		Updates(&dictData).
		Select("status", "remark", "sort")

	return nil
}

func (d *DictService) GetDictData(param request.DictDataParam) (response sysRes.DictDataItem, err error) {
	db := global.DB

	var dictData system.DictData
	result := db.
		First(&dictData, param.DictDataId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("字典数据不存在")
		}
		return response, result.Error
	}

	response.ID = dictData.ID
	response.CreateDictDataRequest = request.CreateDictDataRequest{
		DictTypeWrapper: request.DictTypeWrapper{
			DictType: dictData.DictType,
		},
		Label: dictData.Label,
		Value: dictData.Value,
	}
	response.Remark = dictData.Remark
	response.Sort = dictData.Sort
	response.Status = dictData.Status
	response.ControlWrapper = dictData.ControlWrapper
	return
}

func (d *DictService) GetDictDataList(query request.DictDataQueryParams) (response sysRes.DictDataResponse, err error) {
	db := global.DB

	var dict system.Dict
	result := db.
		Unscoped().
		Where("type = ?", query.DictType).
		First(&dict)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("字典不存在")
		}
		return response, result.Error
	}

	if !dict.DeletedAt.Time.IsZero() {
		return response, errors.New("字典已删除")
	}

	var total int64
	var originDictData []system.DictData
	offset := (query.Page - 1) * query.PageSize
	err = db.
		Model(&system.DictData{}).
		Order("sort Desc").
		Order("id ASC").
		Where(fmt.Sprintf("label LIKE '%%%s%%'", query.Label)).
		Where("dict_type = ?", query.DictType).
		Count(&total).Offset(offset).
		Limit(query.PageSize).
		Find(&originDictData).
		Error
	if err != nil {
		return response, err
	}
	response.Meta.Page = query.Page
	response.Meta.PageSize = query.PageSize
	response.Meta.Total = total

	response.Data = make([]sysRes.DictDataItem, len(originDictData))
	for i, dict := range originDictData {
		response.Data[i].ID = dict.ID
		response.Data[i].ControlWrapper = dict.ControlWrapper
		response.Data[i].DictType = dict.DictType
		response.Data[i].Label = dict.Label
		response.Data[i].Value = dict.Value
		response.Data[i].Remark = dict.Remark
		response.Data[i].Sort = dict.Sort
		response.Data[i].Status = dict.Status
	}

	return
}
