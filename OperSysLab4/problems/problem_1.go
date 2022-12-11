package problems

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type Problem1 struct {
	files []string
	mu    *sync.Mutex
}

type Problem1Args struct {
	Files []string
}

func (p *Problem1) Init(args interface{}) {
	if argsStruct, ok := args.(Problem1Args); !ok {
		log.Fatal("casting to problem 1 args")
	} else {
		p.files = argsStruct.Files
	}
	p.mu = new(sync.Mutex)
}

func (p *Problem1) Run() {
	for _, file := range p.files {
		go func(file string) {
			p.mu.Lock()
			defer p.mu.Unlock()
			text, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(`
=== START OF FILE %s ===
%s
=== END OF FILE %s ===`,
				file, string(text), file)
		}(file)
	}
	time.Sleep(time.Second)
}
