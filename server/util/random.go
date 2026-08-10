package util

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomNumber 生成指定位数的随机数字字符串。
// 用 crypto/rand 而非时间种子的 math/rand（CWE-338），
// 该值用于找回密码与邮箱登录验证码。
func GenerateRandomNumber(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			// crypto/rand 失败属不可恢复的系统级错误，回退到 '0' 不现实，
			// 直接 panic 让上层兜底（fail fast）
			panic("crypto/rand unavailable: " + err.Error())
		}
		result[i] = digits[n.Int64()]
	}
	return string(result)
}
