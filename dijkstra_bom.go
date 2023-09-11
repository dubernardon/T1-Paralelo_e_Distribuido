package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var take chan struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func delay(maxRand int) {
	time.Sleep(50*time.Millisecond + time.Duration(rand.Intn(maxRand)))
}

type fork struct {
	sync.Mutex
	id int
}

var output_mutex sync.Mutex

func (f *fork) String() string {
	return fmt.Sprintf("fork %d", f.id)
}

func newForks(n int) []*fork {
	c := make([]*fork, n)
	for i := 0; i < n; i++ {
		c[i] = &fork{id: i}
	}
	return c
}

type philosopher struct {
	id int
	left, right *fork
}

func newPhilosopher(fs []*fork, id, left, right int) *philosopher {
	return &philosopher{id: id, left: fs[left], right: fs[right]}
}

func (p *philosopher) think(iteration int) {
	time_to_think := 2

	// think
	output_mutex.Lock()
	for s := 1; s <= p.id*5; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" T%d\n", iteration) // Thinking
	output_mutex.Unlock()
	time.Sleep(time.Duration(time_to_think) * time.Second)
}

func (p *philosopher) eat(iteration int) {

	time_to_eat := 3

	take <- struct{}{}
	p.left.Lock()
	p.right.Lock()

	output_mutex.Lock()
	for s := 1; s <= p.id*5; s++ {
		fmt.Printf(" ")
	}
    fmt.Printf(" E%d\n", iteration) // Eating
	output_mutex.Unlock()

    time.Sleep(time.Duration(time_to_eat) * time.Second)

	p.left.Unlock()
	p.right.Unlock()
	<-take
}

func main() {
	rounds := 5
	fs := newForks(5)
	ps := []*philosopher{
		newPhilosopher(fs, 0, 0, 1),
		newPhilosopher(fs, 1, 1, 2),
		newPhilosopher(fs, 2, 2, 3),
		newPhilosopher(fs, 3, 3, 4),
		newPhilosopher(fs, 4, 4, 0),
	}
	take = make(chan struct{}, len(ps)-1)
	wg := &sync.WaitGroup{}

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]\n")


	for i := 1; i <= rounds; i++ {
		wg.Add(len(ps))

		for _, p := range ps {
			go func(p *philosopher, i int) {
				p.think(i)
				p.eat(i)
				wg.Done()
			}(p, i)
		}
		wg.Wait()
	}

	
}
