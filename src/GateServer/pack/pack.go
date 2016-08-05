package pack

import "GateServer/types"

//数据包结构
type Pack struct {
	SocketId types.IdType
	Data []byte
}

func New(id types.IdType, b []byte) *Pack {
	p := new(Pack)
	p.SocketId = id
	p.Data = b
	return p
}