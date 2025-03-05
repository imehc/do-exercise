package response

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type DictItem struct {
	common.IDWrapper
	model.ControlWrapper
	request.CreateDictRequest
}

type DictResponse struct {
	Data []DictItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}

type DictDataItem struct {
	common.IDWrapper
	model.ControlWrapper
	request.CreateDictDataRequest
}

type DictDataResponse struct {
	Data []DictDataItem     `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
