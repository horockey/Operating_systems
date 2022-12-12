package problems

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/npipe.v2"
)

const SeabattlePipeName string = `\\.\pipe\seabattle_pipe`

type Problem6 struct {
	// . - empty
	// S - ship
	// X - attacked
	fieldA [][]rune
	fieldB [][]rune
	shipsA int
	shipsB int

	mu *sync.Mutex
}

type Problem6Args struct {
	FieldA [][]rune
	FieldB [][]rune
	ShipsA int
	ShipsB int
}

func (p *Problem6) Init(args interface{}) {
	if argsStruct, ok := args.(Problem6Args); !ok {
		log.Fatal("casting to problem 6 args")
	} else {
		p.fieldA = argsStruct.FieldA
		p.fieldB = argsStruct.FieldB
		p.shipsA = argsStruct.ShipsA
		p.shipsB = argsStruct.ShipsB
	}
	p.mu = new(sync.Mutex)
}

func (p *Problem6) Run() {
	rand.Seed(time.Now().UnixNano())
	wg := new(sync.WaitGroup)
	wg.Add(2)
	for i := 1; i <= 2; i++ {
		go func(idx int) {
			defer wg.Done()
			var name string
			switch idx {
			case 1:
				name = "A"
			case 2:
				name = "B"
			}
			err := p.server(name)
			fatalOnErr(fmt.Sprintf("seabattle server %s: ", name), err)
		}(i)
	}
	wg.Wait()
}

func (p *Problem6) server(name string) error {
	var field [][]rune
	var ships *int
	var opponentsName string
	switch name {
	case "A":
		field = p.fieldA
		ships = &p.shipsA
		opponentsName = "B"
	case "B":
		field = p.fieldB
		ships = &p.shipsB
		opponentsName = "A"
	}
	listener, err := npipe.Listen(SeabattlePipeName + name)
	if err != nil {
		return err
	}

	if name == "A" {
		opponentsPipe, err := npipe.Dial(SeabattlePipeName + opponentsName)
		if err != nil {
			return err
		}
		fmt.Fprint(opponentsPipe, "A0\n")
	}

	for {
		p.mu.Lock()
		fmt.Printf("%s:\n", name)
		for y := 0; y < len(field); y++ {
			for x := 0; x < len(field[0]); x++ {
				ch := field[y][x]
				fmt.Printf("%c ", ch)
			}
			fmt.Println()
		}
		fmt.Println()
		p.mu.Unlock()
		con, err := listener.Accept()
		if err != nil {
			return err
		}

		r := bufio.NewReader(con)
		msg, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		x := rune(msg[0])
		xInt := int(x - 'A')
		y, _ := strconv.Atoi(strings.Trim(strings.TrimPrefix(msg, string(x)), "\n "))

		if field[y][xInt] == 'S' {
			*ships--
		}
		field[y][xInt] = 'X'
		if *ships == 0 {
			log.Fatalf("GAME OVER. %s LOST", name)
		}
		time.Sleep(time.Millisecond * 500)

		xInt = rand.Int() % 10
		x = rune(xInt) + 'A'
		y = rand.Int() % 10

		opponentsPipe, err := npipe.Dial(SeabattlePipeName + opponentsName)
		if err != nil {
			return err
		}
		fmt.Fprintf(opponentsPipe, "%c%d\n", x, y)
	}
}

func (p *Problem6) Description() string {
	return "Sea battle"
}

// func (p *Problem6) clearScreen() {
// 	cmd := exec.Command("cmd", "/c", "cls")
// 	cmd.Stdout = os.Stdout
// 	cmd.Run()
// }
