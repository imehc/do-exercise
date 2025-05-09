package util

import (
	"encoding/json"
	"sync"
)

// Client 表示一个SSE客户端连接
type Client struct {
	ID      string
	Channel chan []byte
}

// Message 表示要发送的消息结构
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// Manager SSE连接管理器
type Manager struct {
	clients    map[string]*Client
	mutex      sync.RWMutex
	broadcasts chan Message
}

// NewManager 创建一个新的SSE管理器
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[string]*Client),
		broadcasts: make(chan Message, 100),
	}
}

// AddClient 添加新的客户端连接
func (m *Manager) AddClient(userID string) *Client {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client := &Client{
		ID:      userID,
		Channel: make(chan []byte, 10),
	}

	m.clients[userID] = client
	return client
}

// RemoveClient 移除客户端连接
func (m *Manager) RemoveClient(userID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if client, exists := m.clients[userID]; exists {
		close(client.Channel)
		delete(m.clients, userID)
	}
}

// SendToUser 发送消息给指定用户
func (m *Manager) SendToUser(userID string, msgType string, content interface{}) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, exists := m.clients[userID]
	if !exists {
		return nil
	}

	msg := Message{Type: msgType, Content: content}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	client.Channel <- data
	return nil
}

// SendToAllExcept 发送消息给除了指定用户之外的所有用户
func (m *Manager) SendToAllExcept(excludeUserID string, msgType string, content interface{}) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	msg := Message{Type: msgType, Content: content}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for userID, client := range m.clients {
		if userID != excludeUserID {
			client.Channel <- data
		}
	}
	return nil
}

// BroadcastToAll 广播消息给所有用户
func (m *Manager) BroadcastToAll(msgType string, content interface{}) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	msg := Message{Type: msgType, Content: content}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, client := range m.clients {
		client.Channel <- data
	}
	return nil
}
