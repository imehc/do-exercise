package system

import (
	"net/http"

	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/labstack/echo/v4"
)

type SysUserApi struct{}

func (s *SysUserApi) Create(ctx echo.Context) error {
	var sys request.CreateSysUserReq
	if err := ctx.Bind(&sys); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "请求参数绑定失败",
			"detail": err.Error(),
		})
	}

	if err := ctx.Validate(&sys); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "参数验证失败",
			"detail": err.Error(),
			"lang":   ctx.Get("lang"),
		})
	}

	return ctx.JSON(http.StatusOK, sys)
}
