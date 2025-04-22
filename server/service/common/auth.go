package common

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
)

type AuthService struct{}

func (s *AuthService) Login(req common.Login) (system.SysUser, error) {
	existUser := &system.SysUser{}
	err := global.DB.
		Preload("Roles").
		Where("username = ?", req.Username).
		First(existUser).
		Error
	if err != nil {
		return system.SysUser{}, err
	}
	hash := util.Hash{
		Value: existUser.Password,
	}
	if !hash.Compare(req.Password) {
		return system.SysUser{}, errors.New("密码错误")
	}

	return *existUser, nil
}
