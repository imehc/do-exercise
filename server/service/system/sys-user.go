package system

import "github.com/imehc/do-exercise/server/model/system/request"

type SysUserService struct{}

func (s *SysUserService) Create() {

}

func (s *SysUserService) Login(request.Login) (bool, error) {
	return true, nil
}
