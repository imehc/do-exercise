package system

import "github.com/imehc/do-exercise/server/model/common"

type SysUserService struct{}

type userId = uint

func (s *SysUserService) Create() {

}

func (s *SysUserService) Login(common.Login) (userId, error) {
	return 111, nil
}
