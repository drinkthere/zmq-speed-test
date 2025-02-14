package container

import (
	"sync"
)

type BroadcastChannel struct {
	subscribers []chan struct{} // 或者使用实际的消息类型
	mu          sync.RWMutex
}

func NewBroadcastChannel() *BroadcastChannel {
	return &BroadcastChannel{
		subscribers: make([]chan struct{}, 0),
	}
}

func (bc *BroadcastChannel) Subscribe() chan struct{} {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	ch := make(chan struct{}, 1) // 使用缓冲通道避免阻塞
	bc.subscribers = append(bc.subscribers, ch)
	return ch
}

func (bc *BroadcastChannel) Unsubscribe(ch chan struct{}) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	for i, subscriber := range bc.subscribers {
		if subscriber == ch {
			bc.subscribers = append(bc.subscribers[:i], bc.subscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

func (bc *BroadcastChannel) Broadcast() {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	for _, subscriber := range bc.subscribers {
		// 非阻塞发送
		select {
		case subscriber <- struct{}{}:
		default:
			// 如果channel已满，跳过
		}
	}
}
