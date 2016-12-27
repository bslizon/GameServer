package GlobalConfig

import (
	logLevel "GameServer/gameLog/level"
)

const (
	LOG_LEVEL                 = logLevel.INFO
	GATESERVER_LOG_FILE_PATH  = "./GateServerLog.txt"
	LOGICSERVER_LOG_FILE_PATH = "./LogicServerLog.txt"
	USE_LOG_FILE              = true
)
