package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type SysTokenService struct{}

// TokenScope 描述令牌管理请求的可见范围。
// 会话保存在 Redis 里，不经过 GORM，租户插件管不到，必须在服务层按
// TokenInfo.TenantId 自行裁剪。
type TokenScope struct {
	// TenantId 调用者所属租户
	TenantId string
	// IsSuperAdmin 平台超级管理员看全部租户的会话
	IsSuperAdmin bool
}

// visible 判断某个会话是否在当前调用者的可见范围内。
// 超管不受限；其余人只能看到与自己同租户的会话。
// 单租户模式下所有会话的 TenantId 一致，判定自然恒真。
func (s TokenScope) visible(tenantId string) bool {
	if s.IsSuperAdmin {
		return true
	}
	return tenantId == s.TenantId
}

// FindAll 获取令牌列表（按调用者的租户可见范围裁剪）
func (s *SysTokenService) FindAll(scope TokenScope) ([]response.SysTokenLogRsp, error) {
	ctx := context.Background()
	log := global.Log
	rdb := global.Redis

	accessTokenKeys := []string{}
	refreshTokenKeys := []string{}
	var cursor uint64
	pattern := fmt.Sprintf("%s*", util.PrefixUserAcessToken)
	// SCAN 所有 userAccessToken_:* 的 key
	for {
		keys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			log.Error("Failed to scan user token sets", zap.Error(err))
			return nil, errors.New("getFailed")
		}
		cursor = nextCursor

		for _, key := range keys {
			tokens, err := rdb.SMembers(ctx, key).Result()
			if err != nil {
				log.Warn("Failed to read token set", zap.String("key", key), zap.Error(err))
				continue
			}
			for _, token := range tokens {
				accessTokenKey := fmt.Sprintf("%s%s", util.PrefixAccessToken, token)
				accessTokenKeys = append(accessTokenKeys, accessTokenKey)
			}
		}

		if cursor == 0 {
			break
		}
	}

	if len(accessTokenKeys) == 0 {
		log.Info("No access tokens found")
		return []response.SysTokenLogRsp{}, nil
	}

	// Step 2: Pipeline 批量 GET 所有 token 值
	pipe := rdb.Pipeline()
	cmds := make([]*redis.StringCmd, len(accessTokenKeys))
	for i, key := range accessTokenKeys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Error("Failed to exec pipeline for access tokens", zap.Error(err))
		return nil, errors.New("getFailed")
	}

	// Step 3: 解析 JSON 为结构体并收集refreshToken keys
	result := make(map[string]model.TokenInfo)
	for i, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil && err != redis.Nil {
			log.Warn("Failed to get access token value", zap.String("key", accessTokenKeys[i]), zap.Error(err))
			continue
		}

		var info model.TokenInfo
		if err := json.Unmarshal([]byte(val), &info); err != nil {
			log.Warn("Failed to parse token JSON", zap.String("key", accessTokenKeys[i]), zap.Error(err))
			continue
		}

		// 跨租户会话直接丢弃，不进入后续 refresh token 查询与响应
		if !scope.visible(info.TenantId) {
			continue
		}

		result[accessTokenKeys[i]] = info
		// 收集refreshToken keys
		refreshTokenKey := fmt.Sprintf("%s%s", util.PrefixRefreshToken, info.RefreshToken)
		refreshTokenKeys = append(refreshTokenKeys, refreshTokenKey)
	}

	// Step 4: Pipeline 批量 GET 所有 refreshToken 值
	refreshPipe := rdb.Pipeline()
	refreshCmds := make([]*redis.StringCmd, len(refreshTokenKeys))
	for i, key := range refreshTokenKeys {
		refreshCmds[i] = refreshPipe.Get(ctx, key)
	}

	_, err = refreshPipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Error("Failed to exec pipeline for refresh tokens", zap.Error(err))
		return nil, errors.New("getFailed")
	}

	// Step 5: 解析 refreshToken JSON 为结构体
	refreshResult := make(map[string]model.RefreshTokenInfo)
	for i, cmd := range refreshCmds {
		val, err := cmd.Result()
		if err != nil && err != redis.Nil {
			log.Warn("Failed to get refresh token value", zap.String("key", refreshTokenKeys[i]), zap.Error(err))
			continue
		}

		var info model.RefreshTokenInfo
		if err := json.Unmarshal([]byte(val), &info); err != nil {
			log.Warn("Failed to parse refresh token JSON", zap.String("key", refreshTokenKeys[i]), zap.Error(err))
			continue
		}

		refreshResult[refreshTokenKeys[i]] = info
	}

	tokens := lo.Map(lo.Entries(result), func(entry lo.Entry[string, model.TokenInfo], _ int) response.SysTokenLogRsp {
		// 获取对应的refreshToken信息
		refreshTokenKey := fmt.Sprintf("%s%s", util.PrefixRefreshToken, entry.Value.RefreshToken)
		refreshInfo, hasRefresh := refreshResult[refreshTokenKey]

		// 如果refreshToken信息不存在，使用默认值
		if !hasRefresh {
			refreshInfo = model.RefreshTokenInfo{
				CreatedTime: entry.Value.CreatedTime, // 使用accessToken的创建时间作为默认值
				ExpiredTime: entry.Value.ExpiredTime, // 使用accessToken的过期时间作为默认值
			}
		}

		return response.SysTokenLogRsp{
			UserId:              entry.Value.UserId,
			TenantId:            entry.Value.TenantId,
			Username:            entry.Value.Username,
			AccessToken:         strings.TrimPrefix(entry.Key, util.PrefixAccessToken),
			RefreshToken:        entry.Value.RefreshToken,
			Disabled:            entry.Value.Disabled,
			AccessTokenCreated:  entry.Value.CreatedTime,
			AccessTokenExpired:  entry.Value.ExpiredTime,
			RefreshTokenCreated: refreshInfo.CreatedTime,
			RefreshTokenExpired: refreshInfo.ExpiredTime,
		}
	})

	slices.SortFunc(tokens, func(a, b response.SysTokenLogRsp) int {
		if a.Username < b.Username {
			return -1
		}
		if a.Username > b.Username {
			return 1
		}
		if a.AccessTokenCreated.Before(b.AccessTokenCreated) {
			return -1
		}
		if a.AccessTokenCreated.After(b.AccessTokenCreated) {
			return 1
		}
		return 0
	})

	return tokens, nil
}

