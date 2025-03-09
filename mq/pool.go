package mq

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ChannelPool Channel池结构体
type ChannelPool struct {
	mu      sync.RWMutex
	pool    chan *amqp.Channel
	factory func() (*amqp.Channel, error)
	closed  bool
}

// NewChannelPool 创建新的Channel池
func NewChannelPool(size int) (*ChannelPool, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid pool size")
	}

	pool := &ChannelPool{
		pool: make(chan *amqp.Channel, size),
		factory: func() (*amqp.Channel, error) {
			if conn == nil {
				return nil, fmt.Errorf("connection is nil")
			}
			return conn.Channel()
		},
	}

	// 预创建Channel
	for i := 0; i < size; i++ {
		ch, err := pool.factory()
		if err != nil {
			return nil, fmt.Errorf("failed to create channel: %v", err)
		}
		pool.pool <- ch
	}

	return pool, nil
}

// Get 获取一个Channel
func (p *ChannelPool) Get() (*amqp.Channel, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, fmt.Errorf("pool is closed")
	}
	p.mu.RUnlock()

	select {
	case ch := <-p.pool:
		if ch.IsClosed() {
			// 如果Channel已关闭，创建新的
			newCh, err := p.factory()
			if err != nil {
				return nil, err
			}
			return newCh, nil
		}
		return ch, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout getting channel from pool")
	}
}

// Put 归还一个Channel
func (p *ChannelPool) Put(ch *amqp.Channel) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return fmt.Errorf("pool is closed")
	}
	p.mu.RUnlock()

	if ch == nil {
		return fmt.Errorf("channel is nil")
	}

	if ch.IsClosed() {
		// 如果Channel已关闭，创建新的
		newCh, err := p.factory()
		if err != nil {
			return err
		}
		ch = newCh
	}

	select {
	case p.pool <- ch:
		return nil
	default:
		// 如果池已满，关闭多余的Channel
		return ch.Close()
	}
}

// Close 关闭Channel池
func (p *ChannelPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	close(p.pool)

	for ch := range p.pool {
		if err := ch.Close(); err != nil {
			return err
		}
	}

	return nil
}
