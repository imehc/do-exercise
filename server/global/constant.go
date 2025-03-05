package global

import "fmt"

const (
	CLAIMS = "claims"
)

// 获取缓存key
func GetCacheKey(username string) (key string) {
	key = fmt.Sprintf("auth:%s", username)
	return
}
