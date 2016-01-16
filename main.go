package main

func main() {
	agent := Agent{
		Addr: "localhost",
		Port: "7023",
	}
	agent.Run()
}
