package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/utils"
)

type AuthService struct{}

func (a *AuthService) Login(req request.Login) (enable bool, err error) {
	db := global.DB

	var user system.User
	err = db.Model(system.User{}).Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		return false, errors.New("用户不存在")
	}

	if !utils.BcryptCheck(req.Password, user.Password) {
		return false, errors.New("密码错误")
	}

	// TODO: 使用字典数据
	if user.Status == 0 {
		return true, nil
	}

	return true, nil
}

func (a *AuthService) Logout() error {
	// 模拟成功
	return nil
}
