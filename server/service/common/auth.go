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
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct{}

// maxLoginCandidates 未指定租户登录时，最多逐行比对的同名账号数。
//
// 「同一个人在 N 个租户里就是 N 行 sys_user」的克隆模型下，未指定租户的登录必须把
// 口令拿去逐行 bcrypt 比对——成本随克隆数线性放大（P2-6）。图形验证码、IP 限流与
// 用户名维度的递增惩罚锁定已经把可放大面压得很小，这里再给单次请求的计算量封顶，
// 顺带把真实用户的登录延迟也钉住（bcrypt cost 10 ≈ 每行数十毫秒）。
//
// 同名账号超过该上限时，超出的行不参与比对：这类用户请在登录页填写租户编码，
// 走「指定租户」路径（只做 1 次比对，且不受上限影响）。
const maxLoginCandidates = 20

// Login 登录。返回命中的用户（带角色）及其归属的启用业务租户列表。
// 未指定租户且账号归属多个启用租户时返回 requiresTenantSelection，
// 由前端弹窗选择后经 select_tenant 完成登录。
func (s *AuthService) Login(db *gorm.DB, req common.Login) (*system.SysUser, []commonResponse.TenantOption, error) {
	tenantId := req.TenantId
	if tenantId == "" && req.TenantCode != "" {
		resolved, err := s.resolveTenantByCode(db, req.TenantCode)
		if err != nil {
			return nil, nil, err
		}
		tenantId = resolved
	}
	if tenantId == "" {
		return s.loginAcrossTenants(db, req)
	}

	if err := s.checkTenantUsable(db, tenantId); err != nil {
		return nil, nil, err
	}
	user, err := s.Authenticate(db, tenantId, req.Username, req.Password)
	if err != nil {
		return nil, nil, err
	}
	return user, s.tenantOptions(db, tenantId), nil
}

// resolveTenantByCode 把登录页填写的租户编码解析为租户 ID。
//
// 编码是给人看的登录入口（租户选择弹窗里已经展示它），不是秘密；解析失败按
// 「租户不存在/已停用」直接答复，而不是伪装成口令错误——否则用户填错编码时
// 只会看到「用户名或密码错误」，无从排查。平台域没有 sys_tenant 行，也就没有
// 编码可填：平台管理员留空编码登录。
func (s *AuthService) resolveTenantByCode(db *gorm.DB, code string) (string, error) {
	var tenant system.SysTenant
	if err := db.Where("code = ?", code).First(&tenant).Error; err != nil {
		return "", errors.New("tenantNotFound")
	}
	if !tenant.Status {
		return "", errors.New("tenantDisabled")
	}
	return tenant.TenantId, nil
}

// tenantOptions 显式指定租户登录时的「可切换租户」列表：只有目标租户本身。
//
// 要判断「同一口令在别的租户下也有账号」，只能把口令拿去逐行 bcrypt 比对，
// 那正是 P2-6 的成本来源。既然用户已经指明了要进哪个租户，就不再为了拼一份
// 切换列表付这笔钱——本次会话被钉在该租户上（my_tenants / switch_tenant 都以
// token 里的授权集合为准）。需要跨租户切换时留空租户编码登录。
func (s *AuthService) tenantOptions(db *gorm.DB, tenantId string) []commonResponse.TenantOption {
	if tenantId == "" || tenantId == global.PlatformTenantID {
		return nil
	}
	var tenant system.SysTenant
	if err := db.Where("tenant_id = ? AND status = ?", tenantId, true).First(&tenant).Error; err != nil {
		return nil
	}
	return []commonResponse.TenantOption{{
		TenantId: tenant.TenantId,
		Name:     tenant.Name,
		Code:     tenant.Code,
	}}
}

