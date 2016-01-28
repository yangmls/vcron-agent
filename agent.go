package main

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
	"net"
	"os"
)

type Agent struct {
	Name string
	Addr string
	Port string
	C    net.Conn
}

func (agent *Agent) Run() {
	err := agent.Connect()

	if err != nil {
		fmt.Println("can not connect to", agent.Addr, "on port", agent.Port)
		return
	}

	agent.Waiting()
}

func (agent *Agent) Connect() error {
	conn, err := net.Dial("tcp", agent.Addr+":"+agent.Port)

	agent.C = conn
	return err
}

func (agent *Agent) Waiting() {
	for {
		request, err := agent.WaitRequest()

		if err != nil {
			break
		}

		response := &vcron.Response{
			Result: true,
		}

		if request.Type == "register" {
			response.Message = agent.Name
		}

		if request.Type == "run" {
			for _, job := range request.Jobs {
				go RunCommand(job.Command)
			}
		}

		agent.SendResponse(response)
	}
}

func (agent *Agent) WaitRequest() (*vcron.Request, error) {
	var (
		n   int
		err error
	)

	fmt.Println("waiting request")

	prefix := make([]byte, 4, 4)

	if n, err = agent.C.Read(prefix); err != nil {
		return nil, err
	}

	fmt.Println(prefix)

	if n == 0 {
		return nil, nil
	}

	if n != 4 {
		return nil, nil
	}

	var (
		size    uint64
		errcode int
	)

	if size, errcode = binary.Uvarint(prefix); errcode <= 0 {
		return nil, nil
	}

	buf := make([]byte, int(size), int(size))

	if n, err = agent.C.Read(buf); err != nil {
		return nil, err
	}

	if uint64(n) != size {
		return nil, nil
	}

	fmt.Println("got request")

	if err != nil {
		return nil, err
	}

	request := &vcron.Request{}

	err = proto.Unmarshal(buf, request)

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (agent *Agent) SendResponse(response *vcron.Response) {
	data, _ := proto.Marshal(response)
	agent.C.Write(data)
}

func NewAgent(name string, addr string, port string) *Agent {
	if name == "" {
		name = defaultName()
	}

	agent := &Agent{
		Name: name,
		Addr: addr,
		Port: port,
	}

	return agent
}

func defaultName() (name string) {
	name, err := os.Hostname()

	if err != nil {
		return ""
	}

	return name
}
