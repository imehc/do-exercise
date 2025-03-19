package router

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	m "github.com/imehc/do-exercise/server/middleware"
	"github.com/imehc/do-exercise/server/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServer() {
	e := echo.New()

	e.Use(m.Logger)
	e.Use(m.I18nMiddleware())
	e.Use(middleware.Recover())

	e.Validator = &util.CustomValidator{}

	g := e.Group("")

	RouterGroupApp.Normal.InitNormalRouter(g)
	RouterGroupApp.System.InitSysUserRouter(g)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", global.Config.System.Port)))
}
