package main

import (
	"sync"
	"time"
	"fmt"
	"runtime"
	"github.com/talkol/scheduling-tests/6_rwmutex_with_channels/prioritymutex"
)

const NUM_CLIENTS = 5000

func main() {

	runtime.GOMAXPROCS(1) // play with values from 1 to 4

	mutex := prioritymutex.NewPriorityMutex()

	go commitBlockLoop(mutex)

	group := sync.WaitGroup{}
	for i := 0; i<NUM_CLIENTS; i++ {
		group.Add(1)
		go clientAddingTransaction(mutex, i, &group)
	}
	group.Wait()
	time.Sleep(1 * time.Millisecond)

}

func commitBlockLoop(mutex *prioritymutex.PriorityMutex) {
	for {
		tryingToCommitStart := time.Now()
		mutex.HighPriorityLock()

		now := time.Now()
		fmt.Printf("\ncommit loop scheduled, waited %v\n", now.Sub(tryingToCommitStart))
		time.Sleep(1 * time.Millisecond)

		mutex.HighPriorityUnlock()
	}
}

func clientAddingTransaction(mutex *prioritymutex.PriorityMutex, index int, group *sync.WaitGroup) {
	mutex.LowPriorityRLock()
	fmt.Printf("  %d", index)
	time.Sleep(1 * time.Millisecond)
	mutex.LowPriorityRUnlock()

	group.Done()
}
