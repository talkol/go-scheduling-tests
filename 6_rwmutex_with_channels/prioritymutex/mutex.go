package prioritymutex

import "sync"

type PriorityMutex struct {
	highPriorityChannel chan bool
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
	m.highPriorityChannel <- true
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
			m.drainLowPriority(50)
			m.waitGroup.Wait()
		}

	}
}

func (m *PriorityMutex) drainLowPriority(maxToDrain int) {
	for i:=0; i<maxToDrain; i++ {

		select {
		case <- m.lowPriorityChannel:
			m.waitGroup.Add(1)
		default:
			return
		}

	}
}