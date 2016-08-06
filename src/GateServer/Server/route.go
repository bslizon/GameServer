package Server

import (
	"GateServer/pack"
	gLog "gameLog"
	"fmt"
	"errors"
	"GateServer/config"
)

// 必须保证是线程安全的
func (svr *TcpPackServer) RoutePackIn(p *pack.Pack) (err error) {
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
func (svr *TcpPackServer) RoutePackOut(p *pack.Pack) (err error) {
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
		for _, lk := range svr.linkMap {
			err := lk.PutBytes(p.Data)
			if err != nil {
				gLog.Warn(err)// 广播包不返回err
			}
		}
	} else if p.Sid == config.DROP_SID {// 丢弃
		gLog.Warn(fmt.Sprintf("zero sid %#v ", p.Data))
	} else {
		if lk, ok := svr.GetLink(p.Sid); ok {
			errr := lk.PutBytes(p.Data)
			if errr != nil {
				gLog.Error(errr)
				err = errors.New(errr.Error() + " put wtSyncChan failed.")
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