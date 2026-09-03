package util

import (
	"fmt"
	"os"
)

// Exit 打印错误信息并以非零码退出。
// message 按 Printf 格式处理；err 为 nil 时不参与格式化，
// 否则「没有占位符的提示语 + nil」会被 Printf 追加成 %!(EXTRA <nil>)。
func Exit(message string, err error) {
	if err == nil {
		fmt.Print(message)
	} else {
		fmt.Printf(message, err)
	}
	os.Exit(1)
}
