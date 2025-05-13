package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
)

type SysInfoApi struct{}

// Get 系统信息
func (s *SysInfoApi) Get(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	info, err := sysInfoService.Get()
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	response.Success(ctx, info)
}
