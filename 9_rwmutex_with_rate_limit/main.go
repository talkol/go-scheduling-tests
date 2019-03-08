package main

import (
	"sync"
	"time"
	"fmt"
	"runtime"
	"github.com/talkol/scheduling-tests/9_rwmutex_with_rate_limit/ratelimiter"
)

const NUM_CLIENTS = 5000

func main() {

	runtime.GOMAXPROCS(1) // play with values from 1 to 4

	mutex := sync.RWMutex{}
	rateLimiter := ratelimiter.NewRateLimiter(50)

	go commitBlockLoop(&mutex)

	group := sync.WaitGroup{}
	for i := 0; i<NUM_CLIENTS; i++ {
		group.Add(1)
		go clientAddingTransaction(rateLimiter, &mutex, i, &group)
	}
	group.Wait()
	time.Sleep(1 * time.Millisecond)

}

func commitBlockLoop(mutex *sync.RWMutex) {
	for {
		tryingToCommitStart := time.Now()
		mutex.Lock()

		now := time.Now()
		fmt.Printf("\ncommit loop scheduled, waited %v\n", now.Sub(tryingToCommitStart))
		time.Sleep(1 * time.Millisecond)

		mutex.Unlock()
	}
}

func clientAddingTransaction(rateLimiter *ratelimiter.RateLimiter, mutex *sync.RWMutex, index int, group *sync.WaitGroup) {
	rateLimiter.RequestSlot()
	defer rateLimiter.ReleaseSlot()

	mutex.RLock()
	fmt.Printf("  %d", index)
	time.Sleep(1 * time.Millisecond)
	mutex.RUnlock()

	group.Done()
}