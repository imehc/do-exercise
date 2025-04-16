package common

import "github.com/imehc/do-exercise/server/model/common"

type AuthService struct{}

func (s *AuthService) Login(common.Login) (uint, error) {
	return 111, nil
}
