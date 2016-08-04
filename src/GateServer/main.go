package main

import (
	GateServer "GateServer/Server"
)

func main() {
	gtSvr := GateServer.New()
	go func() {
		for _ = range gtSvr.ReadPackChan {
		}
	}()
	gtSvr.Start()
}