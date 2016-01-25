package main

import (
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
	fmt.Println("waiting request")
	buf := make([]byte, 2048)
	len, err := agent.C.Read(buf)
	fmt.Println("got request")

	if err != nil {
		return nil, err
	}

	data := buf[0:len]

	request := &vcron.Request{}

	err = proto.Unmarshal(data, request)

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
