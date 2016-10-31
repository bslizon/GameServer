package gameLog

import (
	"log"
	"os"
	"fmt"
	"GlobalConfig"
	logLevel "gameLog/level"
	"bufio"
	"sync"
)

var lgrMtx sync.Mutex	// 解决同时操作文件指针和logger的竞争问题
var logFileWriter *bufio.Writer
var logger *log.Logger

func init() {
	if GlobalConfig.USE_LOG_FILE {
		expFilePtr, err := os.OpenFile(GlobalConfig.GATESERVER_LOG_FILE_PATH, os.O_CREATE | os.O_APPEND, 0600)
		if err != nil {
			fmt.Println(err)
			return
		}
		logFileWriter = bufio.NewWriter(expFilePtr)
		logger = log.New(logFileWriter, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	} else {
		logger = log.New(os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	}
}

func Debug(v interface{}) {
	if logLevel.DEBUG >= GlobalConfig.LOG_LEVEL {
		lgrMtx.Lock()
		defer lgrMtx.Unlock()
		logger.Output(2, fmt.Sprintln("[DEBUG]", v))
	}
}

func Info(v interface{}) {
	if logLevel.INFO >= GlobalConfig.LOG_LEVEL {
		lgrMtx.Lock()
		defer lgrMtx.Unlock()
		logger.Output(2, fmt.Sprintln("[INFO]", v))
	}
}

func Warn(v interface{}) {
	if logLevel.WARN >= GlobalConfig.LOG_LEVEL {
		lgrMtx.Lock()
		defer lgrMtx.Unlock()
		logger.Output(2, fmt.Sprintln("[WARN]", v))
	}
}

func Error(v interface{}) {
	if logLevel.ERROR >= GlobalConfig.LOG_LEVEL {
		lgrMtx.Lock()
		defer lgrMtx.Unlock()
		logger.Output(2, fmt.Sprintln("[ERROR]", v))
		if logFileWriter != nil{
			logFileWriter.Flush()
		}
	}
}

func Panic(v interface{}) {
	if logLevel.PANIC >= GlobalConfig.LOG_LEVEL {
		lgrMtx.Lock()
		defer lgrMtx.Unlock()
		logger.Output(5, fmt.Sprintln("[PANIC]", v))
		if logFileWriter != nil{
			logFileWriter.Flush()
		}
	}
}

func Fatal(v interface{}) {
	if logLevel.FATAL >= GlobalConfig.LOG_LEVEL {
		lgrMtx.Lock()
		defer lgrMtx.Unlock()
		logger.Output(2, fmt.Sprintln("[FATAL]", v))
		if logFileWriter != nil{
			logFileWriter.Flush()
		}
		os.Exit(1)
	}
}

func Printf(format string, v ...interface{}){
	lgrMtx.Lock()
	defer lgrMtx.Unlock()
	logger.Output(5, fmt.Sprintf(format, v...))
}

func Flush() error {
	lgrMtx.Lock()
	defer lgrMtx.Unlock()
	return logFileWriter.Flush()
}