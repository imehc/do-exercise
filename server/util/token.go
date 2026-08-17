package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/redis/go-redis/v9"
)

const (
	PrefixAccessToken      = "accessToken:"
	PrefixRefreshToken     = "refreshToken:"
	PrefixUserAcessToken   = "userAccessToken_"
	PrefixUserRefreshToken = "userRefreshToken_"
)

type Token struct {
	UserId            string
	Username          string
	RoleIds           []uint
	TenantId          string
	ExpireTime        time.Duration // 有效时间
	RefreshExpireTime time.Duration
	Disabled          bool
	CreatedTime       time.Time
	// MustChangePassword 标记该账号仍需强制修改密码
	MustChangePassword bool
	// IsSuperAdmin 平台超级管理员标识（仅平台域账号有效）
	IsSuperAdmin bool
}

func (t *Token) GenerateToken() (*common.Token, error) {
	accessToken, err := Uuid()
	if err != nil {
		return nil, err
	}
	refreshToken, err := Uuid()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tokenInfoJson, err := json.Marshal(model.TokenInfo{
		UserId:             t.UserId,
		Username:           t.Username,
		RoleIds:            t.RoleIds,
		TenantId:           t.TenantId,
		RefreshToken:       refreshToken,
		Disabled:           t.Disabled,
		CreatedTime:        t.CreatedTime,
		ExpiredTime:        t.CreatedTime.Add(t.ExpireTime),
		MustChangePassword: t.MustChangePassword,
		IsSuperAdmin:       t.IsSuperAdmin,
	})
	if err != nil {
		return nil, err
	}
	refreshTokenInfoJson, err := json.Marshal(model.RefreshTokenInfo{
		UserId:             t.UserId,
		Username:           t.Username,
		RoleIds:            t.RoleIds,
		TenantId:           t.TenantId,
		Disabled:           t.Disabled,
		CreatedTime:        t.CreatedTime,
		ExpiredTime:        t.CreatedTime.Add(t.RefreshExpireTime),
		MustChangePassword: t.MustChangePassword,
		IsSuperAdmin:       t.IsSuperAdmin,
	})
	if err != nil {
		return nil, err
	}

	pipe := global.Redis.Pipeline()
	// 保存token相关信息
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixAccessToken, accessToken), tokenInfoJson, t.ExpireTime)
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixRefreshToken, refreshToken), refreshTokenInfoJson, t.RefreshExpireTime)

	// 将token和refreshToken添加到用户的token集合中
	pipe.SAdd(ctx, fmt.Sprintf("%s%s", PrefixUserAcessToken, t.UserId), accessToken)
	pipe.SAdd(ctx, fmt.Sprintf("%s%s", PrefixUserRefreshToken, t.UserId), refreshToken)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &common.Token{
		AccessToken:        accessToken,
		ExpireTime:         int64(t.ExpireTime.Seconds()), // 将毫秒转换为秒
		RefreshToken:       refreshToken,
		RefreshExpireTime:  int64(t.RefreshExpireTime.Seconds()), // 将毫秒转换为秒
		MustChangePassword: t.MustChangePassword,
		TenantId:           t.TenantId,
		IsSuperAdmin:       t.IsSuperAdmin,
	}, nil
}

