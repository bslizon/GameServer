package randomIdGenerator

import (
	"GateServer/types"
	"sync"
)

var nowId types.IdType
var nowMutex sync.Mutex

func Get() types.IdType {
	nowMutex.Lock()
	nowId++
	id := nowId
	nowMutex.Unlock()
	return id
}