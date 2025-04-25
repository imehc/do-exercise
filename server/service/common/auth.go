package common

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/request"
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
		return system.SysUser{}, errors.New("userNotFound")
	}
	hash := util.Hash{
		Value: existUser.Password,
	}
	if !hash.Compare(req.Password) {
		return system.SysUser{}, errors.New("passwordError")
	}

	return *existUser, nil
}

func (s *AuthService) ResetPassword(req request.UserResetPasswordReq) error {
	db := global.DB
	var user *system.SysUser
	result := db.
		Unscoped().
		First(&user, req.Id)
	if result.Error != nil {
		return errors.New("userNotFound")
	}

	if !user.DeletedAt.Time.IsZero() {
		return errors.New("userDeleted")
	}

	hash := util.Hash{Value: req.Password}
	password, err := hash.Hash()
	if err != nil {
		return err
	}
	user.Password = password

	if err := db.
		Model(user).
		Select("Password").
		Updates(user).
		Error; err != nil {
		return errors.New("resetPasswordFailed")
	}

	return nil
}