func (t *Token) RefreshToken(refreshToken string) (*common.Token, error) {
	ctx := context.Background()

	// 获取refreshToken对应的userId
	refreshKey := fmt.Sprintf("%s%s", PrefixRefreshToken, refreshToken)
	refreshTokenString, err := global.Redis.Get(ctx, refreshKey).Result()
	if err != nil {
		return nil, errors.New("refreshTokenNotExist")
	}
	var refreshTokenInfo model.RefreshTokenInfo
	err = json.Unmarshal([]byte(refreshTokenString), &refreshTokenInfo)
	if err != nil {
		return nil, errors.New("refreshTokenNotExist")
	}

	// 重放检测：已被轮转消费过的 refresh token 再次出现，意味着旧凭据可能已泄露，
	// 判定整个 token 家族失陷，吊销该用户全部会话（此时即便攻击者持有新 token 也一起作废）。
	if refreshTokenInfo.Rotated {
		_ = RevokeAllUserTokens(refreshTokenInfo.UserId)
		return nil, errors.New("refreshTokenNotExist")
	}

	// 已禁用或已过期的 refresh token 直接拒绝
	if refreshTokenInfo.Disabled {
		return nil, errors.New("refreshTokenNotExist")
	}
	if !refreshTokenInfo.ExpiredTime.IsZero() && time.Now().After(refreshTokenInfo.ExpiredTime) {
		return nil, errors.New("refreshTokenExpired")
	}
	// 以 Redis TTL 兜底（过期即不可用）
	remainingTTL, err := global.Redis.TTL(ctx, refreshKey).Result()
	if err != nil || remainingTTL <= 0 {
		return nil, errors.New("refreshTokenExpired")
	}

	// 清理该用户的无效 access/refresh token
	_ = CleanUserTokenSet(refreshTokenInfo.UserId, PrefixUserAcessToken, PrefixAccessToken)
	_ = CleanUserTokenSet(refreshTokenInfo.UserId, PrefixUserRefreshToken, PrefixRefreshToken)

	// 轮转：生成全新的 access 与 refresh token，旧 refresh 转为失陷哨兵
	newAccessToken, err := Uuid()
	if err != nil {
		return nil, errors.New("refreshFailed")
	}
	newRefreshToken, err := Uuid()
	if err != nil {
		return nil, errors.New("refreshFailed")
	}

	now := time.Now()
	accessExpire := t.ExpireTime
	if accessExpire <= 0 {
		accessExpire = 2 * time.Hour
	}
	refreshExpire := t.RefreshExpireTime
	if refreshExpire <= 0 {
		refreshExpire = 7 * 24 * time.Hour
	}

	accessInfoJson, err := json.Marshal(model.TokenInfo{
		UserId:             refreshTokenInfo.UserId,
		Username:           refreshTokenInfo.Username,
		RoleIds:            refreshTokenInfo.RoleIds,
		TenantId:           refreshTokenInfo.TenantId,
		RefreshToken:       newRefreshToken, // 新 access 指向新 refresh
		Disabled:           refreshTokenInfo.Disabled,
		CreatedTime:        now,
		ExpiredTime:        now.Add(accessExpire), // 从当前时间重新起算，不再累加原过期时间
		MustChangePassword: refreshTokenInfo.MustChangePassword,
		IsSuperAdmin:       refreshTokenInfo.IsSuperAdmin,
	})
	if err != nil {
		return nil, errors.New("refreshFailed")
	}
	refreshInfoJson, err := json.Marshal(model.RefreshTokenInfo{
		UserId:             refreshTokenInfo.UserId,
		Username:           refreshTokenInfo.Username,
		RoleIds:            refreshTokenInfo.RoleIds,
		TenantId:           refreshTokenInfo.TenantId,
		Disabled:           refreshTokenInfo.Disabled,
		CreatedTime:        now,
		ExpiredTime:        now.Add(refreshExpire),
		MustChangePassword: refreshTokenInfo.MustChangePassword,
		IsSuperAdmin:       refreshTokenInfo.IsSuperAdmin,
	})
	if err != nil {
		return nil, errors.New("refreshFailed")
	}
	// 旧 refresh token 记录改写为 rotated 哨兵，保留剩余 TTL 以便捕获重放；
	// 下一次携带同一 token 到来时走上面的 family-revocation 分支。
	rotatedJson, err := json.Marshal(model.RefreshTokenInfo{
		UserId:             refreshTokenInfo.UserId,
		Username:           refreshTokenInfo.Username,
		RoleIds:            refreshTokenInfo.RoleIds,
		TenantId:           refreshTokenInfo.TenantId,
		Disabled:           true,
		CreatedTime:        refreshTokenInfo.CreatedTime,
		ExpiredTime:        refreshTokenInfo.ExpiredTime,
		MustChangePassword: refreshTokenInfo.MustChangePassword,
		IsSuperAdmin:       refreshTokenInfo.IsSuperAdmin,
		Rotated:            true,
	})
	if err != nil {
		return nil, errors.New("refreshFailed")
	}

	pipe := global.Redis.Pipeline()
	// 旧 refresh token 转失陷哨兵（保留原 TTL）
	pipe.Set(ctx, refreshKey, rotatedJson, remainingTTL)
	// 保存新 access / refresh
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixAccessToken, newAccessToken), accessInfoJson, accessExpire)
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixRefreshToken, newRefreshToken), refreshInfoJson, refreshExpire)
	// 更新用户 token 集合：新 access 入集，旧 refresh 出集、新 refresh 入集
	userAccessSetKey := fmt.Sprintf("%s%s", PrefixUserAcessToken, refreshTokenInfo.UserId)
	userRefreshSetKey := fmt.Sprintf("%s%s", PrefixUserRefreshToken, refreshTokenInfo.UserId)
	pipe.SAdd(ctx, userAccessSetKey, newAccessToken)
	pipe.SRem(ctx, userRefreshSetKey, refreshToken)
	pipe.SAdd(ctx, userRefreshSetKey, newRefreshToken)

	// 执行管道操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, errors.New("refreshFailed")
	}

	return &common.Token{
		AccessToken:        newAccessToken,
		ExpireTime:         int64(accessExpire.Seconds()),
		RefreshToken:       newRefreshToken,
		RefreshExpireTime:  int64(refreshExpire.Seconds()),
		MustChangePassword: refreshTokenInfo.MustChangePassword,
		TenantId:           refreshTokenInfo.TenantId,
		IsSuperAdmin:       refreshTokenInfo.IsSuperAdmin,
	}, nil
}

