package Server

import (
	"GateServer/pack"
	"gameLog"
	"fmt"
	"errors"
	"GateServer/config"
)

// 必须保证是线程安全的
func RouteIn(p *pack.Pack) (err error) {
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

	//gameLog.Debug(fmt.Sprintf("routing. sid: %d data: %v", p.Sid, p.Data))
	RouteOut(p)

	////////////////////////////////////////////////////////////////////
	err = nil
	return
}

// 必须保证是线程安全的
func RouteOut(p *pack.Pack) (err error) {
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
		for _, lk := range GateServer.linkMap {
			err := lk.PutBytes(p.Data)
			if err != nil {
				gameLog.Warn(err)
			}
		}
	} else if p.Sid == 0 {// 为0的sid丢弃
		gameLog.Warn(fmt.Sprintf("zero sid %#v", p.Data))
	} else {
		if lk, ok := GateServer.GetLink(p.Sid); ok {
			lk.PutBytes(p.Data)
		} else {
			gameLog.Warn(fmt.Sprintf("link missing, a pack has benn drop. sid: %d data: %v. ", p.Sid, p.Data))
		}
	}

	////////////////////////////////////////////////////////////////////
	err = nil
	return
}