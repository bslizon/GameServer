package Server

import (
	"net"
	GSConfig "GateServer/config"
	gLog "gameLog"
	"GateServer/types"
	rg "GateServer/randomIdGenerator"
	"sync"
	"GateServer/unPackRingBuffer"
	. "GateServer/config"
	"utils"
	"encoding/binary"
	"GateServer/pack"
	"bytes"
	"errors"
	"time"
)

type TcpServer struct {
	sync.Mutex
	connMap map[types.IdType]*net.TCPConn
	ReadPackChan chan *pack.Pack
	WritePackChan chan *pack.Pack
}

func New() *TcpServer {
	defer utils.PrintPanicStack()
	svr := new(TcpServer)
	svr.connMap = make(map[types.IdType]*net.TCPConn)
	svr.ReadPackChan = make(chan *pack.Pack, PACK_CHAN_SIZE)
	svr.WritePackChan = make(chan *pack.Pack, PACK_CHAN_SIZE)
	return svr
}


func (t *TcpServer) PutConn(i types.IdType, c *net.TCPConn) error {
	defer utils.PrintPanicStack()
	t.Lock()
	defer func() {
		t.Unlock()
	}()

	if len(t.connMap) < MAX_TCP_CONN {
		t.connMap[i] = c
		return nil
	} else {
		return errors.New("tcp conn limit")
	}
}

func (t *TcpServer) GetConn(i types.IdType) (*net.TCPConn, bool) {
	defer utils.PrintPanicStack()
	t.Lock()
	c, ok := t.connMap[i]
	t.Unlock()
	return c, ok
}

func (t *TcpServer) DelConn(i types.IdType) {
	defer utils.PrintPanicStack()
	t.Lock()
	defer func() {
		t.Unlock()
	}()
	delete(t.connMap, i)
}

func (t *TcpServer) KickConn(i types.IdType) {
	defer utils.PrintPanicStack()
	conn, ok := t.GetConn(i)
	if ok {
		t.DelConn(i)
		conn.Close()
	}
}

func (svr *TcpServer) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":" + GSConfig.EXTERNAL_LISTEN_PORT)
	if err != nil {
		gLog.Fatal(err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		gLog.Fatal(err)
	}
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			gLog.Warn(err)
			continue
		}

		id  := rg.Get()
		err = svr.PutConn(id, tcpConn)
		if err != nil {
			tcpConn.Close()
			gLog.Warn(err)
			continue
		}
		gLog.Info("connected : " + tcpConn.RemoteAddr().String())
		go tcpHandle(svr, id, tcpConn)
	}
}

func tcpHandle(svr *TcpServer, id types.IdType, conn *net.TCPConn) {
	defer func() {
		gLog.Info("disconnected :" + conn.RemoteAddr().String())
		conn.Close()
		svr.DelConn(id)
	}()
	defer utils.PrintPanicStack()


	sizeBuf := make([]byte, PACK_DATA_SIZE_TYPE_LEN)

	var dataSize32 int32
	var dataSize64 int64
	var realRdIdx int64
	var realWtIdx int64
	rbuf := new(unPackRingBuffer.UnPackRingBuffer)

	// 一开始先读一口再说
	err := conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
	if err != nil {
		panic(err)
	}
	n,err := conn.Read(rbuf.Buf[:])
	if err != nil {
		panic(err)
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
					err := conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := conn.Read(rbuf.Buf[realWtIdx : ])
					if(err != nil) {
						panic(err)
					}

					rbuf.WtIdx += int64(n)

				} else if realRdIdx > realWtIdx {//间插情况
					err := conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := conn.Read(rbuf.Buf[realWtIdx : realRdIdx])
					if(err != nil) {
						panic(err)
					}

					rbuf.WtIdx += int64(n)
				}
			}
		}

		err := binary.Read(bytes.NewBuffer(sizeBuf), binary.BigEndian, &dataSize32)
		if err != nil {
			panic(err)
		}
		if dataSize32 > MAX_PACK_DATA_SIZE {
			panic("pack data out of limit")
		}else if dataSize32 <= 0 {
			panic("pack data less than or equal 0")
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
				p := pack.New(b)

				svr.ReadPackChan <- p
				break
			} else {// 没有，得填充缓冲区
				if realRdIdx <= realWtIdx { //顺序情况
					err := conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := conn.Read(rbuf.Buf[realWtIdx : ])
					if(err != nil) {
						panic(err)
					}

					rbuf.WtIdx += int64(n)

				} else if realRdIdx > realWtIdx {//间插情况
					err := conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
					if err != nil {
						panic(err)
					}
					n, err := conn.Read(rbuf.Buf[realWtIdx : realRdIdx])
					if(err != nil) {
						panic(err)
					}

					rbuf.WtIdx += int64(n)
				}
			}
		}
	}
}