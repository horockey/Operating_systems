package problems

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/natefinch/npipe.v2"
)

const CryptoServerPipeName string = `\\.\pipe\crypto_server_pipe`

type Problem5 struct {
}

func (p *Problem5) Init(args interface{}) {
}

func (p *Problem5) Run() {
	go func() {
		err := p.server()
		fatalOnErr("crypto server: ", err)
	}()
	var src string
	fmt.Print("Enter src text: ")
	fmt.Fscanf(os.Stdin, "\n")
	r := bufio.NewReader(os.Stdin)
	src, _ = r.ReadString('\n')
	err := p.client(src)
	fatalOnErr("crypto client: ", err)
}

func (p *Problem5) client(src string) error {
	pipe, err := npipe.Dial(CryptoServerPipeName)
	if err != nil {
		return err
	}
	fmt.Fprintf(pipe, "%s\n", src)
	r := bufio.NewReader(pipe)
	response, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("Response: %s\n", response)
	return nil
}

func (p *Problem5) server() error {
	listener, err := npipe.Listen(CryptoServerPipeName)
	if err != nil {
		return err
	}
	for {
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
		encryptedMsg := ""
		for idx := range msg {
			encryptedMsg += string(rune(msg[idx]) + 1)
		}
		fmt.Fprintf(con, "%s\n", encryptedMsg)
	}
}

func (p *Problem5) Description() string {
	return "Crypto server"
}
