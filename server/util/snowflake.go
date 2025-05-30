package util

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/spf13/cast"
)

var (
	node     *snowflake.Node
	nodeOnce sync.Once
)

// InitSnowflake 初始化雪花节点（建议在项目启动时调用一次）
func InitSnowflake(nodeID int64) error {
	var err error
	nodeOnce.Do(func() {
		// 检查 nodeID 合法性（0~1023）
		if nodeID < 0 || nodeID > 1023 {
			err = fmt.Errorf("nodeID 必须在 0~1023 之间，当前为: %d", nodeID)
			return
		}
		node, err = snowflake.NewNode(nodeID)
		if err != nil {
			fmt.Printf("Snowflake 节点初始化失败: %v\n", err)
		} else {
			fmt.Printf("Snowflake 节点初始化成功，nodeID: %d\n", nodeID)
		}
	})
	return err
}

// NextID 生成全局唯一ID（限制在16位以内）
func NextID() string {
	if node == nil {
		panic("Snowflake 节点未初始化，请先调用 InitSnowflake")
	}
	return cast.ToString(node.Generate().Int64())
}
