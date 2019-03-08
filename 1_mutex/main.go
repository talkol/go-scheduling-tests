package main

import (
	"sync"
	"time"
	"fmt"
	"runtime"
)

const NUM_CLIENTS = 5000

func main() {

	runtime.GOMAXPROCS(1) // play with values from 1 to 4

	mutex := sync.Mutex{}

	go commitBlockLoop(&mutex)

	group := sync.WaitGroup{}
	for i := 0; i<NUM_CLIENTS; i++ {
		group.Add(1)
		go clientAddingTransaction(&mutex, i, &group)
	}
	group.Wait()
	time.Sleep(1 * time.Millisecond)

}

func commitBlockLoop(mutex *sync.Mutex) {
	for {
		tryingToCommitStart := time.Now()
		mutex.Lock()

		now := time.Now()
		fmt.Printf("\ncommit loop scheduled, waited %v\n", now.Sub(tryingToCommitStart))
		time.Sleep(1 * time.Millisecond)

		mutex.Unlock()
	}
}

func clientAddingTransaction(mutex *sync.Mutex, index int, group *sync.WaitGroup) {
	mutex.Lock()
	fmt.Printf("  %d", index)
	time.Sleep(1 * time.Millisecond)
	mutex.Unlock()

	group.Done()
}