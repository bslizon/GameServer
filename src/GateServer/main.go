package main

import (
	"GateServer/Server"
)

func main() {
	Server.GateServer = Server.NewTcpPackServer()
	Server.GateServer.Start()
}