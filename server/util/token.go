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
	PrefixAccessToken    = "accessToken:"
	PrefixRefreshToken   = "refreshToken:"
	PrefixUserAcessToken = "userAccessToken_"
	PrefixUserRefresh    = "userRefreshToken_"
)

type Token struct {
	UserId            int64
	RoleIds           []uint
	ExpireTime        time.Duration
	RefreshExpireTime time.Duration
}

type refreshInfo struct {
	UserId  int64  `json:"user_id"`
	RoleIds []uint `json:"role_ids"`
}

func (t *Token) GenerateToken() (common.Token, error) {
	accessToken, err := Uuid()
	if err != nil {
		return common.Token{}, err
	}
	refreshToken, err := Uuid()
	if err != nil {
		return common.Token{}, err
	}

	ctx := context.Background()
	tokenInfoJson, err := json.Marshal(model.Auth{
		UserID:       t.UserId,
		RoleIds:      t.RoleIds,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return common.Token{}, err
	}
	refreshTokenJson, err := json.Marshal(refreshInfo{
		UserId:  t.UserId,
		RoleIds: t.RoleIds,
	})

	pipe := global.Redis.Pipeline()
	// 保存token相关信息
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixAccessToken, accessToken), tokenInfoJson, t.ExpireTime)
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixRefreshToken, refreshToken), refreshTokenJson, t.RefreshExpireTime)

	// 将token和refreshToken添加到用户的token集合中
	pipe.SAdd(ctx, fmt.Sprintf("%s%d", PrefixUserAcessToken, t.UserId), accessToken)
	pipe.SAdd(ctx, fmt.Sprintf("%s%d", PrefixRefreshToken, t.UserId), refreshToken)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return common.Token{}, err
	}

	return common.Token{
		AccessToken:       accessToken,
		ExpireTime:        int64(t.ExpireTime.Seconds()), // 将毫秒转换为秒
		RefreshToken:      refreshToken,
		RefreshExpireTime: int64(t.RefreshExpireTime.Seconds()), // 将毫秒转换为秒
	}, nil
}

func (t *Token) RefreshToken(refreshToken string) (common.Token, error) {
	ctx := context.Background()

	// 获取refreshToken对应的userId
	refreshTokenString, err := global.Redis.Get(ctx, fmt.Sprintf("%s%s", PrefixRefreshToken, refreshToken)).Result()
	if err != nil {
		return common.Token{}, errors.New("refreshTokenNotExist")
	}
	var refreshTokenInfo refreshInfo
	err = json.Unmarshal([]byte(refreshTokenString), &refreshTokenInfo)
	if err != nil {
		return common.Token{}, errors.New("refreshTokenNotExist")
	}

	// 生成新的accessToken
	newAccessToken, err := Uuid()
	if err != nil {
		return common.Token{}, errors.New("refreshFailed")
	}

	// 获取refreshToken的剩余过期时间
	refreshExpire, err := global.Redis.TTL(ctx, fmt.Sprintf("%s%s", PrefixRefreshToken, refreshToken)).Result()
	if err != nil || refreshExpire <= 0 {
		return common.Token{}, errors.New("refreshTokenExpired")
	}

	tokenInfoJson, err := json.Marshal(model.Auth{
		UserID:       refreshTokenInfo.UserId,
		RoleIds:      refreshTokenInfo.RoleIds,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return common.Token{}, errors.New("refreshFailed")
	}
	// 保存新的token信息
	pipe := global.Redis.Pipeline()
	pipe.Set(ctx, fmt.Sprintf("%s%s", PrefixAccessToken, newAccessToken), tokenInfoJson, t.ExpireTime)

	// 执行管道操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return common.Token{}, errors.New("refreshFailed")
	}

	return common.Token{
		AccessToken:       newAccessToken,
		ExpireTime:        int64(t.ExpireTime.Seconds()), // 将毫秒转换为秒
		RefreshToken:      refreshToken,
		RefreshExpireTime: int64(refreshExpire.Seconds()), // 将毫秒转换为秒
	}, nil
}

// updateTokenRoles 更新指定类型token的角色信息
func updateTokenRoles(ctx context.Context, userId int64, roleIds []uint, tokenType string, prefix string) error {
	tokens, err := global.Redis.SMembers(ctx, fmt.Sprintf("%s%d", tokenType, userId)).Result()
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
			tokenData = &model.Auth{}
		} else {
			tokenData = &refreshInfo{}
		}

		if err = json.Unmarshal([]byte(tokenInfo), tokenData); err != nil {
			continue
		}

		// 更新角色ID
		switch v := tokenData.(type) {
		case *model.Auth:
			v.RoleIds = roleIds
		case *refreshInfo:
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
func UpdateUserRoleInCache(userId int64, roleIds []uint) error {
	ctx := context.Background()

	// 更新访问令牌的角色信息
	if err := updateTokenRoles(ctx, userId, roleIds, PrefixUserAcessToken, PrefixAccessToken); err != nil {
		return err
	}

	// 更新刷新令牌的角色信息
	if err := updateTokenRoles(ctx, userId, roleIds, PrefixUserRefresh, PrefixRefreshToken); err != nil {
		return err
	}

	return nil
}
