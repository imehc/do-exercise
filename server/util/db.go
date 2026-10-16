package util

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// DB 返回当前请求携带上下文（用户ID、取消信号等）的数据库连接。
// 若在非请求上下文（如定时任务、异步协程）调用，则返回全局基础连接 global.DB。
func DB(c *gin.Context) *gorm.DB {
	if c != nil {
		if v, exists := c.Get(global.ContextDBKey); exists {
			if db, ok := v.(*gorm.DB); ok && db != nil {
				return db
			}
		}
	}
	return global.DB
}

// BypassTenantDB 返回绕过租户隔离的请求级 DB。
// 平台层管理跨租户数据（创建租户、租户管理员、跨租户查询）时使用；
// 调用方必须显式设置目标数据的 tenant_id，否则将落入空租户。
//
// 派生自 DB(c) 已有的语句上下文，而不是裸的 c.Request.Context()：
// 后者不带 ContextMiddleware 注入的 userId / tenantId / isSuperAdmin，
// 一旦覆盖掉，模型钩子写入的 created_by / updated_by 会全部变成空串，
// 服务层的超管判定也会退化成「受限身份」。
func BypassTenantDB(c *gin.Context) *gorm.DB {
	base := DB(c)
	if c == nil {
		return base
	}
	parent := c.Request.Context()
	if base.Statement != nil && base.Statement.Context != nil {
		parent = base.Statement.Context
	}
	ctx := context.WithValue(parent, global.ContextTenantBypassKey, true)
	return base.WithContext(ctx)
}
