package util

import (
	"sync"
	"testing"
	"time"
)

// 同一用户的多条连接必须互相独立：
// 旧实现用 map[userID]*Client，第二条连接会覆盖第一条，
// 且任一条断开都会误伤另一条。
func TestMultipleConnectionsPerUser(t *testing.T) {
	m := NewManager()
	a := m.AddClient("u1", "conn-a")
	b := m.AddClient("u1", "conn-b")

	if a.Channel == b.Channel {
		t.Fatal("同一用户的两条连接不应共用 channel")
	}
	if got := m.ConnCount(); got != 2 {
		t.Fatalf("期望 2 条连接，实际 %d", got)
	}

	// 广播应同时送达两条
	if err := m.SendToUser("u1", "msg", "hello"); err != nil {
		t.Fatal(err)
	}
	for name, ch := range map[string]chan []byte{"conn-a": a.Channel, "conn-b": b.Channel} {
		select {
		case <-ch:
		default:
			t.Errorf("%s 未收到消息", name)
		}
	}

	// 移除 conn-a 不应影响 conn-b
	m.RemoveClient("u1", "conn-a")
	if got := m.ConnCount(); got != 1 {
		t.Fatalf("移除一条后期望剩 1 条，实际 %d", got)
	}
	if err := m.SendToUser("u1", "msg", "still here"); err != nil {
		t.Fatal(err)
	}
	select {
	case <-b.Channel:
	default:
		t.Error("conn-a 断开后 conn-b 应仍能收到消息")
	}

	m.RemoveClient("u1", "conn-b")
	if got := m.ConnCount(); got != 0 {
		t.Fatalf("期望 0 条连接，实际 %d", got)
	}
}

// 慢消费者（缓冲写满且无人读取）不得阻塞发送方。
// 旧实现持 RLock 做阻塞发送，一个滞后客户端就能让整个 SSE 子系统卡死。
func TestSlowConsumerDoesNotBlockSender(t *testing.T) {
	m := NewManager()
	m.AddClient("slow", "conn-1") // 只加不读

	done := make(chan struct{})
	go func() {
		defer close(done)
		// 远超 channel 容量（10）
		for i := 0; i < 200; i++ {
			_ = m.BroadcastToAll("msg", i)
		}
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("慢消费者阻塞了发送方——非阻塞投递未生效")
	}
}

// 发送与增删连接并发进行时不得死锁，也不得 panic。
func TestConcurrentSendAndRemove(t *testing.T) {
	m := NewManager()

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// 持续发送
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_ = m.BroadcastToAll("tick", 1)
			}
		}
	}()

	// 持续增删连接
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			connID := string(rune('a' + n))
			for j := 0; j < 500; j++ {
				m.AddClient("u", connID)
				m.RemoveClient("u", connID)
			}
		}(i)
	}

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()

	time.Sleep(500 * time.Millisecond)
	close(stop)

	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("并发收发与增删连接发生死锁")
	}

	if got := m.ConnCount(); got != 0 {
		t.Fatalf("全部移除后期望 0 条连接，实际 %d", got)
	}
}
