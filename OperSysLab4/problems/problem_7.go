package problems

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Problem7 struct {
	arr []int

	mu *sync.Mutex
}

type Problem7Args struct {
	Arr []int
}

func (p *Problem7) Init(args interface{}) {
	if argsStruct, ok := args.(Problem7Args); !ok {
		log.Fatal("casting to problem 7 args")
	} else {
		p.arr = argsStruct.Arr
	}
	p.mu = new(sync.Mutex)
}

func (p *Problem7) Run() {
	waitChan := make(chan struct{})
	go p.addDelProcess()
	go p.sortAndPrintProcess()
	<-waitChan
}

func (p *Problem7) addDelProcess() {
	rand.Seed(time.Now().UnixNano())
	for range time.Tick(time.Millisecond * 500) {
		p.mu.Lock()
		p.arr = append(p.arr, rand.Int()%100-50)
		p.arr = append(p.arr, rand.Int()%100-50)

		idx := rand.Int() % len(p.arr)
		p.arr = append(p.arr[:idx], p.arr[idx+1:]...)
		p.mu.Unlock()
	}
}

func (p *Problem7) sortAndPrintProcess() {
	for range time.Tick(time.Second * 3) {
		p.mu.Lock()
		sort.Slice(p.arr, func(i, j int) bool {
			return p.arr[i] < p.arr[j]
		})
		for _, el := range p.arr {
			fmt.Printf("%d ", el)
		}
		fmt.Println()
		p.mu.Unlock()
	}
}

func (p *Problem7) Description() string {
	return "Parallel sort and +/-"
}
