package Server

import (
	"net"
	"time"
	"encoding/binary"
	"bytes"
	. "GateServer/config"
	"GateServer/pack"
	"utils"
	gLog "gameLog"
	"fmt"
	"io"
	"errors"
)

type tcpPackLink struct {
	sid        SocketIdType
	server     *TcpPackServer
	conn       *net.TCPConn
	wtSyncChan chan []byte
}

func NewPackLink(sid SocketIdType, svr *TcpPackServer, co *net.TCPConn) *tcpPackLink {
	lk := new(tcpPackLink)
	lk.sid = sid
	lk.server = svr
	lk.conn = co
	lk.wtSyncChan = make(chan []byte, WRITE_PACK_SYNC_CHAN_SIZE)
	return lk
}

func (lk *tcpPackLink) Close() {
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	lk.conn.Close()
	gLog.Info(fmt.Sprintf("disconnected: %s sid: %d mapCount: %d ", lk.conn.RemoteAddr().String(), lk.sid, len(lk.server.linkMap)))
	close(lk.wtSyncChan)
}

func (lk *tcpPackLink) PutBytes(b []byte) (err error) {
	// panic转error
	defer func() {
		if x := recover(); x != nil {
			switch value := x.(type) {
			case error:
				err = value
			case string:
				err = errors.New(value)
			default:
				err = errors.New(fmt.Sprintf("unknown panic: %#v. ", value))
			}
		}
	}()
	////////////////////////////////////////////////////////////////////

	select {
	case lk.wtSyncChan <- b:
		err = nil
		return
	case <-time.After(time.Second * WRITE_PACK_SYNC_CHAN_TIMEOUT):
		err = errors.New("put wtSyncChan timeout.")
		return
	}

	////////////////////////////////////////////////////////////////////
	err = nil
	return
}


func (lk *tcpPackLink) StartReadPack() {
	defer func() {
		lk.server.RemoveLink(lk.sid)
	}()
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////


	sizeBuf := make([]byte, PACK_DATA_SIZE_TYPE_LEN)
	var dataSize PackDataSizeType
	var sizeBufIdx PackDataSizeType
	var dataBufIdx PackDataSizeType
	for {
		sizeBufIdx = 0
		for {
			err := lk.conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
			if err != nil {
				panic(err)
			}
			n,err := lk.conn.Read(sizeBuf[sizeBufIdx:])
			if err != nil {
				if err == io.EOF {
					gLog.Info(fmt.Sprintf("tcp read EOF. sid: %d ", lk.sid))
					return
				} else {
					panic(err)
				}
			}
			sizeBufIdx += PackDataSizeType(n)
			if sizeBufIdx == PACK_DATA_SIZE_TYPE_LEN {
				break
			} else if sizeBufIdx < PACK_DATA_SIZE_TYPE_LEN && sizeBufIdx > 0 {
				continue
			} else {
				panic("sizeBufIdx error.")
			}
		}

		err := binary.Read(bytes.NewBuffer(sizeBuf), binary.BigEndian, &dataSize)
		if err != nil {
			panic(err)
		}
		if dataSize > MAX_INBOUND_PACK_DATA_SIZE {
			panic("read pack data out of limit.")
		}else if dataSize <= 0 {
			panic("read pack data less than or equal 0.")
		}

		dataBufIdx = 0
		data := make([]byte, dataSize)
		for {
			err := lk.conn.SetReadDeadline(time.Now().Add(TCP_READ_TIMEOUT * time.Second))
			if err != nil {
				panic(err)
			}
			n,err := lk.conn.Read(data[dataBufIdx:])
			if err != nil {
				if err == io.EOF {
					gLog.Info(fmt.Sprintf("tcp read EOF. sid: %d ", lk.sid))
					return
				} else {
					panic(err)
				}
			}

			dataBufIdx += PackDataSizeType(n)

			if dataBufIdx == dataSize {
				lk.server.RoutePackIn(pack.NewPack(lk.sid, data))
				break
			} else if dataBufIdx < dataSize && dataBufIdx > 0 {
				continue
			} else {
				panic("dataBufIdx error.")
			}
		}

	}
}

func (lk *tcpPackLink) StartWritePack() {
	defer func() {
		lk.server.RemoveLink(lk.sid)
	}()
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	var wCount int
	var rawDataSize int
	for rawData := range lk.wtSyncChan {
		n := len(rawData)
		if n <= 0 {
			panic("write pack data less than or equals 0.")
		} else if n > MAX_OUTBOUND_PACK_DATA_SIZE {
			panic("write pack data out of limit.")
		}
		rawDataSize = PACK_DATA_SIZE_TYPE_LEN + n
		data := make([]byte, rawDataSize)
		binary.BigEndian.PutUint32(data, uint32(n))
		copy(data[PACK_DATA_SIZE_TYPE_LEN:], rawData)

		wCount = 0
		for {
			err := lk.conn.SetWriteDeadline(time.Now().Add(TCP_WRITE_TIMEOUT * time.Second))
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
				panic("write byte count error.")
			}
		}
	}
}