// loginAcrossTenants 未指定租户时的登录解析：
// 校验用户名密码后按归属的启用租户分发——唯一归属直接进入；
// 多个启用租户时返回 requiresTenantSelection 引导前端弹窗选择；
func (s *AuthService) loginAcrossTenants(db *gorm.DB, req common.Login) (*system.SysUser, []commonResponse.TenantOption, error) {
	business, platform, options, err := s.matchEnabledTenants(db, req)
	if err != nil {
		return nil, nil, err
	}

	switch len(business) {
	case 0:
		// 仅平台归属：平台管理员直接登录平台域，不出现在租户选择器中
		return &platform[0], options, nil
	case 1:
		return &business[0], options, nil
	default:
		return nil, options, errors.New("requiresTenantSelection")
	}
}

// matchEnabledTenants 返回本次口令实际验证通过的启用租户账号及其租户选项。
// 同名但口令不同的账号不会进入集合，也不会出现在任何租户选择/切换界面中。
// 逐行比对的数量以 maxLoginCandidates 封顶，见该常量的说明。
func (s *AuthService) matchEnabledTenants(db *gorm.DB, req common.Login) ([]system.SysUser, []system.SysUser, []commonResponse.TenantOption, error) {
	var rows []system.SysUser
	if err := db.Preload("Roles").
		Where("username = ?", req.Username).
		Order("created_at ASC").
		Limit(maxLoginCandidates).
		Find(&rows).Error; err != nil {
		return nil, nil, nil, errors.New("userNotFound")
	}
	if len(rows) == 0 {
		return nil, nil, nil, errors.New("userNotFound")
	}
	if len(rows) == maxLoginCandidates {
		// 命中上限：可能还有更靠后的同名账号没参与比对。记一条日志，
		// 让运维知道这些用户需要改用「租户编码」登录，而不是把它当成口令错误查。
		global.Log.Warn("同名账号数达到跨租户登录比对上限，超出的账号需指定租户编码登录",
			zap.String("username", req.Username),
			zap.Int("limit", maxLoginCandidates))
	}

	var (
		business []system.SysUser
		platform []system.SysUser
		options  []commonResponse.TenantOption
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
		// 到期视同停用：过期的租户不出现在可切换列表里。
		if util.IsTenantExpired(tenant.ExpireTime) {
			continue
		}
		business = append(business, u)
		options = append(options, commonResponse.TenantOption{
			TenantId: tenant.TenantId,
			Name:     tenant.Name,
			Code:     tenant.Code,
		})
	}

	if len(business) == 0 && len(platform) == 0 {
		if !matched {
			return nil, nil, nil, errors.New("passwordError")
		}
		return nil, nil, nil, errors.New("tenantDisabled")
	}
	return business, platform, options, nil
}

// checkTenantUsable 校验指定租户存在且启用。空租户与平台保留租户始终可用：
// 前者是公共端点/迁移等无租户上下文，后者是跨租户管理域，都不在 sys_tenant 里。
func (s *AuthService) checkTenantUsable(db *gorm.DB, tenantId string) error {
	if tenantId == "" || tenantId == global.PlatformTenantID {
		return nil
	}
	var tenant system.SysTenant
	if err := db.Where("tenant_id = ?", tenantId).First(&tenant).Error; err != nil {
		return errors.New("tenantNotFound")
	}
	// 到期视同停用：过期租户与停用租户同样拒绝登录/切换。
	if !tenant.Status || util.IsTenantExpired(tenant.ExpireTime) {
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

// EnterTenantByUserIds 校验租户可用并从候选账号集合中取出该租户下的账号（含角色）。
//
// 邮箱登录不能复用按用户名的 EnterTenant：同一邮箱在不同租户下可以是不同的用户名，
// 用 username+tenant 反查会指向另一个人的账号。因此这一路只按候选账号 ID 收敛，
// 候选集来自本次验证码绑定的集合，客户端无法扩大。
func (s *AuthService) EnterTenantByUserIds(db *gorm.DB, tenantId string, userIds []string) (*system.SysUser, error) {
	if err := s.checkTenantUsable(db, tenantId); err != nil {
		return nil, err
	}
	if len(userIds) == 0 {
		return nil, errors.New("userNotFound")
	}
	user := &system.SysUser{}
	err := db.
		Preload("Roles").
		Where("tenant_id = ?", tenantId).
		Where("id IN ?", userIds).
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
