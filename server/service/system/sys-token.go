package system

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
)

type SysTokenService struct{}
type TokenScope struct {
	TenantId     string
	IsSuperAdmin bool
}

func (s TokenScope) visible(tid string) bool {
	return s.IsSuperAdmin && s.TenantId == global.PlatformTenantID || s.TenantId != "" && s.TenantId == tid
}

func (s *SysTokenService) FindAll(scope TokenScope) ([]response.SysTokenLogRsp, error) {
	ctx := context.Background()
	var ids []string
	if scope.IsSuperAdmin && scope.TenantId == global.PlatformTenantID {
		iter := global.Redis.Scan(ctx, 0, util.PrefixSession+"*", 100).Iterator()
		for iter.Next(ctx) {
			ids = append(ids, strings.TrimPrefix(iter.Val(), util.PrefixSession))
		}
		if err := iter.Err(); err != nil {
			return nil, err
		}
	} else {
		if scope.TenantId == "" {
			return nil, errors.New("invalidTenant")
		}
		var err error
		ids, err = global.Redis.SMembers(ctx, util.PrefixTenantSessions+scope.TenantId).Result()
		if err != nil {
			return nil, err
		}
	}
	rows := make([]response.SysTokenLogRsp, 0, len(ids))
	for _, id := range ids {
		session, err := util.ReadSession(ctx, id)
		if err != nil {
			_ = global.Redis.SRem(ctx, util.PrefixTenantSessions+scope.TenantId, id).Err()
			continue
		}
		if !scope.visible(session.TenantId) {
			continue
		}
		if _, err := util.ValidateSessionVersion(global.DB, session.UserId, session.TenantId, session.AuthVersion, session.TenantVersion); err != nil {
			continue
		}
		rows = append(rows, response.SysTokenLogRsp{SessionId: id, UserId: session.UserId, TenantId: session.TenantId, Username: session.Username, Disabled: session.Disabled, AccessTokenCreated: session.CreatedTime, AccessTokenExpired: session.ExpiredTime, RefreshTokenCreated: session.CreatedTime, RefreshTokenExpired: session.RefreshExpiredTime})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].AccessTokenCreated.After(rows[j].AccessTokenCreated) })
	return rows, nil
}
func (s *SysTokenService) Delete(req request.SysTokenDeleteReq, scope TokenScope) error {
	session, err := util.ReadSession(context.Background(), req.SessionId)
	if err != nil || !scope.visible(session.TenantId) {
		return errors.New("accessTokenNotExist")
	}
	return util.RevokeSession(context.Background(), req.SessionId)
}
func (s *SysTokenService) ModityStatus(req request.SysTokenModityStatusReq, scope TokenScope) error {
	session, err := util.ReadSession(context.Background(), req.SessionId)
	if err != nil || !scope.visible(session.TenantId) || req.Disabled == nil {
		return errors.New("accessTokenNotExist")
	}
	if _, err := util.ValidateSessionVersion(global.DB, session.UserId, session.TenantId, session.AuthVersion, session.TenantVersion); err != nil {
		return err
	}
	return util.SetSessionDisabled(context.Background(), req.SessionId, *req.Disabled)
}
