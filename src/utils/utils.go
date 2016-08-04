package utils

import (
	gLog "gameLog"
)


func PrintPanicStack() {
	if x := recover(); x != nil {
		gLog.Panic(x)
		//i := 3
		//funcName, file, line, ok := runtime.Caller(i)
		//for ok {
		//	log.Printf("[func:%v, file:%v, line:%v]\n",runtime.FuncForPC(funcName).Name(), file, line)
		//	i++
		//	funcName, file, line, ok = runtime.Caller(i)
		//}
	}
}
