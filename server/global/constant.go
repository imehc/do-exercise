package global

type contextKey string

const (
	// ContextUserIDKey 上下文用户ID
	ContextUserIDKey contextKey = "userId"
	// ContextDBKey 请求绑定的数据库连接
	ContextDBKey contextKey = "gormDB"
)
