package util

import (
	"encoding/json"
	"sync"
)

// Client 表示一个SSE客户端连接
type Client struct {
	UserID  string // 所属用户
	ConnID  string // 同一用户的不同连接（多标签页/多设备）需各自独立
	Channel chan []byte
}

// Message 表示要发送的消息结构
type Message struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

// Manager SSE连接管理器
type Manager struct {
	// userID -> connID -> client
	clients map[string]map[string]*Client
	mutex   sync.RWMutex
}

// NewManager 创建一个新的SSE管理器
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]map[string]*Client),
	}
}

// AddClient 为用户新增一条连接，返回该连接专属的 Client。
//
// 同一用户允许多条连接。早期实现用 map[userID]*Client，第二个标签页会直接覆盖
// 第一个且不关闭其 channel：前者的 handler goroutine 永久泄漏，
// 而任一连接断开时的 RemoveClient 又会误关另一条仍在使用的 channel。
func (m *Manager) AddClient(userID, connID string) *Client {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client := &Client{
		UserID:  userID,
		ConnID:  connID,
		Channel: make(chan []byte, 10),
	}

	if m.clients[userID] == nil {
		m.clients[userID] = make(map[string]*Client)
	}
	m.clients[userID][connID] = client
	return client
}

// RemoveClient 移除指定连接。
//
// 刻意不 close(client.Channel)：发送在锁外进行（见 snapshot），
// 关闭 channel 会与并发发送形成 "send on closed channel" panic。
// 读取方 handler 依靠请求上下文取消退出，不依赖 channel 关闭信号；
// 连接移出 map 后其 channel 即成为垃圾，由 GC 回收。
func (m *Manager) RemoveClient(userID, connID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	conns, ok := m.clients[userID]
	if !ok {
		return
	}
	delete(conns, connID)
	if len(conns) == 0 {
		delete(m.clients, userID)
	}
}

// ConnCount 返回当前连接总数，便于监控泄漏
func (m *Manager) ConnCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := 0
	for _, conns := range m.clients {
		total += len(conns)
	}
	return total
}

// snapshot 在锁内拷贝目标 channel 列表。
//
// 发送必须在锁外进行：早期实现持 RLock 对容量 10 的 channel 做阻塞发送，
// 一个慢消费者填满缓冲即可让发送方卡在读锁上，
// 随后 RemoveClient（写锁）排在其后，整个 SSE 子系统连带卡死。
func (m *Manager) snapshot(include func(userID string) bool) []chan []byte {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	chans := make([]chan []byte, 0, len(m.clients))
	for userID, conns := range m.clients {
		if include != nil && !include(userID) {
			continue
		}
		for _, c := range conns {
			chans = append(chans, c.Channel)
		}
	}
	return chans
}

// deliver 非阻塞投递：缓冲已满说明该消费者滞后，丢弃这条消息而不是阻塞发送方
func deliver(chans []chan []byte, data []byte) {
	for _, ch := range chans {
		select {
		case ch <- data:
		default:
		}
	}
}

// SendToUser 发送消息给指定用户的全部连接
func (m *Manager) SendToUser(userID string, msgType string, content any) error {
	data, err := json.Marshal(Message{Type: msgType, Content: content})
	if err != nil {
		return err
	}
	deliver(m.snapshot(func(id string) bool { return id == userID }), data)
	return nil
}

// SendToAllExcept 发送消息给除了指定用户之外的所有用户
func (m *Manager) SendToAllExcept(excludeUserID string, msgType string, content any) error {
	data, err := json.Marshal(Message{Type: msgType, Content: content})
	if err != nil {
		return err
	}
	deliver(m.snapshot(func(id string) bool { return id != excludeUserID }), data)
	return nil
}

// BroadcastToAll 广播消息给所有用户
func (m *Manager) BroadcastToAll(msgType string, content any) error {
	data, err := json.Marshal(Message{Type: msgType, Content: content})
	if err != nil {
		return err
	}
	deliver(m.snapshot(nil), data)
	return nil
}
