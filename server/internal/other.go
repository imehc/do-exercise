package internal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/imehc/do-exercise/server/util"
)

// InitOther 初始化其他工具
func InitOther() {
	nodeID := int64(1) // 默认值
	if v := os.Getenv("SNOWFLAKE_NODE_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			nodeID = id
		}
	}
	if err := util.InitSnowflake(nodeID); err != nil {
		panic(fmt.Sprintf("初始化雪花算法失败: %v", err))
	}
}