// Delete 删除token（仅限调用者可见范围内的会话）
func (s *SysTokenService) Delete(req request.SysTokenDeleteReq, scope TokenScope) error {
	ctx := context.Background()
	rdb := global.Redis
	log := global.Log

	// 获取令牌详情
	tokenData, err := getJsonAndParse[model.TokenInfo](parseParam{
		key:     fmt.Sprintf("%s%s", util.PrefixAccessToken, req.AccessToken),
		log:     log,
		rdb:     rdb,
		context: ctx,
	})
	if err != nil {
		if err == redis.Nil {
			return errors.New("accessTokenNotExist")
		}
		return errors.New("getFailed")
	}

	// 越权防护：令牌 ID 是随机 UUID，但列表接口之外仍可能被猜测/横传，
	// 跨租户一律按「不存在」处理，既拒绝操作也不泄露该令牌是否存在。
	if !scope.visible(tokenData.TenantId) {
		return errors.New("accessTokenNotExist")
	}

	err = deleteWithTransaction(
		ctx, rdb,
		deleteParam{ // 删除令牌
			key:   fmt.Sprintf("%s%s", util.PrefixUserAcessToken, tokenData.UserId),
			token: req.AccessToken,
			isSet: true,
		},
		deleteParam{ // 删除刷新令牌
			key:   fmt.Sprintf("%s%s", util.PrefixUserRefreshToken, tokenData.UserId),
			token: tokenData.RefreshToken,
			isSet: true,
		},
		deleteParam{ // 删除令牌详情
			key:   fmt.Sprintf("%s%s", util.PrefixAccessToken, req.AccessToken),
			token: tokenData.RefreshToken,
			isSet: false,
		},
		deleteParam{ // 删除刷新令牌详情
			key:   fmt.Sprintf("%s%s", util.PrefixRefreshToken, tokenData.RefreshToken),
			token: tokenData.RefreshToken,
			isSet: false,
		},
	)
	if err != nil {
		log.Error("Failed to delete token", zap.Error(err))
		return errors.New("deleteFailed")
	}

	return nil
}

