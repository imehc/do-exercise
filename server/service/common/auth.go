package common

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
)

type AuthService struct{}

func (s *AuthService) Login(req common.Login) (system.SysUser, error) {
	existUser := &system.SysUser{}
	err := global.DB.Where("username = ?", req.Username).
		First(existUser).
		Error
	if err != nil {
		return system.SysUser{}, err
	}
	return *existUser, nil
}
