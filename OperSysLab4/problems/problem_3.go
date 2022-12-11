package problems

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Problem3 struct {
	stationCapacity int
	refillMinSec    int
	refillMaxSec    int
	newCarMinSec    int
	newCarMaxSec    int

	refillSec int
	newCarSec int

	carsOnService int

	carsServed int
	carsPassed int

	mu *sync.RWMutex
}

type Problem3Args struct {
	StationCapacity int
	RefillMinSec    int
	RefillMaxSec    int
	NewCarMinSec    int
	NewCarMaxSec    int
}

func (p *Problem3) Init(args interface{}) {
	if argsStruct, ok := args.(Problem3Args); !ok {
		log.Fatal("casting to problem 3 args")
	} else {
		p.stationCapacity = argsStruct.StationCapacity
		p.refillMinSec = argsStruct.RefillMinSec
		p.refillMaxSec = argsStruct.RefillMaxSec
		p.newCarMinSec = argsStruct.NewCarMinSec
		p.newCarMaxSec = argsStruct.NewCarMaxSec
	}
	rand.Seed(time.Now().UnixNano())
	p.refillSec = p.refillMinSec + rand.Int()%(p.refillMaxSec-p.refillMinSec)
	p.newCarSec = p.newCarMinSec + rand.Int()%(p.newCarMaxSec-p.newCarMinSec)
	p.mu = new(sync.RWMutex)
}

func (p *Problem3) Run() {
	fmt.Printf("New car interval: %ds\nRefill interval %ds\n", p.newCarSec, p.refillSec)

	newCarTicker := time.NewTicker(time.Duration(p.newCarSec) * time.Second)
	for range newCarTicker.C {
		p.mu.Lock()
		if p.carsOnService < p.stationCapacity {
			p.carsOnService++
			p.mu.Unlock()
			go func() {
				refillTimer := time.NewTimer(time.Duration(p.refillSec) * time.Second)
				<-refillTimer.C
				p.mu.Lock()
				p.carsServed++
				p.carsOnService--
				p.mu.Unlock()
			}()
		} else {
			p.carsPassed++
			p.mu.Unlock()
		}
		fmt.Printf(`======
New car!
Cars on service: %d,
Cars served: %d
Cars passed: %d
`, p.carsOnService, p.carsServed, p.carsPassed)
	}
}

func (p *Problem3) Description() string {
	return fmt.Sprintf("Auto-station modelling")
}
