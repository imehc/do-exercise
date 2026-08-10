package system

import "github.com/imehc/do-exercise/server/global/shared"

// notifySessionRevoked 通过 SSE 通知指定用户的全部在线连接：会话已吊销，需重新登录。
//
// 全局 Manager 在 router 启动时才初始化（配置迁移等启动早期为 nil），
// 此时静默跳过，不影响主流程；发送为快照式非阻塞投递，失败也无需回滚业务状态。
func notifySessionRevoked(userID string, reason string) {
	if shared.SSEManager == nil {
		return
	}
	_ = shared.SSEManager.SendToUser(userID, "session_revoked", reason)
}
