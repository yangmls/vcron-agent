package main

import (
	"fmt"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
)

type Agent struct {
	Name string
	Addr string
	Port string
}

func (agent *Agent) Run() {
	if agent.Name == "" {
		agent.Name = defaultName()
	}

	conn, err := net.Dial("tcp", agent.Addr+":"+agent.Port)

	if err != nil {
		fmt.Println("can not connect to" + agent.Addr + ":" + agent.Port)
		return
	}

	defer conn.Close()

	agent.Register(conn)

	for {
		handle(conn)
	}
}

func (agent *Agent) Register(conn net.Conn) {
	message := &vcron.Message{
		Type: proto.String("register"),
		Name: &agent.Name,
	}

	data, _ := proto.Marshal(message)

	fmt.Println(data)

	conn.Write(data)
}

func defaultName() (name string) {
	name, err := os.Hostname()

	if err != nil {
		return ""
	}

	return name
}

func handle(conn net.Conn) {
	data := make([]byte, 4096)
	len, readErr := conn.Read(data)

	if readErr != nil {
		return
	}

	message := &vcron.Message{}
	uncodeErr := proto.Unmarshal(data[0:len], message)

	if uncodeErr != nil {
		fmt.Println(uncodeErr)
		return
	}

	if message.GetType() == "run" {
		go RunCommand(message.GetCommand())
	}
}
