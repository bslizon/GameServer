package main

import (
	"GateServer/Server"
)

func main() {
	Server.GateServer = Server.NewTcpServer()
	Server.GateServer.Start()
}