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
	PrefixAccessToken    = "accessToken:"
	PrefixRefreshToken   = "refreshToken:"
	PrefixSession        = "session:"
	PrefixUserSessions   = "userSessions:"
	PrefixTenantSessions = "tenantSessions:"
	// Legacy indexes are retained only for rolling-upgrade cleanup.
	PrefixUserAcessToken     = "userAccessToken_"
	PrefixUserRefreshToken   = "userRefreshToken_"
	PrefixTenantAccessTokens = "tenantAccessTokens:"
)

func TenantAccessSetKey(tid string) string { return PrefixTenantAccessTokens + tid }

type SessionRecord struct {
	model.TokenInfo
	AccessToken        string    `json:"access_token"`
	RefreshExpiredTime time.Time `json:"refresh_expired_time"`
}

type Token struct {
	UserId              string
	Username            string
	TenantId            string
	AuthVersion         uint64
	AuthorizedTenantIds []string
	ExpireTime          time.Duration
	RefreshExpireTime   time.Duration
	Disabled            bool
	CreatedTime         time.Time
	MustChangePassword  bool
	IsSuperAdmin        bool
}

var saveSession = redis.NewScript(`
if redis.call('EXISTS', KEYS[1]) ~= 0 then return 0 end
redis.call('SET', KEYS[1], ARGV[1], 'PX', ARGV[4])
redis.call('SET', KEYS[2], ARGV[2], 'PX', ARGV[5])
redis.call('SET', KEYS[3], ARGV[3], 'PX', ARGV[4])
for i=4,5 do
 redis.call('SADD',KEYS[i],ARGV[6])
 if redis.call('PTTL',KEYS[i]) < tonumber(ARGV[4]) then redis.call('PEXPIRE',KEYS[i],ARGV[4]) end
end
return 1
`)

func (t *Token) GenerateToken() (*common.Token, error) {
	principal, err := LoadSessionPrincipal(global.DB, t.UserId, t.TenantId)
	if err != nil {
		return nil, err
	}
	if t.AuthVersion != 0 && t.AuthVersion != principal.AuthVersion {
		return nil, errors.New("refreshTokenNotExist")
	}
	if t.ExpireTime <= 0 || t.RefreshExpireTime <= 0 {
		return nil, errors.New("invalid token lifetime")
	}
	if t.ExpireTime > t.RefreshExpireTime {
		t.ExpireTime = t.RefreshExpireTime
	}
	a, err := Uuid()
	if err != nil {
		return nil, err
	}
	r, err := Uuid()
	if err != nil {
		return nil, err
	}
	sid, err := Uuid()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	info := model.TokenInfo{SessionId: sid, AuthVersion: principal.AuthVersion, TenantVersion: principal.TenantVersion, UserId: principal.UserId, Username: principal.Username, TenantId: principal.TenantId, AuthorizedTenantIds: t.AuthorizedTenantIds, RefreshToken: r, Disabled: t.Disabled, CreatedTime: now, ExpiredTime: now.Add(t.ExpireTime), MustChangePassword: principal.MustChangePassword, IsSuperAdmin: principal.IsSuperAdmin}
	refresh := refreshInfo(info, now.Add(t.RefreshExpireTime))
	session := SessionRecord{TokenInfo: info, AccessToken: a, RefreshExpiredTime: refresh.ExpiredTime}
	aj, _ := json.Marshal(info)
	rj, _ := json.Marshal(refresh)
	sj, _ := json.Marshal(session)
	result, err := saveSession.Run(context.Background(), global.Redis, []string{PrefixSession + sid, PrefixAccessToken + a, PrefixRefreshToken + r, PrefixUserSessions + t.UserId, PrefixTenantSessions + t.TenantId}, sj, aj, rj, t.RefreshExpireTime.Milliseconds(), t.ExpireTime.Milliseconds(), sid).Int()

	if err != nil || result != 1 {
		return nil, errors.New("refreshFailed")
	}
	return tokenResponse(a, r, info, t.RefreshExpireTime), nil
}

func refreshInfo(info model.TokenInfo, expiry time.Time) model.RefreshTokenInfo {
	return model.RefreshTokenInfo{SessionId: info.SessionId, AuthVersion: info.AuthVersion, TenantVersion: info.TenantVersion, UserId: info.UserId, Username: info.Username, TenantId: info.TenantId, AuthorizedTenantIds: info.AuthorizedTenantIds, Disabled: info.Disabled, CreatedTime: info.CreatedTime, ExpiredTime: expiry, MustChangePassword: info.MustChangePassword, IsSuperAdmin: info.IsSuperAdmin}
}
func tokenResponse(a, r string, info model.TokenInfo, refreshRemaining time.Duration) *common.Token {
	return &common.Token{AccessToken: a, RefreshToken: r, ExpireTime: int64(time.Until(info.ExpiredTime).Seconds()), RefreshExpireTime: int64(refreshRemaining.Seconds()), TenantId: info.TenantId, IsSuperAdmin: info.IsSuperAdmin, MustChangePassword: info.MustChangePassword}
}

