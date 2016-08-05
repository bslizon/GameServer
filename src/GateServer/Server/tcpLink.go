package Server

import (
	"net"
	"time"
	"GateServer/unPackRingBuffer"
	"encoding/binary"
	"bytes"
	. "GateServer/config"
	"GateServer/pack"
	"utils"
	gLog "gameLog"
	"strconv"
	"fmt"
	"io"
)

type TcpLink struct {
	sid        SocketIdType
	server     *TcpServer
	conn       *net.TCPConn
	WtSyncChan chan []byte
}

func NewTcpLink(sid SocketIdType, svr *TcpServer, co *net.TCPConn) *TcpLink {
	lk := new(TcpLink)
	lk.sid = sid
	lk.server = svr
	lk.conn = co
	lk.WtSyncChan = make(chan []byte, WRITE_PACK_SYNC_CHAN_SIZE)
	return lk
}

func (lk *TcpLink) Close() {
	defer utils.PrintPanicStack()
	lk.conn.Close()
	gLog.Info("disconnected: " + lk.conn.RemoteAddr().String() + " socketid: " + fmt.Sprintf("%d", lk.sid) + " " + " mapCount: " + strconv.Itoa(len(lk.server.linkMap)))
	close(lk.WtSyncChan)
}

func (lk *TcpLink) StartRead() {
	defer func() {
		lk.server.KickLink(lk.sid)
	}()
	defer utils.PrintPanicStack()

	sizeBuf := make([]byte, PACK_DATA_SIZE_TYPE_LEN)

	var dataSize32 int32
	var dataSize64 int64
	var realRdIdx int64
	var realWtIdx int64
	rbuf := new(unPackRingBuffer.UnPackRingBuffer)

	// 一开始先读一口再说
	err := lk.conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
	if err != nil {
		panic(err)
	}
	n,err := lk.conn.Read(rbuf.Buf[:])
	if err != nil {
		if err == io.EOF {
			gLog.Info("tcp read EOF")
			return
		} else {
			panic(err)
		}
	}
	rbuf.WtIdx += int64(n)

	for {
		// 判断是否已经存在一个整的size在缓冲区中
		for {
			realRdIdx = rbuf.RdIdx % RING_BUFFER_SIZE
			realWtIdx = rbuf.WtIdx % RING_BUFFER_SIZE

			if rbuf.WtIdx - rbuf.RdIdx >= PACK_DATA_SIZE_TYPE_LEN {// 有，取出来
				realEndIdx := (rbuf.RdIdx + PACK_DATA_SIZE_TYPE_LEN) % RING_BUFFER_SIZE
				if realRdIdx < realEndIdx {
					copy(sizeBuf, rbuf.Buf[realRdIdx : realEndIdx])
				} else {
					copy(sizeBuf, rbuf.Buf[realRdIdx : ])
					copy(sizeBuf[RING_BUFFER_SIZE - realRdIdx: ], rbuf.Buf[ : PACK_DATA_SIZE_TYPE_LEN - (RING_BUFFER_SIZE - realRdIdx)])
				}
				rbuf.RdIdx += PACK_DATA_SIZE_TYPE_LEN
				break
			} else {//没有，继续从conn读
				if realRdIdx <= realWtIdx { //顺序情况
					err := lk.conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := lk.conn.Read(rbuf.Buf[realWtIdx : ])
					if(err != nil) {
						if err == io.EOF {
							gLog.Info("tcp read EOF")
							return
						} else {
							panic(err)
						}
					}

					rbuf.WtIdx += int64(n)

				} else if realRdIdx > realWtIdx {//间插情况
					err := lk.conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := lk.conn.Read(rbuf.Buf[realWtIdx : realRdIdx])
					if(err != nil) {
						if err == io.EOF {
							gLog.Info("tcp read EOF")
							return
						} else {
							panic(err)
						}
					}

					rbuf.WtIdx += int64(n)
				}
			}
		}

		err := binary.Read(bytes.NewBuffer(sizeBuf), binary.BigEndian, &dataSize32)
		if err != nil {
			panic(err)
		}
		if dataSize32 > MAX_INBOUND_PACK_DATA_SIZE {
			panic("read pack data out of limit")
		}else if dataSize32 <= 0 {
			panic("read pack data less than or equal 0")
		}

		dataSize64 = int64(dataSize32)

		// 判断是否已经存在一个整的packdata在缓冲区中
		for {
			realRdIdx = rbuf.RdIdx % RING_BUFFER_SIZE
			realWtIdx = rbuf.WtIdx % RING_BUFFER_SIZE

			//已经有了，解析一个packdata出来
			if rbuf.WtIdx - rbuf.RdIdx >= dataSize64 {
				b := make([]byte, dataSize64)
				realEndIdx := (rbuf.RdIdx + dataSize64) % RING_BUFFER_SIZE
				if realRdIdx < realEndIdx {
					copy(b, rbuf.Buf[realRdIdx : realEndIdx])
				} else {
					copy(b, rbuf.Buf[realRdIdx : ])
					copy(b[RING_BUFFER_SIZE - realRdIdx: ], rbuf.Buf[ : dataSize64 - (RING_BUFFER_SIZE - realRdIdx)])
				}

				rbuf.RdIdx += dataSize64
				RouteIn(pack.NewPack(lk.sid, b))
				break
			} else {// 没有，得填充缓冲区
				if realRdIdx <= realWtIdx { //顺序情况
					err := lk.conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := lk.conn.Read(rbuf.Buf[realWtIdx : ])
					if(err != nil) {
						if err == io.EOF {
							gLog.Info("tcp read EOF")
							return
						} else {
							panic(err)
						}
					}

					rbuf.WtIdx += int64(n)

				} else if realRdIdx > realWtIdx {//间插情况
					err := lk.conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := lk.conn.Read(rbuf.Buf[realWtIdx : realRdIdx])
					if(err != nil) {
						if err == io.EOF {
							gLog.Info("tcp read EOF")
							return
						} else {
							panic(err)
						}
					}

					rbuf.WtIdx += int64(n)
				}
			}
		}
	}
}

func (lk *TcpLink) StartWrite() {
	defer func() {
		lk.server.KickLink(lk.sid)
	}()
	defer utils.PrintPanicStack()

	var wCount int
	var rawDataSize int
	for rawData := range lk.WtSyncChan {
		n := len(rawData)
		if n <= 0 {
			panic("write pack data less than or equal 0")
		} else if n > MAX_OUTBOUND_PACK_DATA_SIZE{
			panic("write pack data out of limit")
		}
		rawDataSize = PACK_DATA_SIZE_TYPE_LEN + n
		data := make([]byte, rawDataSize)
		binary.BigEndian.PutUint32(data, uint32(n))
		copy(data[PACK_DATA_SIZE_TYPE_LEN:], rawData)

		wCount = 0
		for {
			err := lk.conn.SetReadDeadline(time.Now().Add(TCP_WRITE_TIMEOUT * time.Second))
			if err != nil {
				panic(err)
			}

			wn, err := lk.conn.Write(data)
			if(err != nil) {
				panic(err)
			}

			wCount += wn
			if wCount == rawDataSize {
				break
			} else if wCount > rawDataSize {
				panic("write byte count error")
			}
		}
	}
}