// updateTokenRoles 更新指定类型token的角色信息
func updateTokenRoles(ctx context.Context, userId string, roleIds []uint, tokenType string, prefix string) error {
	setKey := fmt.Sprintf("%s%s", tokenType, userId)
	tokens, err := global.Redis.SMembers(ctx, setKey).Result()
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}

	// 批量 GET（1 次往返），替代逐 token 串行 Get
	tokenKeys := make([]string, len(tokens))
	getPipe := global.Redis.Pipeline()
	cmds := make([]*redis.StringCmd, len(tokens))
	for i, token := range tokens {
		tokenKeys[i] = fmt.Sprintf("%s%s", prefix, token)
		cmds[i] = getPipe.Get(ctx, tokenKeys[i])
	}
	if _, err := getPipe.Exec(ctx); err != nil && err != redis.Nil {
		return err
	}

	// 更新角色 ID 后批量回写，用 KeepTTL 保留原过期时间（省掉逐 token 的 TTL 往返）
	setPipe := global.Redis.Pipeline()
	updated := false
	for i, cmd := range cmds {
		tokenInfo, err := cmd.Result()
		if err != nil {
			continue
		}

		var tokenData interface{}
		if prefix == PrefixAccessToken {
			tokenData = &model.TokenInfo{}
		} else {
			tokenData = &model.RefreshTokenInfo{}
		}

		if err = json.Unmarshal([]byte(tokenInfo), tokenData); err != nil {
			continue
		}

		// 更新角色ID
		switch v := tokenData.(type) {
		case *model.TokenInfo:
			v.RoleIds = roleIds
		case *model.RefreshTokenInfo:
			v.RoleIds = roleIds
		}

		tokenInfoJson, err := json.Marshal(tokenData)
		if err != nil {
			continue
		}

		setPipe.SetArgs(ctx, tokenKeys[i], tokenInfoJson, redis.SetArgs{
			TTL:     0,
			KeepTTL: true, // 保留原过期时间，省掉逐 token 的 TTL 往返
		})
		updated = true
	}
	if !updated {
		return nil
	}
	_, err = setPipe.Exec(ctx)
	return err
}

// UpdateUserRoleInCache 更新用户角色缓存信息
func UpdateUserRoleInCache(userId string, roleIds []uint) error {
	ctx := context.Background()

	// 更新访问令牌的角色信息
	if err := updateTokenRoles(ctx, userId, roleIds, PrefixUserAcessToken, PrefixAccessToken); err != nil {
		return err
	}

	// 更新刷新令牌的角色信息
	if err := updateTokenRoles(ctx, userId, roleIds, PrefixUserRefreshToken, PrefixRefreshToken); err != nil {
		return err
	}

	return nil
}

