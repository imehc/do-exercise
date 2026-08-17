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
	commonResponse "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
	systemService "github.com/imehc/do-exercise/server/service/system"
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct{}

// Login 登录。返回命中的用户（带角色）及其归属的启用业务租户列表。
// 多租户模式下未指定租户且账号归属多个启用租户时返回 requiresTenantSelection，
// 由前端弹窗选择后经 select_tenant 完成登录。
func (s *AuthService) Login(db *gorm.DB, req common.Login) (*system.SysUser, []commonResponse.TenantOption, error) {
	tenantId := req.TenantId
	if !global.Config.Tenant.IsMulti() {
		// 单租户模式始终以默认租户为准，忽略客户端传入的租户ID
		tenantId = global.Config.Tenant.DefaultTenantId
	} else if tenantId == "" {
		return s.loginAcrossTenants(db, req)
	}

	if err := s.checkTenantUsable(db, tenantId); err != nil {
		return nil, nil, err
	}
	user, err := s.Authenticate(db, tenantId, req.Username, req.Password)
	if err != nil {
		return nil, nil, err
	}
	options, _ := s.listEnabledTenants(db, req.Username)
	return user, options, nil
}

// loginAcrossTenants 多租户模式未指定租户时的登录解析：
// 校验用户名密码后按归属的启用租户分发——唯一归属直接进入；
// 多个启用租户时返回 requiresTenantSelection 引导前端弹窗选择；
func (s *AuthService) loginAcrossTenants(db *gorm.DB, req common.Login) (*system.SysUser, []commonResponse.TenantOption, error) {
	var rows []system.SysUser
	if err := db.Preload("Roles").Where("username = ?", req.Username).Find(&rows).Error; err != nil {
		return nil, nil, errors.New("userNotFound")
	}
	if len(rows) == 0 {
		return nil, nil, errors.New("userNotFound")
	}

	var (
		business []system.SysUser
		platform []system.SysUser
		matched  bool
	)
	for _, u := range rows {
		hash := util.Hash{Value: u.Password}
		if !hash.Compare(req.Password) {
			continue
		}
		matched = true
		if u.TenantId == global.PlatformTenantID {
			platform = append(platform, u)
			continue
		}
		var tenant system.SysTenant
		if err := db.Where("tenant_id = ? AND status = ?", u.TenantId, true).First(&tenant).Error; err != nil {
			continue
		}
		business = append(business, u)
	}

	if len(business) == 0 && len(platform) == 0 {
		if !matched {
			return nil, nil, errors.New("passwordError")
		}
		return nil, nil, errors.New("tenantDisabled")
	}

	options, _ := s.listEnabledTenants(db, req.Username)

	switch len(business) {
	case 0:
		// 仅平台归属：平台管理员 直接登录平台域，不出现在租户选择器中
		return &platform[0], options, nil
	case 1:
		return &business[0], options, nil
	default:
		return nil, options, errors.New("requiresTenantSelection")
	}
}

// listEnabledTenants 返回用户名归属的启用业务租户（不含平台保留租户）
func (s *AuthService) listEnabledTenants(db *gorm.DB, username string) ([]commonResponse.TenantOption, error) {
	return (&systemService.SysTenantService{}).ListEnabledForUsername(db, username)
}

// checkTenantUsable 校验指定租户存在且启用；单租户模式或平台保留租户始终可用。
func (s *AuthService) checkTenantUsable(db *gorm.DB, tenantId string) error {
	if !global.Config.Tenant.IsMulti() {
		return nil
	}
	if tenantId == "" || tenantId == global.PlatformTenantID {
		return nil
	}
	var tenant system.SysTenant
	if err := db.Unscoped().Where("tenant_id = ?", tenantId).First(&tenant).Error; err != nil {
		return errors.New("tenantNotFound")
	}
	if !tenant.Status {
		return errors.New("tenantDisabled")
	}
	return nil
}

// Authenticate 校验指定租户下的用户名密码，返回带角色信息的用户。
func (s *AuthService) Authenticate(db *gorm.DB, tenantId, username, password string) (*system.SysUser, error) {
	existUser := &system.SysUser{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.WithContext(ctx).
		Preload("Roles").
		Where("tenant_id = ?", tenantId).
		Where("username = ?", username).
		First(existUser).
		Error
	if err != nil {
		global.Log.Error("登录失败", zap.Error(err), zap.String("username", username))
		errorKey := util.TranslateDBError(err, "userNotFound")
		return nil, errors.New(errorKey)
	}
	hash := util.Hash{Value: existUser.Password}
	if !hash.Compare(password) {
		return nil, errors.New("passwordError")
	}
	return existUser, nil
}

// EnterTenant 校验租户可用并加载该租户下指定用户名的账号（含角色）。
// 用于多租户登录选择/切换：密码已在前一阶段验证，这里只做归属与可用性校验。
func (s *AuthService) EnterTenant(db *gorm.DB, tenantId, username string) (*system.SysUser, error) {
	if err := s.checkTenantUsable(db, tenantId); err != nil {
		return nil, err
	}
	user := &system.SysUser{}
	err := db.
		Preload("Roles").
		Where("tenant_id = ?", tenantId).
		Where("username = ?", username).
		First(user).
		Error
	if err != nil {
		errorKey := util.TranslateDBError(err, "userNotFound")
		return nil, errors.New(errorKey)
	}
	return user, nil
}

// ResetPassword 重置密码。id 必须由调用方从已校验的来源（邮箱查找、管理员路径参数）取得，
// 请求结构体不再承载 id 字段，杜绝客户端可控 id 被误用的可能。
func (s *AuthService) ResetPassword(db *gorm.DB, id string, req request.UserResetPasswordReq) error {
	hash := util.Hash{Value: req.Password}
	password, err := hash.Hash()
	if err != nil {
		return err
	}

	// user.Password=pa
	err = db.Model(&system.SysUser{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Password":           password,
		"MustChangePassword": false,
	}).Error
	if err != nil {
		global.Log.Error("reset password failed", zap.String("userId", id), zap.Error(err))
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
