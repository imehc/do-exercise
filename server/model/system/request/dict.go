package request

import (
	"github.com/imehc/do-exercise/server/model"
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

// 通用信息，不需要验证
type DictRequest struct {
	model.RemarkWrapper
	model.StatusWrapper
	model.ControlWrapper
}

// 创建
type CreateDictRequest struct {
	Name        string `json:"name" binding:"required"` // 字典名称
	Type        string `json:"type" binding:"required"` // 字典类型
	DictRequest `json:",inline"`
}

func (d CreateDictRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Name.required": "字典名称不能为空",
		"Type.required": "字典类型不能为空",
	}
}

// 更新
type UpdateDictRequest struct {
	DictRequest
}

func (d UpdateDictRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type DictParam struct {
	DictId int `json:"dict_id" binding:"required"`
}

type DictQueryParams struct {
	commonReq.QueryParams `json:",inline"`
	Name                  string `json:"name" form:"name"`
}

type DictTypeWrapper struct {
	DictType string `json:"dict_type" form:"dict_type" binding:"required"` // 字典类型
}

func (d DictTypeWrapper) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"DictType.required": "字典类型不能为空",
	}
}

// 字典数据
type DictDataRequest struct {
	model.RemarkWrapper
	model.SortWrapper
	model.StatusWrapper
}

type CreateDictDataRequest struct {
	DictTypeWrapper `json:",inline"`
	Label           string `json:"label" binding:"required"` // 数据标签
	Value           string `json:"value" binding:"required"` // 数据键值
	DictDataRequest `json:",inline"`
}

func (d CreateDictDataRequest) GetMessage() utils.ValidatorMessages {
	messages := utils.ValidatorMessages{
		"Label.required": "数据标签不能为空",
		"Value.required": "数据键值不能为空",
	}

	if dtm := d.DictTypeWrapper.GetMessage(); dtm != nil {
		for k, v := range dtm {
			messages[k] = v
		}
	}

	return messages
}

type UpdateDictDataRequest struct {
	DictDataRequest
}

func (d UpdateDictDataRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type DictDataParam struct {
	DictDataId int `json:"dict_data_id" binding:"required"`
}

type DictDataQueryParams struct {
	commonReq.QueryParams `json:",inline"`
	DictTypeWrapper       `json:",inline"`
	Label                 string `json:"label" form:"label"`
}

// TIP: 如果结构体内部都实现了同一个接口，则需要指明分别调用那个一个接口

func (d DictDataQueryParams) GetMessage() utils.ValidatorMessages {
	// 初始化一个空的 ValidatorMessages
	messages := make(utils.ValidatorMessages)
	// 获取 QueryParams 的消息并合并
	if qm := d.QueryParams.GetMessage(); qm != nil {
		for k, v := range qm {
			messages[k] = v
		}
	}

	// 获取 DictTypeWrapper 的消息并合并
	if dtm := d.DictTypeWrapper.GetMessage(); dtm != nil {
		for k, v := range dtm {
			messages[k] = v
		}
	}

	return messages
}
