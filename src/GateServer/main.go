package main

import (
)
import "GateServer/Server"

func main() {
	Server.GateServer = Server.NewTcpServer()
	Server.GateServer.Start()
}