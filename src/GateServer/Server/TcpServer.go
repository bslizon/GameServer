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
	"strconv"
)

type TcpServer struct {
	sync.RWMutex
	linkMap map[GSConfig.SocketIdType]*TcpLink
}

func NewTcpServer() *TcpServer {
	defer utils.PrintPanicStack()
	svr := new(TcpServer)
	svr.linkMap = make(map[GSConfig.SocketIdType]*TcpLink)
	return svr
}


func (svr *TcpServer) PutLink(i GSConfig.SocketIdType, lk *TcpLink) error {
	defer utils.PrintPanicStack()
	svr.Lock()
	defer func() {
		svr.Unlock()
	}()

	if len(svr.linkMap) < MAX_TCP_CONN {
		svr.linkMap[i] = lk
		return nil
	} else {
		return errors.New("tcp conn limit")
	}
}

func (svr *TcpServer) GetLink(i GSConfig.SocketIdType) (*TcpLink, bool) {
	defer utils.PrintPanicStack()
	svr.RLock()
	defer func() {
		svr.RUnlock()
	}()
	c, ok := svr.linkMap[i]
	return c, ok
}

// 仅仅只是从map移除，不关闭链接
func (svr *TcpServer) DelLink(i GSConfig.SocketIdType) {
	defer utils.PrintPanicStack()

	svr.Lock()
	defer func() {
		svr.Unlock()
	}()
	delete(svr.linkMap, i)
}

// 会关闭连接
func (svr *TcpServer) KickLink(i GSConfig.SocketIdType) {
	defer utils.PrintPanicStack()
	lk, ok := svr.GetLink(i)
	if ok {
		svr.DelLink(i)
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

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			gLog.Warn(err)
			continue
		}

		gLog.Info("connected: " + tcpConn.RemoteAddr().String() + "  " + strconv.Itoa(len(svr.linkMap)))
		go handleTcpConn(svr, tcpConn)
	}
}

func handleTcpConn(svr *TcpServer, tcpConn *net.TCPConn) {
	defer utils.PrintPanicStack()
	sid  := g.Get()

	lk := NewTcpLink(sid, svr, tcpConn)
	err := svr.PutLink(sid, lk)
	if err != nil {
		lk.conn.Close()
		close(lk.WtSyncChan)
		gLog.Warn(err)
		return
	}

	go lk.StartRead()
	go lk.StartWrite()
	gLog.Info("serving: " + tcpConn.RemoteAddr().String() +  " mapCount: " + strconv.Itoa(len(svr.linkMap)))
}