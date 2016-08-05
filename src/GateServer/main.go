package main

import (
	"GateServer/pack"
	"GateServer/Server"
	"time"
)

var testChan chan *pack.Pack

func main() {
	gtSvr := Server.NewTcpServer()
	go gtSvr.Start()

	time.Sleep(10 * time.Second)

	lk, ok := gtSvr.GetLink(1)
	if ok {
		lk.WtSyncChan <- []byte("abcdefghijklmn")
	}
	select {

	}
}