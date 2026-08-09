package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct{}

func (s *AuthService) Login(db *gorm.DB, req common.Login) (*system.SysUser, error) {
	existUser := &system.SysUser{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.WithContext(ctx).
		Preload("Roles").
		Where("username = ?", req.Username).
		First(existUser).
		Error
	if err != nil {
		global.Log.Error("登录失败", zap.Error(err), zap.String("username", req.Username))
		errorKey := util.TranslateDBError(err, "userNotFound")
		return nil, errors.New(errorKey)
	}
	hash := util.Hash{
		Value: existUser.Password,
	}
	if !hash.Compare(req.Password) {
		return nil, errors.New("passwordError")
	}

	return existUser, nil
}

func (s *AuthService) ResetPassword(db *gorm.DB, req request.UserResetPasswordReq) error {
	hash := util.Hash{Value: req.Password}
	password, err := hash.Hash()
	if err != nil {
		return err
	}

	// user.Password=pa
	err = db.Model(&system.SysUser{}).Where("id = ?", req.Id).Update("Password", password).Error
	if err != nil {
		global.Log.Error("reset password failed", zap.String("userId", req.Id), zap.Error(err))
		return err
	}

	return nil
}

func (s *AuthService) Logout(userId string, accessToken string) error {
	redis := global.Redis
	log := global.Log
	accesskey := fmt.Sprintf("%s%s", util.PrefixAccessToken, accessToken)
	ctx := context.Background()
	tokenInfo, err := redis.Get(ctx, accesskey).Result()
	if err != nil {
		log.Error("Failed to get token info from Redis", zap.String("userId", userId), zap.Error(err))

		return errors.New("operationFailed")
	}

	var tokenData model.TokenInfo
	if err := json.Unmarshal([]byte(tokenInfo), &tokenData); err != nil {
		log.Error("Failed to unmarshal token info", zap.String("userId", userId), zap.Error(err))
		return errors.New("operationFailed")
	}

	refreshKey := fmt.Sprintf("%s%s", util.PrefixRefreshToken, tokenData.RefreshToken)

	pipe := redis.Pipeline()
	if err := pipe.Del(ctx, accesskey).Err(); err != nil {
		log.Error("Failed to delete access token", zap.String("userId", userId), zap.Error(err))
		return errors.New("operationFailed")
	}
	if err := pipe.Del(ctx, refreshKey).Err(); err != nil {
		log.Error("Failed to delete refresh token", zap.String("userId", userId), zap.Error(err))
		return errors.New("operationFailed")
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("Failed to execute pipeline", zap.Error(err))
		return errors.New("operationFailed")
	}

	// 清理该用户的无效 access/refresh token
	err = util.CleanUserTokenSet(tokenData.UserId, util.PrefixUserAcessToken, util.PrefixAccessToken)
	if err != nil {
		log.Error("Failed to clean user accessToken set", zap.String("userId", userId), zap.Error(err))
		return errors.New("operationFailed")
	}
	err = util.CleanUserTokenSet(tokenData.UserId, util.PrefixUserRefreshToken, util.PrefixRefreshToken)
	if err != nil {
		log.Error("Failed to clean user refreshToken set", zap.String("userId", userId), zap.Error(err))
		return errors.New("operationFailed")
	}

	return nil
}
