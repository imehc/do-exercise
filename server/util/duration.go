package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDurationString 解析表示时间长度的字符串并将其转换为time.Duration类型。
// 该函数支持解析格式如"1d"（1天）、"2h"（2小时）、"3m"（3分钟）和"4s"（4秒）。
// 如果字符串格式不符合预期，函数将返回一个错误。
// 参数:
//
//	durationStr - 表示时间长度的字符串。
//
// 返回值:
//
//	time.Duration - 解析后的时间长度。
//	error - 如果解析失败，返回错误信息。
func ParseDurationString(durationStr string) (time.Duration, error) {
	// 使用正则表达式匹配输入字符串，以确定其是否符合预期的格式。
	re := regexp.MustCompile(`(\d+)([dhms])`)
	matches := re.FindStringSubmatch(durationStr)

	// 如果匹配成功且包含两个捕获组（数字和单位）：
	if len(matches) == 3 {
		// 将捕获的数字部分转换为整数。
		value, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, err
		}
		// 获取时间单位，并转换为小写以统一处理。
		unit := strings.ToLower(matches[2])
		// 根据不同的时间单位，将数字转换为对应的time.Duration值。
		switch unit {
		case "d":
			return time.Duration(value) * 24 * time.Hour, nil
		case "h":
			return time.Duration(value) * time.Hour, nil
		case "m":
			return time.Duration(value) * time.Minute, nil
		case "s":
			return time.Duration(value) * time.Second, nil
		default:
			// 如果时间单位不识别，返回错误。
			return 0, fmt.Errorf("invalid time unit: %s", unit)
		}
	}

	// 如果输入字符串不符合预期格式但能直接转换为整数，则假设其单位为秒。
	if value, err := strconv.Atoi(durationStr); err == nil {
		return time.Duration(value) * time.Second, nil // 假设纯数字以秒为单位
	}

	// 如果输入字符串既不符合预期格式也无法直接转换为整数，返回错误。
	return 0, fmt.Errorf("invalid input format: %s", durationStr)
}
