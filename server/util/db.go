package util

import (
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