package main

import (
	"GateServer/pack"
	"GateServer/Server"
)

var testChan chan *pack.Pack

func main() {
	gtSvr := Server.NewTcpServer()
	gtSvr.Start()
}