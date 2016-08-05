package main

import (
	GateServer "GateServer/Server"
	"GateServer/pack"
)

var testChan chan *pack.Pack

func main() {
	gtSvr := GateServer.NewTcpServer()
	gtSvr.Start()
}