// RevokeAllUserTokens 吊销某用户的全部会话。
// 删除用户、禁用账号、修改密码后必须调用：AuthMiddleware 只校验 Redis、从不回查数据库，
// 不清理的话已签发的 access token 会在 TTL 内继续有效，
// refresh token 更可在其有效期内持续换发新的 access token。
func RevokeAllUserTokens(userId string) error {
	if userId == "" {
		return nil
	}
	ctx := context.Background()

	var firstErr error
	for _, pair := range []struct {
		setPrefix   string
		tokenPrefix string
	}{
		{PrefixUserAcessToken, PrefixAccessToken},
		{PrefixUserRefreshToken, PrefixRefreshToken},
	} {
		setKey := fmt.Sprintf("%s%s", pair.setPrefix, userId)
		tokens, err := global.Redis.SMembers(ctx, setKey).Result()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		pipe := global.Redis.Pipeline()
		for _, token := range tokens {
			pipe.Del(ctx, fmt.Sprintf("%s%s", pair.tokenPrefix, token))
		}
		// 连同索引集合一起删除，避免留下悬空成员
		pipe.Del(ctx, setKey)
		if _, err := pipe.Exec(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// RevokeAllUserTokensExcept 吊销某用户除当前会话外的全部会话，并清除保留会话的强制改密标记。
// 用户自己验证原密码改密后调用：当前会话保留（免重新登录），其余设备全部下线；
// 同时清除保留会话上的 must_change_password，否则 AuthMiddleware 仍会拦截该会话。
func RevokeAllUserTokensExcept(userId string, keepAccessToken string) error {
	if userId == "" || keepAccessToken == "" {
		return RevokeAllUserTokens(userId)
	}
	ctx := context.Background()

	// 找到与保留的 access token 配对的 refresh token（access token 记录里存了它）
	var keepRefreshToken string
	tokenInfoString, err := global.Redis.Get(ctx, fmt.Sprintf("%s%s", PrefixAccessToken, keepAccessToken)).Result()
	if err == nil {
		var tokenInfo model.TokenInfo
		if json.Unmarshal([]byte(tokenInfoString), &tokenInfo) == nil {
			keepRefreshToken = tokenInfo.RefreshToken
		}
	}

	var firstErr error
	for _, pair := range []struct {
		setPrefix   string
		tokenPrefix string
		keep        string
	}{
		{PrefixUserAcessToken, PrefixAccessToken, keepAccessToken},
		{PrefixUserRefreshToken, PrefixRefreshToken, keepRefreshToken},
	} {
		setKey := fmt.Sprintf("%s%s", pair.setPrefix, userId)
		tokens, err := global.Redis.SMembers(ctx, setKey).Result()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		pipe := global.Redis.Pipeline()
		for _, token := range tokens {
			if token != pair.keep {
				pipe.Del(ctx, fmt.Sprintf("%s%s", pair.tokenPrefix, token))
			}
		}
		// 重建集合，只保留当前会话的 token
		pipe.Del(ctx, setKey)
		if pair.keep != "" {
			pipe.SAdd(ctx, setKey, pair.keep)
		}
		if _, err := pipe.Exec(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	// 清除保留会话上的强制改密标记
	for _, key := range []string{
		fmt.Sprintf("%s%s", PrefixAccessToken, keepAccessToken),
		fmt.Sprintf("%s%s", PrefixRefreshToken, keepRefreshToken),
	} {
		if err := setTokenMustChangePassword(ctx, key, false); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// setTokenMustChangePassword 更新 Redis 中 token 记录的 must_change_password 字段，保留其余字段与 TTL。
func setTokenMustChangePassword(ctx context.Context, key string, mustChange bool) error {
	tokenString, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(tokenString), &data); err != nil {
		return err
	}
	data["must_change_password"] = mustChange
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ttl, err := global.Redis.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttl <= 0 {
		return nil
	}
	return global.Redis.Set(ctx, key, b, ttl).Err()
}

// CleanUserTokenSet 清理 userAccessToken 或 userRefreshToken Set 中的无效 token
func CleanUserTokenSet(userId string, tokenType string, prefix string) error {
	ctx := context.Background()

	// 拼接用户的 token 集合 key，如 userAccessToken_<userId>
	setKey := fmt.Sprintf("%s%s", tokenType, userId)

	// 获取集合中的所有 token 字符串
	tokens, err := global.Redis.SMembers(ctx, setKey).Result()
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}

	// 批量 EXISTS，一次往返判断哪些 token 已过期
	tokenKeys := make([]string, len(tokens))
	existPipe := global.Redis.Pipeline()
	cmds := make([]*redis.IntCmd, len(tokens))
	for i, token := range tokens {
		tokenKeys[i] = fmt.Sprintf("%s%s", prefix, token)
		cmds[i] = existPipe.Exists(ctx, tokenKeys[i])
	}
	if _, err := existPipe.Exec(ctx); err != nil {
		return err
	}

	// 只对已失效的 token 执行 SRem
	rmPipe := global.Redis.Pipeline()
	removed := false
	for i, token := range tokens {
		exists, e := cmds[i].Result()
		if e != nil || exists != 0 {
			continue
		}
		rmPipe.SRem(ctx, setKey, token)
		removed = true
	}
	if !removed {
		return nil
	}
	_, err = rmPipe.Exec(ctx)
	return err
}
