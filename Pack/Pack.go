package Pack

import (
	"GameServer/GateServer/config"
)

//数据包结构
type Pack struct {
	Sid  config.SocketIdType
	Data []byte
}

func NewPack(id config.SocketIdType, b []byte) *Pack {
	p := new(Pack)
	p.Sid = id
	p.Data = b
	return p
}