func ReadSession(ctx context.Context, sid string) (*SessionRecord, error) {
	if sid == "" {
		return nil, errors.New("accessTokenNotExist")
	}
	data, err := global.Redis.Get(ctx, PrefixSession+sid).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("%w: %v", ErrAuthDependencyUnavailable, err)
		}
		return nil, errors.New("accessTokenNotExist")
	}
	var session SessionRecord
	if json.Unmarshal(data, &session) != nil || session.SessionId != sid || time.Now().After(session.RefreshExpiredTime) {
		return nil, errors.New("accessTokenNotExist")
	}
	return &session, nil
}
func ValidateAccessToken(ctx context.Context, access string) (*model.TokenInfo, error) {
	if access == "" {
		return nil, errors.New("accessTokenNotExist")
	}
	data, err := global.Redis.Get(ctx, PrefixAccessToken+access).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("%w: %v", ErrAuthDependencyUnavailable, err)
		}
		return nil, err
	}
	var info model.TokenInfo
	if json.Unmarshal(data, &info) != nil || info.SessionId == "" || info.ExpiredTime.IsZero() || time.Now().After(info.ExpiredTime) {
		return nil, errors.New("accessTokenNotExist")
	}
	session, err := ReadSession(ctx, info.SessionId)
	if err != nil {
		return nil, err
	}
	if session.AccessToken != access || session.Disabled {
		return nil, errors.New("accessTokenNotExist")
	}
	user, err := ValidateSessionVersion(global.DB.WithContext(ctx), info.UserId, info.TenantId, info.AuthVersion, info.TenantVersion)
	if err != nil {
		return nil, err
	}
	info.Username = user.Username
	info.MustChangePassword = user.MustChangePassword
	info.IsSuperAdmin = user.IsSuperAdmin
	return &info, nil
}

var rotateSession = redis.NewScript(`
if redis.call('GET',KEYS[1]) ~= ARGV[1] or redis.call('GET',KEYS[2]) ~= ARGV[2] then return 0 end
if redis.call('PTTL',KEYS[1]) <= 0 or redis.call('PTTL',KEYS[2]) <= 0 then return 0 end
redis.call('SET',KEYS[1],ARGV[3],'KEEPTTL')
redis.call('DEL',KEYS[3])
redis.call('SET',KEYS[4],ARGV[4],'PX',ARGV[7])
redis.call('SET',KEYS[5],ARGV[5],'PX',ARGV[8])
redis.call('SET',KEYS[2],ARGV[6],'PX',ARGV[8])
return 1
`)

func (t *Token) RefreshToken(token string) (*common.Token, error) {
	ctx := context.Background()
	raw, err := global.Redis.Get(ctx, PrefixRefreshToken+token).Result()
	if err != nil {
		return nil, errors.New("refreshTokenNotExist")
	}
	var old model.RefreshTokenInfo
	if json.Unmarshal([]byte(raw), &old) != nil || old.SessionId == "" {
		return nil, errors.New("refreshTokenNotExist")
	}
	if old.Rotated {
		if old.RotatedAt.IsZero() || time.Since(old.RotatedAt) > 5*time.Second {
			_ = RevokeSession(ctx, old.SessionId)
		}
		return nil, errors.New("refreshTokenNotExist")
	}
	if old.Disabled || time.Now().After(old.ExpiredTime) {
		return nil, errors.New("refreshTokenExpired")
	}
	user, err := ValidateSessionVersion(global.DB, old.UserId, old.TenantId, old.AuthVersion, old.TenantVersion)
	if err != nil {
		return nil, err
	}
	sraw, err := global.Redis.Get(ctx, PrefixSession+old.SessionId).Result()
	if err != nil {
		return nil, errors.New("refreshTokenNotExist")
	}
	var session SessionRecord
	if json.Unmarshal([]byte(sraw), &session) != nil || session.Disabled || session.RefreshToken != token {
		return nil, errors.New("refreshTokenNotExist")
	}
	a, err := Uuid()
	if err != nil {
		return nil, err
	}
	r, err := Uuid()
	if err != nil {
		return nil, err
	}
	lifetime := t.ExpireTime
	if lifetime <= 0 {
		lifetime = AuthAccessExpire()
	}
	remaining := time.Until(old.ExpiredTime)
	if remaining <= 0 {
		return nil, errors.New("refreshTokenExpired")
	}
	if lifetime > remaining {
		lifetime = remaining
	}
	now := time.Now()
	info := session.TokenInfo
	info.RefreshToken = r
	info.CreatedTime = now
	info.ExpiredTime = now.Add(lifetime)
	info.MustChangePassword = user.MustChangePassword
	info.IsSuperAdmin = user.IsSuperAdmin
	fresh := refreshInfo(info, old.ExpiredTime)
	old.Rotated = true
	old.RotatedAt = now
	old.Disabled = true
	next := SessionRecord{TokenInfo: info, AccessToken: a, RefreshExpiredTime: old.ExpiredTime}
	oj, _ := json.Marshal(old)
	aj, _ := json.Marshal(info)
	rj, _ := json.Marshal(fresh)
	sj, _ := json.Marshal(next)
	result, err := rotateSession.Run(ctx, global.Redis, []string{PrefixRefreshToken + token, PrefixSession + old.SessionId, PrefixAccessToken + session.AccessToken, PrefixAccessToken + a, PrefixRefreshToken + r}, raw, sraw, oj, aj, rj, sj, lifetime.Milliseconds(), remaining.Milliseconds()).Int()
	if err != nil || result != 1 {
		return nil, errors.New("refreshTokenNotExist")
	}
	return tokenResponse(a, r, info, remaining), nil
}

