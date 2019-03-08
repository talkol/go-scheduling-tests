package limitedrwmutex

import "sync"

type LimitedRWMutex struct {
	sync.RWMutex
	channel chan struct{}
}

func NewLimitedRWMutex() *LimitedRWMutex {
	m := &LimitedRWMutex{
		channel: make(chan struct{}, 50),
	}
	return m
}

func (m *LimitedRWMutex) Lock() {
	m.RWMutex.Lock()
}

func (m *LimitedRWMutex) Unlock() {
	m.RWMutex.Unlock()
}

func (m *LimitedRWMutex) RLock() {
	m.channel <- struct{}{}
	m.RWMutex.RLock()
}

func (m *LimitedRWMutex) RUnlock() {
	m.RWMutex.RUnlock()
	<- m.channel
}