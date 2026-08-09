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
	ExpireTime        time.Duration // 有效时间
	RefreshExpireTime time.Duration
	Disabled          bool
	CreatedTime       time.Time
	// MustChangePassword 标记该账号仍需强制修改密码
	MustChangePassword bool
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
		RefreshToken:       refreshToken,
		Disabled:           t.Disabled,
		CreatedTime:        t.CreatedTime,
		ExpiredTime:        t.CreatedTime.Add(t.ExpireTime),
		MustChangePassword: t.MustChangePassword,
	})
	if err != nil {
		return nil, err
	}
	refreshTokenInfoJson, err := json.Marshal(model.RefreshTokenInfo{
		UserId:             t.UserId,
		Username:           t.Username,
		RoleIds:            t.RoleIds,
		Disabled:           t.Disabled,
		CreatedTime:        t.CreatedTime,
		ExpiredTime:        t.CreatedTime.Add(t.RefreshExpireTime),
		MustChangePassword: t.MustChangePassword,
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
	}, nil
}

func (t *Token) RefreshToken(refreshToken string) (*common.Token, error) {
	ctx := context.Background()

	// 获取refreshToken对应的userId
	refreshTokenString, err := global.Redis.Get(ctx, fmt.Sprintf("%s%s", PrefixRefreshToken, refreshToken)).Result()
	if err != nil {
		return nil, errors.New("refreshTokenNotExist")
	}
	var refreshTokenInfo model.RefreshTokenInfo
	err = json.Unmarshal([]byte(refreshTokenString), &refreshTokenInfo)
	if err != nil {
		return nil, errors.New("refreshTokenNotExist")
	}

	// 清理该用户的无效 access/refresh token
	_ = CleanUserTokenSet(refreshTokenInfo.UserId, PrefixUserAcessToken, PrefixAccessToken)
	_ = CleanUserTokenSet(refreshTokenInfo.UserId, PrefixUserRefreshToken, PrefixRefreshToken)

	// 生成新的accessToken
	newAccessToken, err := Uuid()
	if err != nil {
		return nil, errors.New("refreshFailed")
	}

	// 获取refreshToken的剩余过期时间
	refreshExpire, err := global.Redis.TTL(ctx, fmt.Sprintf("%s%s", PrefixRefreshToken, refreshToken)).Result()
	if err != nil || refreshExpire <= 0 {
		return nil, errors.New("refreshTokenExpired")
	}

	tokenInfoJson, err := json.Marshal(model.TokenInfo{
		UserId:             refreshTokenInfo.UserId,
		Username:           refreshTokenInfo.Username,
		RoleIds:            refreshTokenInfo.RoleIds,
		RefreshToken:       refreshToken,
		Disabled:           refreshTokenInfo.Disabled,
		CreatedTime:        refreshTokenInfo.CreatedTime,
		ExpiredTime:        refreshTokenInfo.ExpiredTime.Add(t.ExpireTime),
		MustChangePassword: refreshTokenInfo.MustChangePassword,
	})
	if err != nil {
		return nil, errors.New("refreshFailed")
	}
	// 保存新的token信息
	pipe := global.Redis.Pipeline()
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixAccessToken, newAccessToken), tokenInfoJson, t.ExpireTime)
	pipe.SAdd(ctx, fmt.Sprintf("%s%s", PrefixUserAcessToken, refreshTokenInfo.UserId), newAccessToken) // 添加到用户token集合

	// 执行管道操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, errors.New("refreshFailed")
	}

	return &common.Token{
		AccessToken:        newAccessToken,
		ExpireTime:         int64(t.ExpireTime.Seconds()), // 将毫秒转换为秒
		RefreshToken:       refreshToken,
		RefreshExpireTime:  int64(refreshExpire.Seconds()), // 将毫秒转换为秒
		MustChangePassword: refreshTokenInfo.MustChangePassword,
	}, nil
}

// updateTokenRoles 更新指定类型token的角色信息
func updateTokenRoles(ctx context.Context, userId string, roleIds []uint, tokenType string, prefix string) error {
	tokens, err := global.Redis.SMembers(ctx, fmt.Sprintf("%s%s", tokenType, userId)).Result()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		tokenKey := fmt.Sprintf("%s%s", prefix, token)
		tokenInfo, err := global.Redis.Get(ctx, tokenKey).Result()
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

		// 获取原有token的过期时间
		ttl, err := global.Redis.TTL(ctx, tokenKey).Result()
		if err != nil {
			continue
		}

		// 使用原有的过期时间更新token
		if err = global.Redis.Set(ctx, tokenKey, tokenInfoJson, ttl).Err(); err != nil {
			continue
		}
	}
	return nil
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

	pipe := global.Redis.Pipeline()
	for _, token := range tokens {
		tokenKey := fmt.Sprintf("%s%s", prefix, token)

		// 检查 token 对应的 Redis 键是否存在
		exists, e := global.Redis.Exists(ctx, tokenKey).Result()
		if e != nil {
			continue
		}

		// 如果 token 已过期（Redis 中不存在），则从用户集合中移除
		if exists == 0 {
			pipe.SRem(ctx, setKey, token)
		}
	}
	_, err = pipe.Exec(ctx)
	return err
}