const deleteSessionLua = `
local function remove(id)
 local key='session:'..id
 local raw=redis.call('GET',key)
 if raw then
  local s=cjson.decode(raw)
  redis.call('DEL','accessToken:'..s.access_token,'refreshToken:'..s.refresh_token,key)
  redis.call('SREM','userSessions:'..s.user_id,id)
  redis.call('SREM','tenantSessions:'..s.tenant_id,id)
 end
end
`

// Redis Lua cjson turns empty arrays into objects when encoding a table. An
// empty authorized tenant list must remain decodable as a Go slice.
const encodeSessionLua = `
local function encodeSession(record)
 if type(record.authorized_tenant_ids)=='table' and next(record.authorized_tenant_ids)==nil then
  record.authorized_tenant_ids=cjson.null
 end
 return cjson.encode(record)
end
`

var revokeSession = redis.NewScript(deleteSessionLua + `remove(ARGV[1]); return 1`)

func RevokeSession(ctx context.Context, sid string) error {
	return revokeSession.Run(ctx, global.Redis, nil, sid).Err()
}

var revokeUserSessions = redis.NewScript(encodeSessionLua + deleteSessionLua + `
local ids=redis.call('SMEMBERS',KEYS[1])
for _,id in ipairs(ids) do
 local raw=redis.call('GET','session:'..id)
 if raw then
  local s=cjson.decode(raw)
  if s.access_token ~= ARGV[1] then remove(id)
  else
   s.auth_version=tonumber(ARGV[2]);s.must_change_password=false
   redis.call('SET','session:'..id,encodeSession(s),'KEEPTTL')
   for _,key in ipairs({'accessToken:'..s.access_token,'refreshToken:'..s.refresh_token}) do
    local data=redis.call('GET',key)
    if data then local t=cjson.decode(data);t.auth_version=tonumber(ARGV[2]);t.must_change_password=false;redis.call('SET',key,encodeSession(t),'KEEPTTL') end
   end
  end
 else redis.call('SREM',KEYS[1],id) end
end
return 1
`)

func RevokeAllUserTokens(userID string) error {
	if userID == "" {
		return nil
	}
	return revokeUserSessions.Run(context.Background(), global.Redis, []string{PrefixUserSessions + userID}, "", 0).Err()
}
func RevokeAllUserTokensExcept(userID, keepAccess string, expectedVersion uint64) error {
	if keepAccess == "" {
		return RevokeAllUserTokens(userID)
	}
	data, err := global.Redis.Get(context.Background(), PrefixAccessToken+keepAccess).Bytes()
	if err != nil {
		return err
	}
	var info model.TokenInfo
	if json.Unmarshal(data, &info) != nil || info.UserId != userID {
		return errors.New("accessTokenNotExist")
	}
	user, err := LoadSessionPrincipal(global.DB, userID, info.TenantId)
	if err != nil {
		return err
	}
	if expectedVersion == 0 || user.AuthVersion != expectedVersion {
		return errors.New("accessTokenNotExist")
	}
	return revokeUserSessions.Run(context.Background(), global.Redis, []string{PrefixUserSessions + userID}, keepAccess, expectedVersion).Err()
}

var setSessionDisabled = redis.NewScript(encodeSessionLua + `
local raw=redis.call('GET',KEYS[1]);if not raw then return 0 end
local s=cjson.decode(raw);s.disabled=(ARGV[1]=='true')
redis.call('SET',KEYS[1],encodeSession(s),'KEEPTTL')
return 1
`)

func SetSessionDisabled(ctx context.Context, sid string, disabled bool) error {
	result, err := setSessionDisabled.Run(ctx, global.Redis, []string{PrefixSession + sid}, fmt.Sprint(disabled)).Int()
	if err != nil {
		return err
	}
	if result != 1 {
		return errors.New("accessTokenNotExist")
	}
	return nil
}

// CleanUserTokenSet only cleans legacy indexes; new session indexes expire and are pruned on listing.
func CleanUserTokenSet(userID, kind, prefix string) error {
	return global.Redis.Del(context.Background(), kind+userID).Err()
}
