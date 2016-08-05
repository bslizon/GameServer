package socketIdGenerator

import (
	"GateServer/config"
	"sync"
)

var nowId config.SocketIdType
var nowMutex sync.Mutex

func Get() config.SocketIdType {
	nowMutex.Lock()
	nowId++
	id := nowId
	nowMutex.Unlock()
	return id
}