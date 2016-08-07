package Server

import (
	"net"
	gLog "gameLog"
	g "GateServer/socketIdGenerator"
	"sync"
	"GateServer/config"
	"utils"
	"errors"
	"fmt"
)

type TcpPackServer struct {
	sync.RWMutex
	linkMap map[config.SocketIdType]*tcpPackLink
}

func NewTcpPackServer() *TcpPackServer {
	svr := new(TcpPackServer)
	svr.linkMap = make(map[config.SocketIdType]*tcpPackLink)
	return svr
}

func (svr *TcpPackServer) PutLink(i config.SocketIdType, lk *tcpPackLink) (err error) {
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

	if len(svr.linkMap) > config.MAX_TCP_CONN {
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

func (svr *TcpPackServer) GetLink(i config.SocketIdType) (*tcpPackLink, bool) {
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
func (svr *TcpPackServer) RemoveLink(i config.SocketIdType) {
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	lk, ok := svr.GetLink(i)
	if ok {
		svr.Lock()
		defer func() {
			svr.Unlock()
		}()
		delete(svr.linkMap, i)
		gLog.Info(fmt.Sprintf("has been removed: %s sid: %d mapCount: %d ", lk.conn.RemoteAddr().String(), lk.sid, len(lk.server.linkMap)))
		lk.Close()
	}
}

// 复制一份linkMap，用于广播
func (svr *TcpPackServer) GetLinkMapCopy() map[config.SocketIdType]*tcpPackLink {
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	linkMap := make(map[config.SocketIdType]*tcpPackLink)
	svr.RWMutex.Lock()
	defer func() {
		svr.RWMutex.Unlock()
	}()
	for k, v := range svr.linkMap {
		linkMap[k] = v
	}
	return linkMap
}

func (svr *TcpPackServer) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":" + config.EXTERNAL_LISTEN_PORT)
	if err != nil {
		gLog.Fatal(err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		gLog.Fatal(err)
	}
	defer tcpListener.Close()

	gLog.Info("listen on: " + config.EXTERNAL_LISTEN_PORT)

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			gLog.Warn(err)
			continue
		}

		gLog.Info(fmt.Sprintf("connected: %s mapCount: %d ", tcpConn.RemoteAddr().String(), len(svr.linkMap)))
		go svr.handleTcpConn(tcpConn)
	}
}

func (svr *TcpPackServer) handleTcpConn(tcpConn *net.TCPConn) {
	defer utils.PrintPanicStack()
	////////////////////////////////////////////////////////////////////

	sid  := g.Get()

	lk := NewPackLink(sid, svr, tcpConn)
	err := svr.PutLink(sid, lk)
	if err != nil {
		lk.Close()
		gLog.Warn(fmt.Sprintf("%s disconnected: %s sid: %d mapCount: %d ", err.Error(), lk.conn.RemoteAddr().String(), lk.sid, len(lk.server.linkMap)))
		return
	}

	go lk.StartReadPack()
	go lk.StartWritePack()
	gLog.Info(fmt.Sprintf("serving: %s sid: %d mapCount: %d ", tcpConn.RemoteAddr().String(),  lk.sid, len(lk.server.linkMap)))
}