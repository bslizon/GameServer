package unPackRingBuffer

import "GateServer/config"

type UnPackRingBuffer struct {
	RdIdx int64
	WtIdx int64
	Buf [config.RING_BUFFER_SIZE]byte
}