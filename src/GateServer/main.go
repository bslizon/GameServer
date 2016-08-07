package main

import (
	"GateServer/Server"
	"os/signal"
	"os"
	//"runtime/pprof"
	//"GateServer/config"
	//"time"
	"syscall"
	"gameLog"
)

func main() {
	//f, _ := os.Create(config.PROFILE_FILE)
	//pprof.StartCPUProfile(f)  // 开始cpu profile，结果写到文件f中
	//defer pprof.StopCPUProfile()  // 结束profile
	Server.GateServer = Server.NewTcpPackServer()
	go Server.GateServer.Start()

	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT)
	select {
	case <- c:
		gameLog.Info("server stop.")
	}
}