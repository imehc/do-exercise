package system

import (
	"github.com/imehc/do-exercise/server/model/system/request"
)

type AuthService struct{}

func (a *AuthService) Login(req request.Login) (enable bool, err error) {
	enable = true
	err = nil
	return
}

func (a *AuthService) Logout() error {
	// 模拟成功
	return nil
}
