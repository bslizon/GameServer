package utils

import (
	gLog "GameServer/gameLog"
	"fmt"
	"runtime"
)

func PrintPanicStack() {
	if x := recover(); x != nil {
		var stackString string
		i := 3
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			stackString += fmt.Sprintf("[func:%s %s:%d] ", runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		switch value := x.(type) {
		case error:
			gLog.Panic(value.Error() + stackString)
		case string:
			gLog.Panic(value + stackString)
		default:
			gLog.Printf("[PANIC] unknown panic: %#v.%s", value, stackString)
		}
	}
}
