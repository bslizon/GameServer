package gameLog

import (
	"log"
	"os"
	"fmt"
	"globalConfig"
	logLevel "gameLog/level"
	"bufio"
)

var logFileWriter *bufio.Writer
var logger *log.Logger

//log输出等级
var level = globalConfig.LOG_LEVEL

func init() {
	expFilePtr, err := os.OpenFile("E:/result.txt", os.O_CREATE | os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	logFileWriter = bufio.NewWriter(expFilePtr)
	logger = log.New(logFileWriter, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}

func Debug(v interface{}) {
	if logLevel.DEBUG >= level {
		logger.Output(2, fmt.Sprintln("[DEBUG]", v))
	}
}

func Info(v interface{}) {
	if logLevel.INFO >= level {
		logger.Output(2, fmt.Sprintln("[INFO]", v))
	}
}

func Warn(v interface{}) {
	if logLevel.WARN >= level {
		logger.Output(2, fmt.Sprintln("[WARN]", v))
	}
}

func Error(v interface{}) {
	if logLevel.ERROR >= level {
		logger.Output(2, fmt.Sprintln("[ERROR]", v))
		logFileWriter.Flush()
	}
}

func Panic(v interface{}) {
	if logLevel.PANIC >= level {
		logger.Output(5, fmt.Sprintln("[PANIC]", v))
		logFileWriter.Flush()
	}
}

func Fatal(v interface{}) {
	if logLevel.FATAL >= level {
		logger.Output(2, fmt.Sprintln("[FATAL]", v))
		logFileWriter.Flush()
		os.Exit(1)
	}
}

func Printf(format string, v ...interface{}){
	logger.Output(5, fmt.Sprintf(format, v...))
}