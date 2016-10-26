package socketIdGenerator

import (
	"GateServer/config"
	"sync"
)

var nowId config.SocketIdType
var nowMutex sync.Mutex

//默认id类型的最大值为向全连接发送广播
func Get() config.SocketIdType {
	nowMutex.Lock()
	nowId++
	id := nowId
	nowMutex.Unlock()
	return id
}