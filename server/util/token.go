package util

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
)

type Token struct {
	UserId            uint
	ExpireTime        time.Duration
	RefreshExpireTime time.Duration
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
	tokenInfoJson, err := json.Marshal(map[string]string{
		"userId":       fmt.Sprintf("%d", t.UserId),
		"refreshToken": refreshToken,
	})
	if err != nil {
		return common.Token{}, err
	}
	pipe := global.Redis.Pipeline()
	// 保存token相关信息
	pipe.Set(ctx, fmt.Sprintf("accessToken:%s", accessToken), tokenInfoJson, t.ExpireTime)
	pipe.Set(ctx, fmt.Sprintf("refreshToken:%s", refreshToken), t.UserId, t.RefreshExpireTime)

	// 将token和refreshToken添加到用户的token集合中
	pipe.SAdd(ctx, fmt.Sprintf("userAccessToken_%d", t.UserId), accessToken)
	pipe.SAdd(ctx, fmt.Sprintf("userRefreshToken_%d", t.UserId), refreshToken)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return common.Token{}, err
	}

	return common.Token{
		AccessToken:       accessToken,
		ExpireTime:        int64(t.ExpireTime.Milliseconds()),
		RefreshToken:      refreshToken,
		RefreshExpireTime: int64(t.RefreshExpireTime.Milliseconds()),
	}, nil
}

func (t *Token) RefreshToken(refreshToken string) (common.Token, error) {
	ctx := context.Background()

	// 获取refreshToken对应的userId
	userId, err := global.Redis.Get(ctx, fmt.Sprintf("refreshToken:%s", refreshToken)).Result()
	if err != nil {
		return common.Token{}, err
	}

	// 生成新的accessToken
	newAccessToken, err := Uuid()
	if err != nil {
		return common.Token{}, err
	}

	// 获取refreshToken的剩余过期时间
	refreshExpire, err := global.Redis.TTL(ctx, fmt.Sprintf("refreshToken:%s", refreshToken)).Result()
	if err != nil {
		return common.Token{}, err
	}

	tokenInfoJson, err := json.Marshal(map[string]string{
		"userId":       userId,
		"refreshToken": refreshToken,
	})
	if err != nil {
		return common.Token{}, err
	}
	// 保存新的token信息
	pipe := global.Redis.Pipeline()
	pipe.Set(ctx, fmt.Sprintf("accessToken:%s", newAccessToken), tokenInfoJson, t.ExpireTime)

	// 执行管道操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return common.Token{}, err
	}

	return common.Token{
		AccessToken:       newAccessToken,
		ExpireTime:        int64(t.ExpireTime.Milliseconds()),
		RefreshToken:      refreshToken,
		RefreshExpireTime: int64(refreshExpire.Milliseconds()),
	}, nil
}
