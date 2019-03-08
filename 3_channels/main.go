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

	clientsChannel := make(chan bool)
	commitsChannel := make(chan bool)
	go mainLoop(commitsChannel, clientsChannel)

	go commitBlockLoop(commitsChannel)

	group := sync.WaitGroup{}
	for i := 0; i<NUM_CLIENTS; i++ {
		group.Add(1)
		go clientAddingTransaction(clientsChannel, i, &group)
	}
	group.Wait()
	time.Sleep(1 * time.Millisecond)

}

func mainLoop(commitsChannel chan bool, clientsChannel chan bool) {
	for {

		select {
		case <- commitsChannel:
			time.Sleep(1 * time.Millisecond)
		case <- clientsChannel:
			time.Sleep(1 * time.Millisecond)
		}

	}
}

func commitBlockLoop(commitsChannel chan bool) {
	for {
		tryingToCommitStart := time.Now()
		commitsChannel <- true

		now := time.Now()
		fmt.Printf("\ncommit loop scheduled, waited %v\n", now.Sub(tryingToCommitStart))
		time.Sleep(1 * time.Millisecond)
	}
}

func clientAddingTransaction(clientsChannel chan bool, index int, group *sync.WaitGroup) {
	clientsChannel <- true
	fmt.Printf("  %d", index)
	time.Sleep(1 * time.Millisecond)

	group.Done()
}