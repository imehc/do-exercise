package system

import "github.com/imehc/do-exercise/server/model/common"

type SysUserService struct{}

func (s *SysUserService) Create() {

}

func (s *SysUserService) Login(common.Login) (bool, error) {
	return true, nil
}
