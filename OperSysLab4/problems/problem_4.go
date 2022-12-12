package problems

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/natefinch/npipe.v2"
)

const CopyingServerPipeName string = `\\.\pipe\copying_server_pipe`

type Problem4 struct {
}

func (p *Problem4) Init(args interface{}) {
}

func (p *Problem4) Run() {
	go func() {
		err := p.server()
		fatalOnErr("copying server: ", err)
	}()
	var src, dst string
	fmt.Println("Enter src file:")
	fmt.Fscan(os.Stdin, &src)
	fmt.Println("Enter dst file name:")
	fmt.Fscan(os.Stdin, &dst)
	err := p.client(src, dst)
	fatalOnErr("copying client: ", err)
}

func (p *Problem4) client(src, dst string) error {
	pipe, err := npipe.Dial(CopyingServerPipeName)
	if err != nil {
		return err
	}
	fmt.Fprintf(pipe, "%s | %s\n", src, dst)
	r := bufio.NewReader(pipe)
	response, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("Response: %s\n", response)
	return nil
}

func (p *Problem4) server() error {
	listener, err := npipe.Listen(CopyingServerPipeName)
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
		splittedMsg := strings.Split(msg, " | ")
		src, dst := splittedMsg[0], splittedMsg[1]
		err = exec.Command("powershell.exe", "cp", src, dst).Start()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Fprint(con, "copied successfully\n")
	}
}

func (p *Problem4) Description() string {
	return "Copying server"
}
