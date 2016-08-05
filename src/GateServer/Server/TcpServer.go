package Server

import (
	"net"
	GSConfig "GateServer/config"
	gLog "gameLog"
	g "GateServer/socketIdGenerator"
	"sync"
	. "GateServer/config"
	"utils"
	"errors"
	"fmt"
)

type TcpServer struct {
	sync.RWMutex
	linkMap map[GSConfig.SocketIdType]*TcpLink
}

func NewTcpServer() *TcpServer {
	svr := new(TcpServer)
	svr.linkMap = make(map[GSConfig.SocketIdType]*TcpLink)
	return svr
}

func (svr *TcpServer) PutLink(i GSConfig.SocketIdType, lk *TcpLink) (err error) {
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

	if len(svr.linkMap) > MAX_TCP_CONN {
		return errors.New("tcp conn limit")
	}

	if _, ok := svr.GetLink(i); ok {
		return errors.New("sid conflict")
	}

	svr.Lock()
	defer func() {
		svr.Unlock()
	}()
	svr.linkMap[i] = lk

	////////////////////////////////////////////////////////////////////
	err = nil
	return
}

func (svr *TcpServer) GetLink(i GSConfig.SocketIdType) (*TcpLink, bool) {
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	svr.RLock()
	defer func() {
		svr.RUnlock()
	}()
	c, ok := svr.linkMap[i]
	return c, ok
}

// 会关闭连接
func (svr *TcpServer) RemoveLink(i GSConfig.SocketIdType) {
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	lk, ok := svr.GetLink(i)
	if ok {
		svr.Lock()
		defer func() {
			svr.Unlock()
		}()
		delete(svr.linkMap, i)
		gLog.Info(fmt.Sprintf("has been removed: %s socketid: %d mapCount: %d ", lk.conn.RemoteAddr().String(), lk.sid, len(lk.server.linkMap)))
		lk.Close()
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

	gLog.Info("listen on: " + GSConfig.EXTERNAL_LISTEN_PORT)

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			gLog.Warn(err)
			continue
		}

		gLog.Info(fmt.Sprintf("connected: %s mapCount: %d", tcpConn.RemoteAddr().String(), len(svr.linkMap)))
		go handleTcpConn(svr, tcpConn)
	}
}

func handleTcpConn(svr *TcpServer, tcpConn *net.TCPConn) {
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	sid  := g.Get()

	lk := NewTcpLink(sid, svr, tcpConn)
	err := svr.PutLink(sid, lk)
	if err != nil {
		lk.Close()
		gLog.Warn(fmt.Sprintf("%s disconnected: %s socketid: %d mapCount: %d ", err.Error(), lk.conn.RemoteAddr().String(), lk.sid, len(lk.server.linkMap)))
		return
	}

	go lk.StartRead()
	go lk.StartWrite()
	gLog.Info(fmt.Sprintf("serving: %s socketid: %d mapCount: %d ", tcpConn.RemoteAddr().String(),  lk.sid, len(lk.server.linkMap)))
}