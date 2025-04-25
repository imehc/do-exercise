package util

import (
	"math/rand"
	"time"
)

// GenerateRandomNumber 生成指定位数的随机数字字符串
func GenerateRandomNumber(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const digits = "0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = digits[r.Intn(len(digits))]
	}

	return string(result)
}
