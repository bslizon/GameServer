package Server

import (
	"Pack"
	gLog "gameLog"
	"fmt"
	"errors"
	"GateServer/config"
)

// 必须保证是线程安全的
func (svr *TcpPackServer) RoutePackIn(p *Pack.Pack) (err error) {
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

	svr.RoutePackOut(p)
	gLog.Debug(fmt.Sprintf("routing. sid: %d data: %v", p.Sid, p.Data))
	////////////////////////////////////////////////////////////////////
	err = nil
	return
}

// 必须保证是线程安全的
func (svr *TcpPackServer) RoutePackOut(p *Pack.Pack) (err error) {
	// panic转error
	defer func() {
		if x := recover(); x != nil {
			switch value := x.(type) {
			case error:
				err = value
			case string:
				err = errors.New(value)
			default:
				err = errors.New(fmt.Sprintf("unknown panic: %#v ", value))
			}
		}
	}()
	////////////////////////////////////////////////////////////////////

	if p.Sid == config.BROCASTING_SID {
		lkMap := svr.GetLinkMapCopy()
		for _, lk := range lkMap {
			er := lk.PutBytes(p.Data)
			if er != nil {
				gLog.Warn(er)// 广播包不返回err
			}
		}
	} else if p.Sid == config.DROP_SID {// 丢弃
		gLog.Warn(fmt.Sprintf("zero sid %#v ", p.Data))
	} else {
		if lk, ok := svr.GetLink(p.Sid); ok {
			er := lk.PutBytes(p.Data)
			if er != nil {
				gLog.Error(er)
				err = errors.New(error.Error(er) + " put wtSyncChan failed.")
				return
			}
		} else {
			err = errors.New("invalid sid.")
			return
		}
	}

	////////////////////////////////////////////////////////////////////
	err = nil
	return
}