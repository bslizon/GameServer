package main

import (
	"GateServer/Server"
	//"os"
	//"runtime/pprof"
	//"GateServer/config"
	//"time"
)

func main() {
	//f, _ := os.Create(config.PROFILE_FILE)
	//pprof.StartCPUProfile(f)  // 开始cpu profile，结果写到文件f中
	//defer pprof.StopCPUProfile()  // 结束profile
	Server.GateServer = Server.NewTcpPackServer()
	Server.GateServer.Start()

	//select {
	//case <-time.After(30 * time.Second):
	//}
}