// ModityStatus 修改token状态（仅限调用者可见范围内的会话）
func (s *SysTokenService) ModityStatus(req request.SysTokenModityStatusReq, scope TokenScope) error {
	ctx := context.Background()
	rdb := global.Redis
	log := global.Log

	aKey := fmt.Sprintf("%s%s", util.PrefixAccessToken, req.AccessToken)
	tokenData, err := getJsonAndParse[model.TokenInfo](parseParam{
		key:     aKey,
		log:     log,
		rdb:     rdb,
		context: ctx,
	})
	if err != nil {
		if err == redis.Nil {
			return errors.New("accessTokenNotExist")
		}
		return errors.New("getFailed")
	}

	// 越权防护：跨租户会话按「不存在」处理，理由同 Delete
	if !scope.visible(tokenData.TenantId) {
		return errors.New("accessTokenNotExist")
	}

	rKey := fmt.Sprintf("%s%s", util.PrefixRefreshToken, tokenData.RefreshToken)
	refreshTokenData, err := getJsonAndParse[model.RefreshTokenInfo](parseParam{
		key:     rKey,
		log:     log,
		rdb:     rdb,
		context: ctx,
	})
	if err != nil {
		if err == redis.Nil {
			return errors.New("refreshTokenNotExist")
		}
		return errors.New("getFailed")
	}

	tokenData.Disabled = *req.Disabled
	refreshTokenData.Disabled = *req.Disabled

	pipe := rdb.TxPipeline()
	err = updateWithTransaction(
		updateParam[model.TokenInfo]{
			key:     aKey,
			data:    *tokenData,
			context: ctx,
			pipe:    pipe,
			log:     log,
		},
	)
	if err != nil {
		log.Error("Failed to update token", zap.Error(err))
		return errors.New("updateFailed")
	}

	err = updateWithTransaction(
		updateParam[model.RefreshTokenInfo]{
			key:     rKey,
			data:    *refreshTokenData,
			context: ctx,
			pipe:    pipe,
			log:     log,
		},
	)
	if err != nil {
		log.Error("Failed to update token", zap.Error(err))
		return errors.New("updateFailed")
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("Failed to execute pipeline", zap.Error(err))
		return errors.New("updateFailed")
	}

	// 停用 token 后其在线连接已失去凭证，推送强制下线让前端立即退出重登
	if *req.Disabled {
		notifySessionRevoked(tokenData.UserId, "您的会话已被管理员终止")
	}

	return nil
}

type deleteParam struct {
	key   string
	token string
	isSet bool
}

// 删除redis中的数据
func deleteWithTransaction(ctx context.Context, rdb *redis.Client, params ...deleteParam) error {
	pipe := rdb.TxPipeline()

	for _, param := range params {
		if param.isSet {
			// 检查是否存在再删除（可以在 pipeline 外做检查）
			pipe.SRem(ctx, param.key, param.token)
		} else {
			pipe.Del(ctx, param.key)
		}
	}

	_, err := pipe.Exec(ctx)
	return err
}

type updateParam[T any] struct {
	key     string
	data    T
	context context.Context
	pipe    redis.Pipeliner
	log     *zap.Logger
}

// 更新结构体并保存到redis
func updateWithTransaction[T any](param updateParam[T]) error {
	tokenInfoJson, err := json.Marshal(param.data)
	if err != nil {
		param.log.Error("Failed to marshal token JSON", zap.Error(err))
		return err
	}

	// 获取原有token的过期时间
	ttl, err := param.pipe.TTL(param.context, param.key).Result()
	if err != nil {
		param.log.Error("Failed to get token TTL", zap.Error(err))
		return err
	}

	// 使用原有的过期时间更新token
	err = param.pipe.Set(param.context, param.key, tokenInfoJson, ttl).Err()
	return err
}

type parseParam struct {
	key     string
	context context.Context
	rdb     *redis.Client
	log     *zap.Logger
}

// 从redis获取json并解析为结构体
func getJsonAndParse[T any](param parseParam) (*T, error) {
	tokenInfo, err := param.rdb.Get(param.context, param.key).Result()
	if err != nil {
		param.log.Error("Failed to get token info", zap.Error(err))
		return nil, err
	}
	var tokenData T
	if err = json.Unmarshal([]byte(tokenInfo), &tokenData); err != nil {
		param.log.Error("Failed to parse token JSON", zap.Error(err))
		return nil, err
	}

	return &tokenData, nil
}
