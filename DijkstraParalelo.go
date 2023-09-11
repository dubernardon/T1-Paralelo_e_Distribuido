package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const N = 5 // número de filósofos (e garfos)

type State int

const (
	THINKING State = iota
	HUNGRY
	EATING
)

var (
	state              [N]State
	bothForksAvailable [N]chan bool
	criticalRegionMtx  sync.Mutex
	outputMtx          sync.Mutex
)

func left(i int) int {
	return (i - 1 + N) % N
}

func right(i int) int {
	return (i + 1) % N
}

func myRand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func test(i int) {
	if state[i] == HUNGRY && state[left(i)] != EATING && state[right(i)] != EATING {
		state[i] = EATING
		bothForksAvailable[i] <- true
	}
}

func think(i int) {
	duration := myRand(400, 800)
	outputMtx.Lock()
	defer outputMtx.Unlock()
	fmt.Printf("%d is thinking %dms\n", i, duration)
	time.Sleep(time.Millisecond * time.Duration(duration))
}

func takeForks(i int) {
	criticalRegionMtx.Lock()
	defer criticalRegionMtx.Unlock()
	state[i] = HUNGRY
	outputMtx.Lock()
	fmt.Printf("\t\t%d is HUNGRY\n", i)
	outputMtx.Unlock()
	test(i)
	criticalRegionMtx.Unlock()
	<-bothForksAvailable[i]
}

func eat(i int) {
	duration := myRand(400, 800)
	outputMtx.Lock()
	defer outputMtx.Unlock()
	fmt.Printf("\t\t\t\t%d is eating %dms\n", i, duration)
	time.Sleep(time.Millisecond * time.Duration(duration))
}

func putForks(i int) {
	criticalRegionMtx.Lock()
	defer criticalRegionMtx.Unlock()
	state[i] = THINKING
	test(left(i))
	test(right(i))
}

func philosopher(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		think(i)
		takeForks(i)
		eat(i)
		putForks(i)
	}
}

func main() {
	fmt.Println("dp_14")

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		bothForksAvailable[i] = make(chan bool)
		wg.Add(1)
		go philosopher(i, &wg)
	}
	wg.Wait()
}
