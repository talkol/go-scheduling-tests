package prioritymutex

type PriorityMutex struct {
	highPriorityChannel chan bool
	lowPriorityChannel chan bool
	mutexReleased chan bool
}

func NewPriorityMutex() *PriorityMutex {
	m := &PriorityMutex{
		highPriorityChannel: make(chan bool),
		lowPriorityChannel: make(chan bool),
		mutexReleased: make(chan bool),
	}
	go m.mainLoop()
	return m
}

func (m *PriorityMutex) HighPriorityLock() {
	m.highPriorityChannel <- true
}

func (m *PriorityMutex) HighPriorityUnlock() {
	m.mutexReleased <- true
}

func (m *PriorityMutex) LowPriorityLock() {
	m.lowPriorityChannel <- true
}

func (m *PriorityMutex) LowPriorityUnlock() {
	m.mutexReleased <- true
}

func (m *PriorityMutex) mainLoop() {
	for {

		// this added select will guarantee that the high priority wait the minimum time possible
		select {
		case <- m.highPriorityChannel:
			<- m.mutexReleased
		default:
		}

		select {
		case <- m.highPriorityChannel:
			<- m.mutexReleased
		case <- m.lowPriorityChannel:
			<- m.mutexReleased
		}

	}
}