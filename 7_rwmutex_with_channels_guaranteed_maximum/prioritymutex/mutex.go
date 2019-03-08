package prioritymutex

import (
	"sync"
	"sync/atomic"
)

// this is problematic implementation, don't use it
// the implementation with drainLowPriority(maxToDrain) is better

type PriorityMutex struct {
	highPriorityChannel chan bool
	numHighPriorityWaiting int32
	lowPriorityChannel chan bool
	highPriorityMutexReleased chan bool
	waitGroup sync.WaitGroup
}

func NewPriorityMutex() *PriorityMutex {
	m := &PriorityMutex{
		highPriorityChannel: make(chan bool),
		lowPriorityChannel: make(chan bool),
		highPriorityMutexReleased: make(chan bool),
	}
	go m.mainLoop()
	return m
}

func (m *PriorityMutex) HighPriorityLock() {
	atomic.AddInt32(&m.numHighPriorityWaiting, 1)
	m.highPriorityChannel <- true
	atomic.AddInt32(&m.numHighPriorityWaiting, -1)
}

func (m *PriorityMutex) HighPriorityUnlock() {
	m.highPriorityMutexReleased <- true
}

func (m *PriorityMutex) LowPriorityRLock() {
	m.lowPriorityChannel <- true
}

func (m *PriorityMutex) LowPriorityRUnlock() {
	m.waitGroup.Done()
}

func (m *PriorityMutex) mainLoop() {
	for {

		// this added select will guarantee that the high priority wait the minimum time possible
		select {
		case <- m.highPriorityChannel:
			<- m.highPriorityMutexReleased
		default:
		}

		select {
		case <- m.highPriorityChannel:
			<- m.highPriorityMutexReleased
		case <- m.lowPriorityChannel:
			m.waitGroup.Add(1)
			m.drainLowPriority()
			m.waitGroup.Wait()
		}

	}
}

func (m *PriorityMutex) drainLowPriority() {
	for {

		if atomic.LoadInt32(&m.numHighPriorityWaiting) > 0 {
			return
		}

		select {
		case <- m.lowPriorityChannel:
			m.waitGroup.Add(1)
		default:
			return
		}

	}
}