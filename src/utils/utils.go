package utils

import (
	gLog "gameLog"
	"runtime"
)


func PrintPanicStack() {
	if x := recover(); x != nil {
		switch value := x.(type) {
		case error:
			gLog.Panic(value.Error())
		case string:
			gLog.Panic(value)
		default:
			gLog.Printf("[PANIC] unknown panic: %#v.", value)
		}

		i := 3
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			gLog.Printf("[func:%v, file:%v, line:%v]", runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}
	}
}
